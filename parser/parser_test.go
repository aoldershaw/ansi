package parser_test

import (
	"github.com/aoldershaw/ansi/action"
	"github.com/aoldershaw/ansi/parser"
	"github.com/aoldershaw/ansi/style"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"testing"
)

type spyHandler struct {
	actions []action.Action
}

func (s *spyHandler) Action(a action.Action) {
	s.actions = append(s.actions, a)
}

func TestParser_Actions(t *testing.T) {
	format.UseStringerRepresentation = true

	for _, tt := range []struct {
		description string
		input       []byte
		actions     []action.Action
	}{
		{
			description: "no escape sequences",
			input:       []byte("just some text here!"),
			actions: []action.Action{
				action.Print("just some text here!"),
			},
		},
		{
			description: "linebreak",
			input:       []byte("hello\nworld"),
			actions: []action.Action{
				action.Print("hello"),
				action.Linebreak{},
				action.Print("world"),
			},
		},
		{
			description: "carriage return",
			input:       []byte("hello\rworld"),
			actions: []action.Action{
				action.Print("hello"),
				action.CarriageReturn{},
				action.Print("world"),
			},
		},
		{
			description: "colours",
			input:       []byte("normal\x1b[31mred fg\x1b[42mgreen bg\x1b[91mbright red fg\x1b[102mbright green bg"),
			actions: []action.Action{
				action.Print("normal"),
				action.SetForeground(style.Red),
				action.Print("red fg"),
				action.SetBackground(style.Green),
				action.Print("green bg"),
				action.SetForeground(style.BrightRed),
				action.Print("bright red fg"),
				action.SetBackground(style.BrightGreen),
				action.Print("bright green bg"),
			},
		},
		{
			description: "resetting",
			input:       []byte("some text\x1b[0mreset\x1b[mreset again\x1b[;31mreset to red"),
			actions: []action.Action{
				action.Print("some text"),
				action.Reset{},
				action.Print("reset"),
				action.Reset{},
				action.Print("reset again"),
				action.Reset{},
				action.SetForeground(style.Red),
				action.Print("reset to red"),
			},
		},
		{
			description: "text styling",
			input:       []byte("normal\x1b[1mbold\x1b[2mfaint\x1b[3mitalic\x1b[4munderline\x1b[5mblink\x1b[7minverted\x1b[20mfraktur"),
			actions: []action.Action{
				action.Print("normal"),
				action.SetBold(true),
				action.Print("bold"),
				action.SetFaint(true),
				action.Print("faint"),
				action.SetItalic(true),
				action.Print("italic"),
				action.SetUnderline(true),
				action.Print("underline"),
				action.SetBlink(true),
				action.Print("blink"),
				action.SetInverted(true),
				action.Print("inverted"),
				action.SetFraktur(true),
				action.Print("fraktur"),
			},
		},
		{
			description: "multiple arguments to formatting",
			input:       []byte("\x1b[1;31;20mhello\x1b[;46m"),
			actions: []action.Action{
				action.SetBold(true),
				action.SetForeground(style.Red),
				action.SetFraktur(true),
				action.Print("hello"),
				action.Reset{},
				action.SetBackground(style.Cyan),
			},
		},
		{
			description: "multiple arguments to formatting, some invalid, OK",
			input:       []byte("\x1b[1;69mhello"),
			actions: []action.Action{
				action.SetBold(true),
				action.Print("hello"),
			},
		},
		{
			description: "multiple arguments to formatting, all invalid, not OK",
			input:       []byte("\x1b[68;69mhello"),
			actions: []action.Action{
				action.Print("hello"),
			},
		},
		{
			description: "multiple arguments to formatting, last is empty, does not reset",
			input:       []byte("\x1b[1;mhello"),
			actions: []action.Action{
				action.SetBold(true),
				action.Print("hello"),
			},
		},
		{
			description: "cursor movement",
			input:       []byte("\x1b[5A\x1b[50A\x1b[A\x1b[5B\x1b[50B\x1b[B\x1b[5C\x1b[50C\x1b[C\x1b[5D\x1b[50D\x1b[D\x1b[;50H\x1b[50;H\x1b[H\x1b[;H\x1b[50;50H\x1b[;50f\x1b[50;f\x1b[f\x1b[;f\x1b[50;50f"),
			actions: []action.Action{
				action.CursorUp(5),
				action.CursorUp(50),
				action.CursorUp(1),
				action.CursorDown(5),
				action.CursorDown(50),
				action.CursorDown(1),
				action.CursorForward(5),
				action.CursorForward(50),
				action.CursorForward(1),
				action.CursorBack(5),
				action.CursorBack(50),
				action.CursorBack(1),
				action.CursorPosition(action.Pos{1, 50}),
				action.CursorPosition(action.Pos{50, 1}),
				action.CursorPosition(action.Pos{1, 1}),
				action.CursorPosition(action.Pos{1, 1}),
				action.CursorPosition(action.Pos{50, 50}),
				action.CursorPosition(action.Pos{1, 50}),
				action.CursorPosition(action.Pos{50, 1}),
				action.CursorPosition(action.Pos{1, 1}),
				action.CursorPosition(action.Pos{1, 1}),
				action.CursorPosition(action.Pos{50, 50}),
			},
		},
		{
			description: "cursor movement (not ANSI.SYS)",
			input:       []byte("\x1b[E\x1b[5E\x1b[50E\x1b[F\x1b[5F\x1b[50F\x1b[G\x1b[0G\x1b[1G\x1b[5G\x1b[50G"),
			actions: []action.Action{
				action.CursorDown(1),
				action.CursorColumn(0),
				action.CursorDown(5),
				action.CursorColumn(0),
				action.CursorDown(50),
				action.CursorColumn(0),
				action.CursorUp(1),
				action.CursorColumn(0),
				action.CursorUp(5),
				action.CursorColumn(0),
				action.CursorUp(50),
				action.CursorColumn(0),
				action.CursorColumn(0),
				action.CursorColumn(0),
				action.CursorColumn(1),
				action.CursorColumn(5),
				action.CursorColumn(50),
			},
		},
		{
			description: "save/restore cursor",
			input:       []byte("\x1b[s\x1b[u"),
			actions: []action.Action{
				action.SaveCursorPosition{},
				action.RestoreCursorPosition{},
			},
		},
		{
			description: "erasure",
			input:       []byte("\x1b[J\x1b[0J\x1b[1J\x1b[2J\x1b[K\x1b[0K\x1b[1K\x1b[2K"),
			actions: []action.Action{
				action.EraseDisplay(action.EraseToEnd),
				action.EraseDisplay(action.EraseToEnd),
				action.EraseDisplay(action.EraseToBeginning),
				action.EraseDisplay(action.EraseAll),
				action.EraseLine(action.EraseToEnd),
				action.EraseLine(action.EraseToEnd),
				action.EraseLine(action.EraseToBeginning),
				action.EraseLine(action.EraseAll),
			},
		},
		{
			description: "incomplete escape sequence (no bracket)",
			input:       []byte("hello\x1bworld"),
			actions: []action.Action{
				action.Print("hello"),
				action.Print("world"),
			},
		},
		{
			description: "incomplete escape sequence (double bracket)",
			input:       []byte("hello\x1b[[world"),
			actions: []action.Action{
				action.Print("hello"),
				action.Print("world"),
			},
		},
		{
			description: "incomplete escape sequence at end",
			input:       []byte("hello\x1b"),
			actions: []action.Action{
				action.Print("hello"),
			},
		},
		{
			description: "unknown escape sequence",
			input:       []byte("hello\x1b[1Zworld"),
			actions: []action.Action{
				action.Print("hello"),
				action.Print("world"),
			},
		},
		{
			description: "something",
			input:       []byte("hello\x1b\n"),
			actions: []action.Action{
				action.Print("hello"),
				action.Linebreak{},
			},
		},
	} {
		t.Run(tt.description, func(t *testing.T) {
			g := NewGomegaWithT(t)
			handler := &spyHandler{}
			p := parser.New(handler)

			p.Parse(tt.input)

			g.Expect(handler.actions).To(Equal(tt.actions))
		})
	}
}

