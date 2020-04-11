package parser

import (
	"github.com/aoldershaw/ansi/action"
)

const escapeCode = '\x1b'

type stateFn func(p *Parser) stateFn

type maybeInt struct {
	valid bool
	value int
}

func (m maybeInt) withDefault(i int) int {
	if !m.valid {
		return i
	}
	return m.value
}

type Parser struct {
	handler action.Handler

	start int
	pos   int
	input []byte

	currNum maybeInt
	nums    []maybeInt

	state stateFn
}

func New(handler action.Handler) *Parser {
	return &Parser{
		handler: handler,
		// In most cases, this pre-allocation will be plenty
		nums:  make([]maybeInt, 0, 16),
		state: parseBytes,
	}
}

func NewWithChan() (*Parser, <-chan action.Action) {
	c := make(chan action.Action)
	return New(action.HandlerFunc(func(act action.Action) {
		c <- act
	})), c
}

func (p *Parser) Parse(input []byte) {
	p.pos = 0
	p.start = 0
	p.input = input

	for p.pos < len(p.input) {
		p.state = p.state(p)
	}
}

func (p *Parser) emit(action action.Action) {
	p.handler.Action(action)
	p.start = p.pos
}

func (p *Parser) ignore() {
	p.start = p.pos
}

func (p *Parser) next() (byte, bool) {
	if p.pos >= len(p.input) {
		return 0, false
	}
	c := p.input[p.pos]
	p.pos++
	return c, true
}

func (p *Parser) backup() {
	p.pos--
}

func (p *Parser) peek() byte {
	if p.pos >= len(p.input) {
		return 0
	}
	return p.input[p.pos]
}

func parseBytes(p *Parser) stateFn {
	for {
		switch c := p.peek(); c {
		case escapeCode:
			if p.pos > p.start {
				p.emit(action.Print(p.input[p.start:p.pos]))
			}
			p.next()
			return parseEscapeSequence
		case '\n', '\r':
			if p.pos > p.start {
				p.emit(action.Print(p.input[p.start:p.pos]))
			}
			p.next()
			if c == '\n' {
				p.emit(action.Linebreak{})
			} else {
				p.emit(action.CarriageReturn{})
			}
			return parseBytes
		}
		if _, ok := p.next(); !ok {
			break
		}
	}
	if p.pos > p.start {
		p.emit(action.Print(p.input[p.start:p.pos]))
	}
	return parseBytes
}

func parseEscapeSequence(p *Parser) stateFn {
	p.nums = p.nums[:0]
	p.currNum = maybeInt{}
	next, ok := p.next()
	if !ok {
		return parseEscapeSequence
	}
	if next != '[' {
		p.backup()
		p.ignore()
		return parseBytes
	}
	return parseControlSequence
}

func parseControlSequence(p *Parser) stateFn {
	var ok bool
	for {
		var d byte
		d, ok = p.next()
		if !ok {
			return parseControlSequence
		}
		if !isDigit(d) {
			break
		}
		p.currNum.value = 10*p.currNum.value + (int(d) - '0')
		p.currNum.valid = true
	}
	p.nums = append(p.nums, p.currNum)
	p.currNum = maybeInt{}

	p.backup()
	return parseControlSequenceMode
}

func parseControlSequenceMode(p *Parser) stateFn {
	mode, nextOK := p.next()
	if !nextOK {
		return parseControlSequence
	}
	var (
		actions []action.Action
		ok      bool
	)
	if mode == 'm' && len(p.nums) > 2 {
		// TODO: avoid doing dynamic allocation somehow...maybe set a cap on length?
		// 16 would realistically be good for all "normal" cases
		// May not be "correct" but it will likely work fine
		actions = make([]action.Action, len(p.nums))
	} else {
		var actionsArr [2]action.Action
		actions = actionsArr[:]
	}
	var num maybeInt
	if len(p.nums) > 0 {
		num = p.nums[len(p.nums)-1]
	}
	switch mode {
	case 'm':
		for i := 0; i < len(p.nums); i++ {
			var curOk bool
			// If the final parameter is not specified, and it's not the first, don't reset
			// e.g. "\x1b[m" and "\x1b[1;0m" reset, but "\x1b[1;m" sets to bold only (no reset)
			// Not sure where this is in the spec, but it's how iTerm handles it
			if i != 0 && i == len(p.nums)-1 && !p.nums[i].valid {
				break
			}
			actions[i], curOk = sgrLookup(p.nums[i].withDefault(0))
			// "ok" is true if at least one of the actions is valid...otherwise, the whole thing is ignored
			ok = ok || curOk
		}
	case 'A':
		actions[0], ok = action.CursorUp(num.withDefault(1)), true
	case 'B':
		actions[0], ok = action.CursorDown(num.withDefault(1)), true
	case 'C':
		actions[0], ok = action.CursorForward(num.withDefault(1)), true
	case 'D':
		actions[0], ok = action.CursorBack(num.withDefault(1)), true
	case 'E':
		actions[0], ok = action.CursorDown(num.withDefault(1)), true
		actions[1] = action.CursorColumn(0)
	case 'F':
		actions[0], ok = action.CursorUp(num.withDefault(1)), true
		actions[1] = action.CursorColumn(0)
	case 'G':
		// This *should* be 1 according to https://en.wikipedia.org/wiki/ANSI_escape_code#Terminal_output_sequences
		// but to match vito/elm-ansi, use 0
		// Note that 0 and 1 seem to behave in the same way
		actions[0], ok = action.CursorColumn(num.withDefault(0)), true
	case 'H', 'f':
		var (
			firstNum  maybeInt
			secondNum maybeInt
		)
		if len(p.nums) > 0 {
			firstNum = p.nums[0]
		}
		if len(p.nums) > 1 {
			secondNum = p.nums[1]
		}
		actions[0], ok = action.CursorPosition(action.Pos{
			Line: firstNum.withDefault(1),
			Col:  secondNum.withDefault(1),
		}), true
	case 's':
		actions[0], ok = action.SaveCursorPosition{}, true
	case 'u':
		actions[0], ok = action.RestoreCursorPosition{}, true
	case 'J':
		actions[0], ok = action.EraseDisplay(num.withDefault(0)), true
	case 'K':
		actions[0], ok = action.EraseLine(num.withDefault(0)), true
	case ';':
		return parseControlSequence
	default:
		ok = false
	}

	if !ok {
		p.ignore()
		return parseBytes
	}
	for _, act := range actions {
		if act == nil {
			break
		}
		p.emit(act)
	}

	return parseBytes
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
