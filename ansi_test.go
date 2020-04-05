package ansi_test

import (
	"bytes"
	"github.com/aoldershaw/ansi"
	"github.com/aoldershaw/ansi/action"
	. "github.com/onsi/gomega"
	"testing"
)

type printCall struct {
	data  []byte
	style ansi.Style
	pos   action.Pos
}

type clearCall struct {
	pos action.Pos
}

type spyPrinter struct {
	printCalls []printCall
	clearCalls []clearCall
}

func (p *spyPrinter) Print(data []byte, style ansi.Style, pos action.Pos) {
	p.printCalls = append(p.printCalls, printCall{
		data:  data,
		style: style,
		pos:   pos,
	})
}

func (p *spyPrinter) ClearRight(pos action.Pos) {
	p.clearCalls = append(p.clearCalls, clearCall{
		pos: pos,
	})
}

func TestAnsi(t *testing.T) {
	for _, tt := range []struct {
		description    string
		lineDiscipline ansi.LineDiscipline
		actions        []action.Action
		printCalls     []printCall
		clearCalls     []clearCall
	}{
		{
			description: "makes print calls",
			actions: []action.Action{
				action.Print("some bytes here"),
			},
			printCalls: []printCall{
				{
					data: []byte("some bytes here"),
				},
			},
		},
		{
			description: "applies styles",
			actions: []action.Action{
				action.SetForeground(action.Red),
				action.SetBackground(action.Blue),
				action.SetBold(true),
				action.SetFraktur(true),
				action.SetUnderline(true),
				action.SetItalic(true),
				action.SetInverted(true),
				action.SetFaint(true),
				action.SetBlink(true),
				action.SetFramed(true),
				action.Print("some nicely formatted bytes here"),
			},
			printCalls: []printCall{
				{
					data: []byte("some nicely formatted bytes here"),
					style: ansi.Style{
						Foreground: action.Red,
						Background: action.Blue,
						Bold:       true,
						Faint:      true,
						Italic:     true,
						Underline:  true,
						Blink:      true,
						Inverted:   true,
						Fraktur:    true,
						Framed:     true,
					},
				},
			},
		},
		{
			description: "resets styles",
			actions: []action.Action{
				action.SetForeground(action.Red),
				action.SetBold(true),
				action.Reset{},
				action.Print("some unformatted bytes here"),
			},
			printCalls: []printCall{
				{
					data:  []byte("some unformatted bytes here"),
					style: ansi.Style{},
				},
			},
		},
		{
			description: "print calls move cursor",
			actions: []action.Action{
				action.Print("abc"),
				action.Print("123"),
			},
			printCalls: []printCall{
				{
					data: []byte("abc"),
					pos:  action.Pos{Line: 0, Col: 0},
				},
				{
					data: []byte("123"),
					pos:  action.Pos{Line: 0, Col: 3},
				},
			},
		},
		{
			description: "can move cursor",
			actions: []action.Action{
				action.Print("(0,0)"),
				action.CursorDown(2),
				action.Print("(2,5)"),
				action.CursorUp(1),
				action.CursorColumn(15),
				action.Print("(1,15)"),
				action.CursorPosition{Line: 4, Col: 20},
				action.Print("(4,20)"),
				action.CursorColumn(10),
				action.CursorForward(5),
				action.Print("(4,15)"),
				action.CursorColumn(10),
				action.CursorBack(5),
				action.Print("(4,5)"),
				action.CursorBack(10000),
				action.Print("(4,0)"),
			},
			printCalls: []printCall{
				{
					data: []byte("(0,0)"),
					pos:  action.Pos{Line: 0, Col: 0},
				},
				{
					data: []byte("(2,5)"),
					pos:  action.Pos{Line: 2, Col: 5},
				},
				{
					data: []byte("(1,15)"),
					pos:  action.Pos{Line: 1, Col: 15},
				},
				{
					data: []byte("(4,20)"),
					pos:  action.Pos{Line: 4, Col: 20},
				},
				{
					data: []byte("(4,15)"),
					pos:  action.Pos{Line: 4, Col: 15},
				},
				{
					data: []byte("(4,5)"),
					pos:  action.Pos{Line: 4, Col: 5},
				},
				{
					data: []byte("(4,0)"),
					pos:  action.Pos{Line: 4, Col: 0},
				},
			},
		},
		{
			description: "can save/restore cursor position",
			actions: []action.Action{
				action.CursorPosition{Line: 123, Col: 456},
				action.SaveCursorPosition{},
				action.CursorPosition{Line: 0, Col: 0},
				action.RestoreCursorPosition{},
				action.Print("i'm back!"),
			},
			printCalls: []printCall{
				{
					data: []byte("i'm back!"),
					pos:  action.Pos{Line: 123, Col: 456},
				},
			},
		},
		{
			description: "restoring position when nothing is saved does nothing",
			actions: []action.Action{
				action.CursorPosition{Line: 1, Col: 2},
				action.RestoreCursorPosition{},
				action.Print("no change"),
			},
			printCalls: []printCall{
				{
					data: []byte("no change"),
					pos:  action.Pos{Line: 1, Col: 2},
				},
			},
		},
		{
			description: "restoring position when nothing is saved does nothing",
			actions: []action.Action{
				action.CursorPosition{Line: 1, Col: 2},
				action.RestoreCursorPosition{},
				action.Print("no change"),
			},
			printCalls: []printCall{
				{
					data: []byte("no change"),
					pos:  action.Pos{Line: 1, Col: 2},
				},
			},
		},
		{
			description: "erasing lines",
			actions: []action.Action{
				action.Print("some bytes"),
				action.EraseLine(action.EraseToBeginning),
				action.EraseLine(action.EraseToEnd),
				action.EraseLine(action.EraseAll),
				action.Print("some more bytes"),
			},
			printCalls: []printCall{
				{
					data: []byte("some bytes"),
					pos:  action.Pos{Line: 0, Col: 0},
				},
				{
					data: bytes.Repeat([]byte{' '}, 10),
					pos:  action.Pos{Line: 0, Col: 0},
				},
				{
					data: []byte("some more bytes"),
					pos:  action.Pos{Line: 0, Col: 10},
				},
			},
			clearCalls: []clearCall{
				{
					pos: action.Pos{Line: 0, Col: 11},
				},
				{
					pos: action.Pos{Line: 0, Col: 0},
				},
			},
		},
	} {
		t.Run(tt.description, func(t *testing.T) {
			g := NewGomegaWithT(t)
			spyPrinter := &spyPrinter{}
			state := ansi.New(tt.lineDiscipline, spyPrinter)

			for _, act := range tt.actions {
				state.Action(act)
			}

			g.Expect(spyPrinter.printCalls).To(Equal(tt.printCalls))
			g.Expect(spyPrinter.clearCalls).To(Equal(tt.clearCalls))
		})
	}
}
