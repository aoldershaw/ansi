package ansi

import (
	"github.com/aoldershaw/ansi/action"
	"github.com/aoldershaw/ansi/output"
	"github.com/aoldershaw/ansi/style"
)

const (
	defaultLines = 48
	defaultCols  = 80
)

type State struct {
	Style          style.Style
	LineDiscipline LineDiscipline
	Position       action.Pos
	SavedPosition  *action.Pos

	MaxLine int
	MaxCol  int

	output output.Output
}

func New(lineDiscipline LineDiscipline, output output.Output) *State {
	return &State{
		LineDiscipline: lineDiscipline,
		Style:          style.Style{},

		// TODO: functional options for this (and LineDiscipline)
		MaxLine: defaultLines,
		MaxCol:  defaultCols,

		output: output,
	}
}

type LineDiscipline int

const (
	Raw LineDiscipline = iota
	Cooked
)

func (s *State) Action(act action.Action) {
	switch v := act.(type) {
	case action.Print:
		s.output.Print(v, s.Style, s.Position)
		endCol := s.Position.Col + len(v)
		if endCol > s.MaxCol {
			s.MaxCol = endCol
		}
		s.Position.Col = endCol
	case action.Reset:
		s.Style = style.Style{}
	case action.SetForeground:
		s.Style.Foreground = style.Color(v)
	case action.SetBackground:
		s.Style.Background = style.Color(v)
	case action.SetBold:
		s.Style.Bold = bool(v)
	case action.SetFaint:
		s.Style.Faint = bool(v)
	case action.SetItalic:
		s.Style.Italic = bool(v)
	case action.SetUnderline:
		s.Style.Underline = bool(v)
	case action.SetBlink:
		s.Style.Blink = bool(v)
	case action.SetInverted:
		s.Style.Inverted = bool(v)
	case action.SetFraktur:
		s.Style.Fraktur = bool(v)
	case action.SetFramed:
		s.Style.Framed = bool(v)
	case action.CursorPosition:
		s.moveCursorTo(v.Line, v.Col)
	case action.CursorUp:
		s.moveCursor(-int(v), 0)
	case action.CursorDown:
		s.moveCursor(int(v), 0)
	case action.CursorForward:
		s.moveCursor(0, int(v))
	case action.CursorBack:
		s.moveCursor(0, -int(v))
	case action.CursorColumn:
		s.moveCursorTo(s.Position.Line, int(v))
	case action.Linebreak:
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
	case action.CarriageReturn:
		s.Position.Col = 0
	case action.SaveCursorPosition:
		pos := s.Position
		s.SavedPosition = &pos
	case action.RestoreCursorPosition:
		if s.SavedPosition != nil {
			s.Position = *s.SavedPosition
		}
	case action.EraseLine:
		startOfLine := s.Position
		startOfLine.Col = 0
		switch action.EraseMode(v) {
		case action.EraseToBeginning:
			if s.Position.Col == 0 {
				return
			}
			empty := output.Spacer(s.Position.Col)
			s.output.Print(empty, style.Style{}, startOfLine)
		case action.EraseToEnd:
			pos := s.Position
			pos.Col++
			s.output.ClearRight(pos)
		case action.EraseAll:
			s.output.ClearRight(startOfLine)
		}

	case action.EraseDisplay:
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
	s.moveCursorTo(s.Position.Line + dl, s.Position.Col + dc)
}