package ansi

import "unicode/utf8"

const escapeCode = '\x1b'

type stateFn func(p *Parser, input []byte) stateFn

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
	start int
	pos   int

	currNum maybeInt
	nums    []maybeInt

	state stateFn

	actions  []Action
	action_i int

	dangling []byte
}

func NewParser() *Parser {
	return &Parser{
		// In most cases, this pre-allocation will be plenty
		nums:    make([]maybeInt, 0, 8),
		actions: make([]Action, 0, 8),
		state:   parseBytes,
	}
}

func (p *Parser) Parse(input []byte) (Action, bool, []byte) {
	if p.action_i < len(p.actions) {
		return p.nextAction(), true, input
	}
	p.pos = 0
	p.start = 0

	input = p.extractDangling(input)

	for len(p.actions) == 0 && p.pos < len(input) {
		p.state = p.state(p, input)
	}
	if len(p.actions) == 0 {
		return nil, false, nil
	}
	return p.nextAction(), true, input[p.pos:]
}

// Handle cases where a rune is split up over multiple input events - find the
// boundary for the last complete rune, and mark the incomplete rune as dangling
// for the next input event that comes in
func (p *Parser) extractDangling(input []byte) []byte {
	if len(p.dangling) > 0 {
		// This can be an unfortunate allocation, but it shouldn't matter too much
		// as dangling bytes will likely be pretty rare
		input = append(p.dangling, input...)
	}
	leftover := 0
	for ; leftover < utf8.UTFMax && leftover < len(input); leftover++ {
		r, _ := utf8.DecodeLastRune(input[:len(input)-leftover])
		if r != utf8.RuneError {
			break
		}
	}
	p.dangling = input[len(input)-leftover:]
	return input[:len(input)-leftover]
}

func (p *Parser) ParseAll(input []byte) []Action {
	var actions []Action
	for {
		var (
			action Action
			ok     bool
		)
		action, ok, input = p.Parse(input)
		if !ok {
			break
		}
		actions = append(actions, action)
	}
	return actions
}

func (p *Parser) nextAction() Action {
	a := p.actions[p.action_i]
	if p.action_i == len(p.actions)-1 {
		p.action_i = 0
		p.actions = p.actions[:0]
	} else {
		p.action_i++
	}
	return a
}

func (p *Parser) emit(action Action) {
	p.actions = append(p.actions, action)
	p.start = p.pos
}

func (p *Parser) print(input []byte) {
	data := input[p.start:p.pos]
	printAction := make(Print, len(data))
	copy(printAction, data)
	p.emit(printAction)
}

func (p *Parser) ignore() {
	p.start = p.pos
}

func (p *Parser) next(input []byte) (byte, bool) {
	if p.pos >= len(input) {
		return 0, false
	}
	c := input[p.pos]
	p.pos++
	return c, true
}

func (p *Parser) backup() {
	p.pos--
}

func (p *Parser) peek(input []byte) byte {
	if p.pos >= len(input) {
		return 0
	}
	return input[p.pos]
}

func parseBytes(p *Parser, input []byte) stateFn {
	for {
		switch c := p.peek(input); c {
		case escapeCode:
			if p.pos > p.start {
				p.print(input)
			}
			p.next(input)
			return parseEscapeSequence
		case '\n', '\r':
			if p.pos > p.start {
				p.print(input)
			}
			p.next(input)
			if c == '\n' {
				p.emit(Linebreak{})
			} else {
				p.emit(CarriageReturn{})
			}
			return parseBytes
		}
		if _, ok := p.next(input); !ok {
			break
		}
	}
	if p.pos > p.start {
		p.print(input)
	}
	return parseBytes
}

func parseEscapeSequence(p *Parser, input []byte) stateFn {
	p.nums = p.nums[:0]
	p.currNum = maybeInt{}
	next, ok := p.next(input)
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

func parseControlSequence(p *Parser, input []byte) stateFn {
	var ok bool
	for {
		var d byte
		d, ok = p.next(input)
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

func parseControlSequenceMode(p *Parser, input []byte) stateFn {
	mode, nextOK := p.next(input)
	if !nextOK {
		return parseControlSequence
	}
	var (
		actions []Action
		ok      bool
	)
	if mode == 'm' && len(p.nums) > 2 {
		actions = make([]Action, len(p.nums))
	} else {
		var actionsArr [2]Action
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
		actions[0], ok = CursorUp(num.withDefault(1)), true
	case 'B':
		actions[0], ok = CursorDown(num.withDefault(1)), true
	case 'C':
		actions[0], ok = CursorForward(num.withDefault(1)), true
	case 'D':
		actions[0], ok = CursorBack(num.withDefault(1)), true
	case 'E':
		actions[0], ok = CursorDown(num.withDefault(1)), true
		actions[1] = CursorColumn(0)
	case 'F':
		actions[0], ok = CursorUp(num.withDefault(1)), true
		actions[1] = CursorColumn(0)
	case 'G':
		// This *should* be 1 according to https://en.wikipedia.org/wiki/ANSI_escape_code#Terminal_output_sequences
		// but to match vito/elm-ansi, use 0
		// Note that 0 and 1 seem to behave in the same way
		actions[0], ok = CursorColumn(num.withDefault(0)), true
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
		actions[0], ok = CursorPosition(Pos{
			Line: firstNum.withDefault(1),
			Col:  secondNum.withDefault(1),
		}), true
	case 's':
		actions[0], ok = SaveCursorPosition{}, true
	case 'u':
		actions[0], ok = RestoreCursorPosition{}, true
	case 'J':
		actions[0], ok = EraseDisplay(num.withDefault(0)), true
	case 'K':
		actions[0], ok = EraseLine(num.withDefault(0)), true
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
