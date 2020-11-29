package ansi_test

import (
	"bytes"
	"testing"

	"github.com/aoldershaw/ansi"
	. "github.com/onsi/gomega"
)

type printCall struct {
	data  []byte
	style ansi.Style
	pos   ansi.Pos
}

type clearCall struct {
	pos ansi.Pos
}

type spyOutput struct {
	printCalls []printCall
	clearCalls []clearCall
}

func (p *spyOutput) Print(data []byte, style ansi.Style, pos ansi.Pos) error {
	p.printCalls = append(p.printCalls, printCall{
		data:  data,
		style: style,
		pos:   pos,
	})
	return nil
}

func (p *spyOutput) ClearRight(pos ansi.Pos) error {
	p.clearCalls = append(p.clearCalls, clearCall{
		pos: pos,
	})
	return nil
}

func TestState(t *testing.T) {
	for _, tt := range []struct {
		description    string
		lineDiscipline ansi.LineDiscipline
		actions        []ansi.Action
		printCalls     []printCall
		clearCalls     []clearCall
	}{
		{
			description: "makes print calls",
			actions: []ansi.Action{
				ansi.Print("some bytes here"),
			},
			printCalls: []printCall{
				{
					data: []byte("some bytes here"),
				},
			},
		},
		{
			description: "applies styles",
			actions: []ansi.Action{
				ansi.SetForeground(ansi.Red),
				ansi.SetBackground(ansi.Blue),
				ansi.SetBold(true),
				ansi.SetFraktur(true),
				ansi.SetUnderline(true),
				ansi.SetItalic(true),
				ansi.SetInverted(true),
				ansi.SetFaint(true),
				ansi.SetBlink(true),
				ansi.SetFramed(true),
				ansi.Print("some nicely formatted bytes here"),
			},
			printCalls: []printCall{
				{
					data: []byte("some nicely formatted bytes here"),
					style: ansi.Style{
						Foreground: ansi.Red,
						Background: ansi.Blue,
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
			actions: []ansi.Action{
				ansi.SetForeground(ansi.Red),
				ansi.SetBold(true),
				ansi.Reset{},
				ansi.Print("some unformatted bytes here"),
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
			actions: []ansi.Action{
				ansi.Print("abc"),
				ansi.Print("123"),
			},
			printCalls: []printCall{
				{
					data: []byte("abc"),
					pos:  ansi.Pos{Line: 0, Col: 0},
				},
				{
					data: []byte("123"),
					pos:  ansi.Pos{Line: 0, Col: 3},
				},
			},
		},
		{
			description: "can move cursor",
			actions: []ansi.Action{
				ansi.Print("(0,0)"),
				ansi.CursorDown(2),
				ansi.Print("(2,5)"),
				ansi.CursorUp(1),
				ansi.CursorColumn(15),
				ansi.Print("(1,15)"),
				ansi.CursorPosition{Line: 4, Col: 20},
				ansi.Print("(4,20)"),
				ansi.CursorColumn(10),
				ansi.CursorForward(5),
				ansi.Print("(4,15)"),
				ansi.CursorColumn(10),
				ansi.CursorBack(5),
				ansi.Print("(4,5)"),
				ansi.CursorBack(10000),
				ansi.Print("(4,0)"),
			},
			printCalls: []printCall{
				{
					data: []byte("(0,0)"),
					pos:  ansi.Pos{Line: 0, Col: 0},
				},
				{
					data: []byte("(2,5)"),
					pos:  ansi.Pos{Line: 2, Col: 5},
				},
				{
					data: []byte("(1,15)"),
					pos:  ansi.Pos{Line: 1, Col: 15},
				},
				{
					data: []byte("(4,20)"),
					pos:  ansi.Pos{Line: 4, Col: 20},
				},
				{
					data: []byte("(4,15)"),
					pos:  ansi.Pos{Line: 4, Col: 15},
				},
				{
					data: []byte("(4,5)"),
					pos:  ansi.Pos{Line: 4, Col: 5},
				},
				{
					data: []byte("(4,0)"),
					pos:  ansi.Pos{Line: 4, Col: 0},
				},
			},
		},
		{
			description: "can't move cursor beyond current size",
			actions: []ansi.Action{
				ansi.Print("(0,0)"),
				ansi.CursorDown(1000),
				ansi.CursorForward(1000),
				ansi.Print("(48,80)"),
			},
			printCalls: []printCall{
				{
					data: []byte("(0,0)"),
					pos:  ansi.Pos{Line: 0, Col: 0},
				},
				{
					data: []byte("(48,80)"),
					pos:  ansi.Pos{Line: 48, Col: 80},
				},
			},
		},
		{
			description:    "carriage returns/linebreaks in Raw mode",
			lineDiscipline: ansi.Raw,
			actions: []ansi.Action{
				ansi.Print("hello"),
				ansi.Linebreak{},
				ansi.Print("world"),
				ansi.CarriageReturn{},
				ansi.Print("gr8"),
			},
			printCalls: []printCall{
				{
					data: []byte("hello"),
					pos:  ansi.Pos{Line: 0, Col: 0},
				},
				{
					data: []byte("world"),
					pos:  ansi.Pos{Line: 1, Col: 5},
				},
				{
					data: []byte("gr8"),
					pos:  ansi.Pos{Line: 1, Col: 0},
				},
			},
		},
		{
			description:    "carriage returns/linebreaks in Cooked mode",
			lineDiscipline: ansi.Cooked,
			actions: []ansi.Action{
				ansi.Print("hello"),
				ansi.Linebreak{},
				ansi.Print("world"),
				ansi.CarriageReturn{},
				ansi.Print("gr8"),
			},
			printCalls: []printCall{
				{
					data: []byte("hello"),
					pos:  ansi.Pos{Line: 0, Col: 0},
				},
				{
					data: []byte("world"),
					pos:  ansi.Pos{Line: 1, Col: 0},
				},
				{
					data: []byte("gr8"),
					pos:  ansi.Pos{Line: 1, Col: 0},
				},
			},
		},
		{
			description:    "linebreaks can expand max screen height",
			lineDiscipline: ansi.Cooked,
			actions: []ansi.Action{
				ansi.Print("this one won't expand the screen size!"),
				ansi.Linebreak{},
				ansi.CursorDown(1000),
				ansi.Print("but this one will!"),
				ansi.Linebreak{},
				ansi.CursorUp(1000),
				ansi.CursorDown(1000),
				ansi.Print("(49,0)"),
			},
			printCalls: []printCall{
				{
					data: []byte("this one won't expand the screen size!"),
					pos:  ansi.Pos{Line: 0, Col: 0},
				},
				{
					data: []byte("but this one will!"),
					pos:  ansi.Pos{Line: 48, Col: 0},
				},
				{
					data: []byte("(49,0)"),
					pos:  ansi.Pos{Line: 49, Col: 0},
				},
			},
		},
		{
			description:    "prints can expand the screen width",
			lineDiscipline: ansi.Cooked,
			actions: []ansi.Action{
				ansi.Print("this print isn't more than 80 chars!"),
				ansi.CursorColumn(1000),
				ansi.Print("(0,80)"),
				ansi.Linebreak{},
				ansi.CursorForward(1000),
				ansi.Print("(1,86)"),
			},
			printCalls: []printCall{
				{
					data: []byte("this print isn't more than 80 chars!"),
					pos:  ansi.Pos{Line: 0, Col: 0},
				},
				{
					data: []byte("(0,80)"),
					pos:  ansi.Pos{Line: 0, Col: 80},
				},
				{
					data: []byte("(1,86)"),
					pos:  ansi.Pos{Line: 1, Col: 86},
				},
			},
		},
		{
			description: "can save/restore cursor position",
			actions: []ansi.Action{
				ansi.CursorPosition{Line: 12, Col: 34},
				ansi.SaveCursorPosition{},
				ansi.CursorPosition{Line: 0, Col: 0},
				ansi.RestoreCursorPosition{},
				ansi.Print("i'm back!"),
			},
			printCalls: []printCall{
				{
					data: []byte("i'm back!"),
					pos:  ansi.Pos{Line: 12, Col: 34},
				},
			},
		},
		{
			description: "restoring position when nothing is saved does nothing",
			actions: []ansi.Action{
				ansi.CursorPosition{Line: 1, Col: 2},
				ansi.RestoreCursorPosition{},
				ansi.Print("no change"),
			},
			printCalls: []printCall{
				{
					data: []byte("no change"),
					pos:  ansi.Pos{Line: 1, Col: 2},
				},
			},
		},
		{
			description: "restoring position when nothing is saved does nothing",
			actions: []ansi.Action{
				ansi.CursorPosition{Line: 1, Col: 2},
				ansi.RestoreCursorPosition{},
				ansi.Print("no change"),
			},
			printCalls: []printCall{
				{
					data: []byte("no change"),
					pos:  ansi.Pos{Line: 1, Col: 2},
				},
			},
		},
		{
			description: "erasing lines",
			actions: []ansi.Action{
				ansi.Print("some bytes"),
				ansi.EraseLine(ansi.EraseToBeginning),
				ansi.EraseLine(ansi.EraseToEnd),
				ansi.EraseLine(ansi.EraseAll),
				ansi.Print("some more bytes"),
			},
			printCalls: []printCall{
				{
					data: []byte("some bytes"),
					pos:  ansi.Pos{Line: 0, Col: 0},
				},
				{
					data: bytes.Repeat([]byte{' '}, 10),
					pos:  ansi.Pos{Line: 0, Col: 0},
				},
				{
					data: []byte("some more bytes"),
					pos:  ansi.Pos{Line: 0, Col: 10},
				},
			},
			clearCalls: []clearCall{
				{
					pos: ansi.Pos{Line: 0, Col: 11},
				},
				{
					pos: ansi.Pos{Line: 0, Col: 0},
				},
			},
		},
	} {
		t.Run(tt.description, func(t *testing.T) {
			g := NewGomegaWithT(t)
			spyOutput := &spyOutput{}
			writer := ansi.NewWriter(spyOutput, ansi.WithLineDiscipline(tt.lineDiscipline))

			for _, act := range tt.actions {
				writer.State.Action(act)
			}

			g.Expect(spyOutput.printCalls).To(Equal(tt.printCalls))
			g.Expect(spyOutput.clearCalls).To(Equal(tt.clearCalls))
		})
	}
}
