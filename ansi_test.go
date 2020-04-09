package ansi_test

import (
	"bytes"
	"github.com/aoldershaw/ansi"
	"github.com/aoldershaw/ansi/action"
	"github.com/aoldershaw/ansi/style"
	. "github.com/onsi/gomega"
	"testing"
)

type printCall struct {
	data  []byte
	style style.Style
	pos   action.Pos
}

type clearCall struct {
	pos action.Pos
}

type spyOutput struct {
	printCalls []printCall
	clearCalls []clearCall
}

func (p *spyOutput) Print(data []byte, style style.Style, pos action.Pos) {
	p.printCalls = append(p.printCalls, printCall{
		data:  data,
		style: style,
		pos:   pos,
	})
}

func (p *spyOutput) ClearRight(pos action.Pos) {
	p.clearCalls = append(p.clearCalls, clearCall{
		pos: pos,
	})
}

func TestAnsi_State(t *testing.T) {
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
				action.SetForeground(style.Red),
				action.SetBackground(style.Blue),
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
					style: style.Style{
						Foreground: style.Red,
						Background: style.Blue,
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
				action.SetForeground(style.Red),
				action.SetBold(true),
				action.Reset{},
				action.Print("some unformatted bytes here"),
			},
			printCalls: []printCall{
				{
					data:  []byte("some unformatted bytes here"),
					style: style.Style{},
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
			description: "can't move cursor beyond current size",
			actions: []action.Action{
				action.Print("(0,0)"),
				action.CursorDown(1000),
				action.CursorForward(1000),
				action.Print("(48,80)"),
			},
			printCalls: []printCall{
				{
					data: []byte("(0,0)"),
					pos:  action.Pos{Line: 0, Col: 0},
				},
				{
					data: []byte("(48,80)"),
					pos:  action.Pos{Line: 48, Col: 80},
				},
			},
		},
		{
			description:    "carriage returns/linebreaks in Raw mode",
			lineDiscipline: ansi.Raw,
			actions: []action.Action{
				action.Print("hello"),
				action.Linebreak{},
				action.Print("world"),
				action.CarriageReturn{},
				action.Print("gr8"),
			},
			printCalls: []printCall{
				{
					data: []byte("hello"),
					pos:  action.Pos{Line: 0, Col: 0},
				},
				{
					data: []byte("world"),
					pos:  action.Pos{Line: 1, Col: 5},
				},
				{
					data: []byte("gr8"),
					pos:  action.Pos{Line: 1, Col: 0},
				},
			},
		},
		{
			description:    "carriage returns/linebreaks in Cooked mode",
			lineDiscipline: ansi.Cooked,
			actions: []action.Action{
				action.Print("hello"),
				action.Linebreak{},
				action.Print("world"),
				action.CarriageReturn{},
				action.Print("gr8"),
			},
			printCalls: []printCall{
				{
					data: []byte("hello"),
					pos:  action.Pos{Line: 0, Col: 0},
				},
				{
					data: []byte("world"),
					pos:  action.Pos{Line: 1, Col: 0},
				},
				{
					data: []byte("gr8"),
					pos:  action.Pos{Line: 1, Col: 0},
				},
			},
		},
		{
			description:    "linebreaks can expand max screen height",
			lineDiscipline: ansi.Cooked,
			actions: []action.Action{
				action.Print("this one won't expand the screen size!"),
				action.Linebreak{},
				action.CursorDown(1000),
				action.Print("but this one will!"),
				action.Linebreak{},
				action.CursorUp(1000),
				action.CursorDown(1000),
				action.Print("(49,0)"),
			},
			printCalls: []printCall{
				{
					data: []byte("this one won't expand the screen size!"),
					pos:  action.Pos{Line: 0, Col: 0},
				},
				{
					data: []byte("but this one will!"),
					pos:  action.Pos{Line: 48, Col: 0},
				},
				{
					data: []byte("(49,0)"),
					pos:  action.Pos{Line: 49, Col: 0},
				},
			},
		},
		{
			description:    "prints can expand the screen width",
			lineDiscipline: ansi.Cooked,
			actions: []action.Action{
				action.Print("this print isn't more than 80 chars!"),
				action.CursorColumn(1000),
				action.Print("(0,80)"),
				action.Linebreak{},
				action.CursorForward(1000),
				action.Print("(1,86)"),
			},
			printCalls: []printCall{
				{
					data: []byte("this print isn't more than 80 chars!"),
					pos:  action.Pos{Line: 0, Col: 0},
				},
				{
					data: []byte("(0,80)"),
					pos:  action.Pos{Line: 0, Col: 80},
				},
				{
					data: []byte("(1,86)"),
					pos:  action.Pos{Line: 1, Col: 86},
				},
			},
		},
		{
			description: "can save/restore cursor position",
			actions: []action.Action{
				action.CursorPosition{Line: 12, Col: 34},
				action.SaveCursorPosition{},
				action.CursorPosition{Line: 0, Col: 0},
				action.RestoreCursorPosition{},
				action.Print("i'm back!"),
			},
			printCalls: []printCall{
				{
					data: []byte("i'm back!"),
					pos:  action.Pos{Line: 12, Col: 34},
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
			spyOutput := &spyOutput{}
			state := ansi.New(tt.lineDiscipline, spyOutput)

			for _, act := range tt.actions {
				state.Action(act)
			}

			g.Expect(spyOutput.printCalls).To(Equal(tt.printCalls))
			g.Expect(spyOutput.clearCalls).To(Equal(tt.clearCalls))
		})
	}
}
