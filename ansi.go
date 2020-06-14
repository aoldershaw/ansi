package ansi

const (
	defaultLines = 48
	defaultCols  = 80
)

type LogOption func(*Log)

func WithLineDiscipline(d LineDiscipline) LogOption {
	return func(l *Log) {
		l.State.LineDiscipline = d
	}
}

func WithInitialScreenSize(lines, cols int) LogOption {
	return func(l *Log) {
		if lines > 0 {
			l.State.MaxLine = lines
		}
		if cols > 0 {
			l.State.MaxCol = cols
		}
	}
}

type Log struct {
	*Parser
	State *State
}

type LineDiscipline int

const (
	Raw LineDiscipline = iota
	Cooked
)

type Output interface {
	Print(data []byte, style Style, pos Pos)
	ClearRight(pos Pos)
}

type State struct {
	Style          Style
	LineDiscipline LineDiscipline
	Position       Pos
	SavedPosition  *Pos

	MaxLine int
	MaxCol  int

	output Output
}

func New(output Output, opts ...LogOption) *Log {
	state := &State{
		MaxLine: defaultLines,
		MaxCol:  defaultCols,

		LineDiscipline: Cooked,

		output: output,
	}
	log := &Log{
		Parser: NewParser(state),
		State:  state,
	}
	for _, opt := range opts {
		opt(log)
	}
	return log
}

func (s *State) Action(act Action) {
	switch v := act.(type) {
	case Print:
		s.output.Print(v, s.Style, s.Position)
		endCol := s.Position.Col + len(v)
		if endCol > s.MaxCol {
			s.MaxCol = endCol
		}
		s.Position.Col = endCol
	case Reset:
		s.Style = Style{}
	case SetForeground:
		s.Style.Foreground = Color(v)
	case SetBackground:
		s.Style.Background = Color(v)
	case SetBold:
		s.Style.Bold = bool(v)
	case SetFaint:
		s.Style.Faint = bool(v)
	case SetItalic:
		s.Style.Italic = bool(v)
	case SetUnderline:
		s.Style.Underline = bool(v)
	case SetBlink:
		s.Style.Blink = bool(v)
	case SetInverted:
		s.Style.Inverted = bool(v)
	case SetFraktur:
		s.Style.Fraktur = bool(v)
	case SetFramed:
		s.Style.Framed = bool(v)
	case CursorPosition:
		s.moveCursorTo(v.Line, v.Col)
	case CursorUp:
		s.moveCursor(-int(v), 0)
	case CursorDown:
		s.moveCursor(int(v), 0)
	case CursorForward:
		s.moveCursor(0, int(v))
	case CursorBack:
		s.moveCursor(0, -int(v))
	case CursorColumn:
		s.moveCursorTo(s.Position.Line, int(v))
	case Linebreak:
		switch s.LineDiscipline {
		case Raw:
			s.Position.Line++
		case Cooked:
			s.Position.Line++
			s.Position.Col = 0
		}
		if s.Position.Line > s.MaxLine {
			s.MaxLine = s.Position.Line
		}
	case CarriageReturn:
		s.Position.Col = 0
	case SaveCursorPosition:
		pos := s.Position
		s.SavedPosition = &pos
	case RestoreCursorPosition:
		if s.SavedPosition != nil {
			s.Position = *s.SavedPosition
		}
	case EraseLine:
		startOfLine := s.Position
		startOfLine.Col = 0
		switch EraseMode(v) {
		case EraseToBeginning:
			if s.Position.Col == 0 {
				return
			}
			empty := spacer(s.Position.Col)
			s.output.Print(empty, Style{}, startOfLine)
		case EraseToEnd:
			pos := s.Position
			pos.Col++
			s.output.ClearRight(pos)
		case EraseAll:
			s.output.ClearRight(startOfLine)
		}

	case EraseDisplay:
		// unsupported
	}
}

func (s *State) moveCursorTo(l, c int) {
	s.Position.Line = l
	s.Position.Col = c
	if s.Position.Line < 0 {
		s.Position.Line = 0
	}
	if s.Position.Col < 0 {
		s.Position.Col = 0
	}
	if s.Position.Line > s.MaxLine {
		s.Position.Line = s.MaxLine
	}
	if s.Position.Col > s.MaxCol {
		s.Position.Col = s.MaxCol
	}
}

func (s *State) moveCursor(dl, dc int) {
	s.moveCursorTo(s.Position.Line+dl, s.Position.Col+dc)
}