func TestParser_Carryover(t *testing.T) {
	format.UseStringerRepresentation = true

	for _, tt := range []struct {
		description string
		input1      []byte
		input2      []byte
		actions     []action.Action
	}{
		{
			description: "no partial escapes",
			input1:      []byte("\x1b[31mred text\x1b[m"),
			input2:      []byte("\x1b[32mnow it's green"),
			actions: []action.Action{
				action.SetForeground(style.Red),
				action.Print("red text"),
				action.Reset{},
				action.SetForeground(style.Green),
				action.Print("now it's green"),
			},
		},
		{
			description: "partial escape sequence",
			input1:      []byte("hello\x1b"),
			input2:      []byte("[32mgreen"),
			actions: []action.Action{
				action.Print("hello"),
				action.SetForeground(style.Green),
				action.Print("green"),
			},
		},
		{
			description: "partial escape sequence with bracket",
			input1:      []byte("hello\x1b["),
			input2:      []byte("32mgreen"),
			actions: []action.Action{
				action.Print("hello"),
				action.SetForeground(style.Green),
				action.Print("green"),
			},
		},
		{
			description: "partial escape sequence with bracket and code",
			input1:      []byte("hello\x1b[32"),
			input2:      []byte("mgreen"),
			actions: []action.Action{
				action.Print("hello"),
				action.SetForeground(style.Green),
				action.Print("green"),
			},
		},
		{
			description: "partial escape sequence with bracket and codes",
			input1:      []byte("hello\x1b[32;1"),
			input2:      []byte("mgreen and bold"),
			actions: []action.Action{
				action.Print("hello"),
				action.SetForeground(style.Green),
				action.SetBold(true),
				action.Print("green and bold"),
			},
		},
		{
			description: "partial escape sequence with bracket codes split up",
			input1:      []byte("hello\x1b[32;"),
			input2:      []byte("1mgreen and bold"),
			actions: []action.Action{
				action.Print("hello"),
				action.SetForeground(style.Green),
				action.SetBold(true),
				action.Print("green and bold"),
			},
		},
		{
			description: "partial escape sequence with code split up",
			input1:      []byte("hello\x1b[3"),
			input2:      []byte("2;1mgreen and bold"),
			actions: []action.Action{
				action.Print("hello"),
				action.SetForeground(style.Green),
				action.SetBold(true),
				action.Print("green and bold"),
			},
		},
	} {
		t.Run(tt.description, func(t *testing.T) {
			g := NewGomegaWithT(t)
			handler := &spyHandler{}
			p := parser.New(handler)

			p.Parse(tt.input1)
			p.Parse(tt.input2)

			g.Expect(handler.actions).To(Equal(tt.actions))
		})
	}
}
