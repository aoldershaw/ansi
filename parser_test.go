package ansi_test

import (
	"testing"

	"github.com/aoldershaw/ansi"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

func TestParser_Actions(t *testing.T) {
	format.UseStringerRepresentation = true

	for _, tt := range []struct {
		description string
		input       []byte
		actions     []ansi.Action
	}{
		{
			description: "no escape sequences",
			input:       []byte("just some text here!"),
			actions: []ansi.Action{
				ansi.Print("just some text here!"),
			},
		},
		{
			description: "linebreak",
			input:       []byte("hello\nworld"),
			actions: []ansi.Action{
				ansi.Print("hello"),
				ansi.Linebreak{},
				ansi.Print("world"),
			},
		},
		{
			description: "carriage return",
			input:       []byte("hello\rworld"),
			actions: []ansi.Action{
				ansi.Print("hello"),
				ansi.CarriageReturn{},
				ansi.Print("world"),
			},
		},
		{
			description: "colours",
			input:       []byte("normal\x1b[31mred fg\x1b[42mgreen bg\x1b[91mbright red fg\x1b[102mbright green bg"),
			actions: []ansi.Action{
				ansi.Print("normal"),
				ansi.SetForeground(ansi.Red),
				ansi.Print("red fg"),
				ansi.SetBackground(ansi.Green),
				ansi.Print("green bg"),
				ansi.SetForeground(ansi.BrightRed),
				ansi.Print("bright red fg"),
				ansi.SetBackground(ansi.BrightGreen),
				ansi.Print("bright green bg"),
			},
		},
		{
			description: "resetting",
			input:       []byte("some text\x1b[0mreset\x1b[mreset again\x1b[;31mreset to red"),
			actions: []ansi.Action{
				ansi.Print("some text"),
				ansi.Reset{},
				ansi.Print("reset"),
				ansi.Reset{},
				ansi.Print("reset again"),
				ansi.Reset{},
				ansi.SetForeground(ansi.Red),
				ansi.Print("reset to red"),
			},
		},
		{
			description: "text styling",
			input:       []byte("normal\x1b[1mbold\x1b[2mfaint\x1b[3mitalic\x1b[4munderline\x1b[5mblink\x1b[7minverted\x1b[20mfraktur"),
			actions: []ansi.Action{
				ansi.Print("normal"),
				ansi.SetBold(true),
				ansi.Print("bold"),
				ansi.SetFaint(true),
				ansi.Print("faint"),
				ansi.SetItalic(true),
				ansi.Print("italic"),
				ansi.SetUnderline(true),
				ansi.Print("underline"),
				ansi.SetBlink(true),
				ansi.Print("blink"),
				ansi.SetInverted(true),
				ansi.Print("inverted"),
				ansi.SetFraktur(true),
				ansi.Print("fraktur"),
			},
		},
		{
			description: "multiple arguments to formatting",
			input:       []byte("\x1b[1;31;20mhello\x1b[;46m"),
			actions: []ansi.Action{
				ansi.SetBold(true),
				ansi.SetForeground(ansi.Red),
				ansi.SetFraktur(true),
				ansi.Print("hello"),
				ansi.Reset{},
				ansi.SetBackground(ansi.Cyan),
			},
		},
		{
			description: "multiple arguments to formatting, some invalid, OK",
			input:       []byte("\x1b[1;69mhello"),
			actions: []ansi.Action{
				ansi.SetBold(true),
				ansi.Print("hello"),
			},
		},
		{
			description: "multiple arguments to formatting, all invalid, not OK",
			input:       []byte("\x1b[68;69mhello"),
			actions: []ansi.Action{
				ansi.Print("hello"),
			},
		},
		{
			description: "multiple arguments to formatting, last is empty, does not reset",
			input:       []byte("\x1b[1;mhello"),
			actions: []ansi.Action{
				ansi.SetBold(true),
				ansi.Print("hello"),
			},
		},
		{
			description: "cursor movement",
			input:       []byte("\x1b[5A\x1b[50A\x1b[A\x1b[5B\x1b[50B\x1b[B\x1b[5C\x1b[50C\x1b[C\x1b[5D\x1b[50D\x1b[D\x1b[;50H\x1b[50;H\x1b[H\x1b[;H\x1b[50;50H\x1b[;50f\x1b[50;f\x1b[f\x1b[;f\x1b[50;50f"),
			actions: []ansi.Action{
				ansi.CursorUp(5),
				ansi.CursorUp(50),
				ansi.CursorUp(1),
				ansi.CursorDown(5),
				ansi.CursorDown(50),
				ansi.CursorDown(1),
				ansi.CursorForward(5),
				ansi.CursorForward(50),
				ansi.CursorForward(1),
				ansi.CursorBack(5),
				ansi.CursorBack(50),
				ansi.CursorBack(1),
				ansi.CursorPosition(ansi.Pos{1, 50}),
				ansi.CursorPosition(ansi.Pos{50, 1}),
				ansi.CursorPosition(ansi.Pos{1, 1}),
				ansi.CursorPosition(ansi.Pos{1, 1}),
				ansi.CursorPosition(ansi.Pos{50, 50}),
				ansi.CursorPosition(ansi.Pos{1, 50}),
				ansi.CursorPosition(ansi.Pos{50, 1}),
				ansi.CursorPosition(ansi.Pos{1, 1}),
				ansi.CursorPosition(ansi.Pos{1, 1}),
				ansi.CursorPosition(ansi.Pos{50, 50}),
			},
		},
		{
			description: "cursor movement (not ANSI.SYS)",
			input:       []byte("\x1b[E\x1b[5E\x1b[50E\x1b[F\x1b[5F\x1b[50F\x1b[G\x1b[0G\x1b[1G\x1b[5G\x1b[50G"),
			actions: []ansi.Action{
				ansi.CursorDown(1),
				ansi.CursorColumn(0),
				ansi.CursorDown(5),
				ansi.CursorColumn(0),
				ansi.CursorDown(50),
				ansi.CursorColumn(0),
				ansi.CursorUp(1),
				ansi.CursorColumn(0),
				ansi.CursorUp(5),
				ansi.CursorColumn(0),
				ansi.CursorUp(50),
				ansi.CursorColumn(0),
				ansi.CursorColumn(0),
				ansi.CursorColumn(0),
				ansi.CursorColumn(1),
				ansi.CursorColumn(5),
				ansi.CursorColumn(50),
			},
		},
		{
			description: "save/restore cursor",
			input:       []byte("\x1b[s\x1b[u"),
			actions: []ansi.Action{
				ansi.SaveCursorPosition{},
				ansi.RestoreCursorPosition{},
			},
		},
		{
			description: "erasure",
			input:       []byte("\x1b[J\x1b[0J\x1b[1J\x1b[2J\x1b[K\x1b[0K\x1b[1K\x1b[2K"),
			actions: []ansi.Action{
				ansi.EraseDisplay(ansi.EraseToEnd),
				ansi.EraseDisplay(ansi.EraseToEnd),
				ansi.EraseDisplay(ansi.EraseToBeginning),
				ansi.EraseDisplay(ansi.EraseAll),
				ansi.EraseLine(ansi.EraseToEnd),
				ansi.EraseLine(ansi.EraseToEnd),
				ansi.EraseLine(ansi.EraseToBeginning),
				ansi.EraseLine(ansi.EraseAll),
			},
		},
		{
			description: "incomplete escape sequence (no bracket)",
			input:       []byte("hello\x1bworld"),
			actions: []ansi.Action{
				ansi.Print("hello"),
				ansi.Print("world"),
			},
		},
		{
			description: "incomplete escape sequence (double bracket)",
			input:       []byte("hello\x1b[[world"),
			actions: []ansi.Action{
				ansi.Print("hello"),
				ansi.Print("world"),
			},
		},
		{
			description: "incomplete escape sequence at end",
			input:       []byte("hello\x1b"),
			actions: []ansi.Action{
				ansi.Print("hello"),
			},
		},
		{
			description: "unknown escape sequence",
			input:       []byte("hello\x1b[1Zworld"),
			actions: []ansi.Action{
				ansi.Print("hello"),
				ansi.Print("world"),
			},
		},
		{
			description: "something",
			input:       []byte("hello\x1b\n"),
			actions: []ansi.Action{
				ansi.Print("hello"),
				ansi.Linebreak{},
			},
		},
	} {
		t.Run(tt.description, func(t *testing.T) {
			g := NewGomegaWithT(t)
			p := ansi.NewParser()

			actions := p.ParseAll(tt.input)

			g.Expect(actions).To(Equal(tt.actions))
		})
	}
}

func TestParser_Carryover(t *testing.T) {
	format.UseStringerRepresentation = true

	for _, tt := range []struct {
		description string
		input1      []byte
		input2      []byte
		actions     []ansi.Action
	}{
		{
			description: "no partial escapes",
			input1:      []byte("\x1b[31mred text\x1b[m"),
			input2:      []byte("\x1b[32mnow it's green"),
			actions: []ansi.Action{
				ansi.SetForeground(ansi.Red),
				ansi.Print("red text"),
				ansi.Reset{},
				ansi.SetForeground(ansi.Green),
				ansi.Print("now it's green"),
			},
		},
		{
			description: "partial escape sequence",
			input1:      []byte("hello\x1b"),
			input2:      []byte("[32mgreen"),
			actions: []ansi.Action{
				ansi.Print("hello"),
				ansi.SetForeground(ansi.Green),
				ansi.Print("green"),
			},
		},
		{
			description: "partial escape sequence with bracket",
			input1:      []byte("hello\x1b["),
			input2:      []byte("32mgreen"),
			actions: []ansi.Action{
				ansi.Print("hello"),
				ansi.SetForeground(ansi.Green),
				ansi.Print("green"),
			},
		},
		{
			description: "partial escape sequence with bracket and code",
			input1:      []byte("hello\x1b[32"),
			input2:      []byte("mgreen"),
			actions: []ansi.Action{
				ansi.Print("hello"),
				ansi.SetForeground(ansi.Green),
				ansi.Print("green"),
			},
		},
		{
			description: "partial escape sequence with bracket and codes",
			input1:      []byte("hello\x1b[32;1"),
			input2:      []byte("mgreen and bold"),
			actions: []ansi.Action{
				ansi.Print("hello"),
				ansi.SetForeground(ansi.Green),
				ansi.SetBold(true),
				ansi.Print("green and bold"),
			},
		},
		{
			description: "partial escape sequence with bracket codes split up",
			input1:      []byte("hello\x1b[32;"),
			input2:      []byte("1mgreen and bold"),
			actions: []ansi.Action{
				ansi.Print("hello"),
				ansi.SetForeground(ansi.Green),
				ansi.SetBold(true),
				ansi.Print("green and bold"),
			},
		},
		{
			description: "partial escape sequence with code split up",
			input1:      []byte("hello\x1b[3"),
			input2:      []byte("2;1mgreen and bold"),
			actions: []ansi.Action{
				ansi.Print("hello"),
				ansi.SetForeground(ansi.Green),
				ansi.SetBold(true),
				ansi.Print("green and bold"),
			},
		},
	} {
		t.Run(tt.description, func(t *testing.T) {
			g := NewGomegaWithT(t)
			p := ansi.NewParser()

			actions := append(
				p.ParseAll(tt.input1),
				p.ParseAll(tt.input2)...,
			)

			g.Expect(actions).To(Equal(tt.actions))
		})
	}
}
