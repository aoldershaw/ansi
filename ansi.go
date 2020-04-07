package ansi

import (
	"bytes"
	"github.com/aoldershaw/ansi/action"
)

type Output interface {
	Print(data []byte, style Style, pos action.Pos)
	ClearRight(pos action.Pos)
}

type State struct {
	Style          Style
	LineDiscipline LineDiscipline
	Position       action.Pos
	SavedPosition  *action.Pos

	output Output
}

func New(lineDiscipline LineDiscipline, output Output) *State {
	return &State{
		LineDiscipline: lineDiscipline,
		Style:          Style{},

		output: output,
	}
}

type Style struct {
	Foreground action.Color `json:"fg,omitempty"`
	Background action.Color `json:"bg,omitempty"`
	Bold       bool         `json:"bold,omitempty"`
	Faint      bool         `json:"faint,omitempty"`
	Italic     bool         `json:"italic,omitempty"`
	Underline  bool         `json:"underline,omitempty"`
	Blink      bool         `json:"blink,omitempty"`
	Inverted   bool         `json:"inverted,omitempty"`
	Fraktur    bool         `json:"fraktur,omitempty"`
	Framed     bool         `json:"framed,omitempty"`
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
		s.moveCursor(0, len(v))
	case action.Reset:
		s.Style = Style{}
	case action.SetForeground:
		s.Style.Foreground = action.Color(v)
	case action.SetBackground:
		s.Style.Background = action.Color(v)
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
		s.Position = action.Pos(v)
	case action.CursorUp:
		s.moveCursor(-int(v), 0)
	case action.CursorDown:
		s.moveCursor(int(v), 0)
	case action.CursorForward:
		s.moveCursor(0, int(v))
	case action.CursorBack:
		s.moveCursor(0, -int(v))
	case action.CursorColumn:
		s.Position.Col = int(v)
	case action.Linebreak:
		switch s.LineDiscipline {
		case Raw:
			s.moveCursor(1, 0)
		case Cooked:
			s.Position.Line++
			s.Position.Col = 0
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
			empty := bytes.Repeat([]byte{' '}, s.Position.Col)
			s.output.Print(empty, Style{}, startOfLine)
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

func (s *State) moveCursor(r, c int) {
	s.Position.Line += r
	s.Position.Col += c
	if s.Position.Line < 0 {
		s.Position.Line = 0
	}
	if s.Position.Col < 0 {
		s.Position.Col = 0
	}
}