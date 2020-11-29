package ansi

type LineDiscipline int

const (
	Raw LineDiscipline = iota
	Cooked
)

type State struct {
	Style          Style
	LineDiscipline LineDiscipline
	Position       Pos
	SavedPosition  *Pos

	MaxLine int
	MaxCol  int

	output Output
}

type Output interface {
	Print(data []byte, style Style, pos Pos) error
	ClearRight(pos Pos) error
}

func (s *State) Action(act Action) error {
	switch v := act.(type) {
	case Print:
		if err := s.output.Print(v, s.Style, s.Position); err != nil {
			return err
		}
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
		s.Style.Modifier.applyBit(bool(v), Bold)
	case SetFaint:
		s.Style.Modifier.applyBit(bool(v), Faint)
	case SetItalic:
		s.Style.Modifier.applyBit(bool(v), Italic)
	case SetUnderline:
		s.Style.Modifier.applyBit(bool(v), Underline)
	case SetBlink:
		s.Style.Modifier.applyBit(bool(v), Blink)
	case SetInverted:
		s.Style.Modifier.applyBit(bool(v), Inverted)
	case SetFraktur:
		s.Style.Modifier.applyBit(bool(v), Fraktur)
	case SetFramed:
		s.Style.Modifier.applyBit(bool(v), Framed)
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
				return nil
			}
			empty := spacer(s.Position.Col)
			s.output.Print(empty, Style{}, startOfLine)
		case EraseToEnd:
			pos := s.Position
			pos.Col++
			if err := s.output.ClearRight(pos); err != nil {
				return err
			}
		case EraseAll:
			if err := s.output.ClearRight(startOfLine); err != nil {
				return err
			}
		}

	case EraseDisplay:
		// unsupported
	}

	return nil
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