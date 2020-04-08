package output_test

import (
	"github.com/aoldershaw/ansi/action"
	"github.com/aoldershaw/ansi/output"
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

func TestInMemory_Print(t *testing.T) {
	for _, tt := range []struct {
		description string
		initLines   []output.Line
		printCalls  []printCall
		lines       []output.Line
	}{
		{
			description: "printing to a new line creates a new line",
			printCalls: []printCall{
				{
					data: []byte("foo"),
					pos:  action.Pos{Line: 0, Col: 0},
				},
			},
			lines: []output.Line{
				{
					{
						Data: []byte("foo"),
					},
				},
			},
		},
		{
			description: "printing to the second line creates an empty first line",
			printCalls: []printCall{
				{
					data: []byte("foo"),
					pos:  action.Pos{Line: 1, Col: 0},
				},
			},
			lines: []output.Line{
				{},
				{
					{
						Data: []byte("foo"),
					},
				},
			},
		},
		{
			description: "negative position gets set to 0",
			printCalls: []printCall{
				{
					data: []byte("foo"),
					pos:  action.Pos{Line: -1, Col: -1},
				},
			},
			lines: []output.Line{
				{
					{
						Data: []byte("foo"),
					},
				},
			},
		},
		{
			description: "printing to an existing line merges chunks",
			printCalls: []printCall{
				{
					data: []byte("foo"),
					pos:  action.Pos{Line: 0, Col: 0},
				},
				{
					data: []byte("bar"),
					pos:  action.Pos{Line: 0, Col: 3},
				},
			},
			lines: []output.Line{
				{
					{
						Data: []byte("foobar"),
					},
				},
			},
		},
		{
			description: "printing to an empty existing line works",
			initLines: []output.Line{
				{},
			},
			printCalls: []printCall{
				{
					data: []byte("foo"),
					pos:  action.Pos{Line: 0, Col: 0},
				},
			},
			lines: []output.Line{
				{
					{
						Data: []byte("foo"),
					},
				},
			},
		},
		{
			description: "printing to an existing line adds whitespace if cols are not adjacent",
			printCalls: []printCall{
				{
					data: []byte("foo"),
					pos:  action.Pos{Line: 0, Col: 0},
				},
				{
					data: []byte("bar"),
					pos:  action.Pos{Line: 0, Col: 5},
				},
			},
			lines: []output.Line{
				{
					{
						Data: []byte("foo  bar"),
					},
				},
			},
		},
		{
			description: "printing to a line not at Column 0 adds whitespace",
			printCalls: []printCall{
				{
					data: []byte("bar"),
					pos:  action.Pos{Line: 0, Col: 5},
				},
			},
			lines: []output.Line{
				{
					{
						Data: []byte("     bar"),
					},
				},
			},
		},
		{
			description: "overwrites existing chunk if prints overlap",
			printCalls: []printCall{
				{
					data: []byte("foo"),
					pos:  action.Pos{Line: 0, Col: 0},
				},
				{
					data: []byte("bar"),
					pos:  action.Pos{Line: 0, Col: 1},
				},
			},
			lines: []output.Line{
				{
					{
						Data: []byte("fbar"),
					},
				},
			},
		},
		{
			description: "writes inside existing chunk if overlaps",
			printCalls: []printCall{
				{
					data: []byte("foooo"),
					pos:  action.Pos{Line: 0, Col: 0},
				},
				{
					data: []byte("bar"),
					pos:  action.Pos{Line: 0, Col: 1},
				},
			},
			lines: []output.Line{
				{
					{
						Data: []byte("fbaro"),
					},
				},
			},
		},
		{
			description: "does not merge chunks if styles differ",
			printCalls: []printCall{
				{
					data:  []byte("foo"),
					pos:   action.Pos{Line: 0, Col: 0},
					style: style.Style{},
				},
				{
					data:  []byte("bar"),
					pos:   action.Pos{Line: 0, Col: 4},
					style: style.Style{Bold: true},
				},
			},
			lines: []output.Line{
				{
					{
						Data:  []byte("foo "),
						Style: style.Style{},
					},
					{
						Data:  []byte("bar"),
						Style: style.Style{Bold: true},
					},
				},
			},
		},
		{
			description: "write inside first chunk",
			printCalls: []printCall{
				{
					data:  []byte("foo"),
					pos:   action.Pos{Line: 0, Col: 0},
					style: style.Style{},
				},
				{
					data:  []byte("bar"),
					pos:   action.Pos{Line: 0, Col: 3},
					style: style.Style{Bold: true},
				},
				{
					data:  []byte("baz"),
					pos:   action.Pos{Line: 0, Col: 0},
					style: style.Style{},
				},
			},
			lines: []output.Line{
				{
					{
						Data:  []byte("baz"),
						Style: style.Style{},
					},
					{
						Data:  []byte("bar"),
						Style: style.Style{Bold: true},
					},
				},
			},
		},
		{
			description: "if writing inside chunk, but styles differ, splits chunk",
			printCalls: []printCall{
				{
					data:  []byte("abc"),
					pos:   action.Pos{Line: 0, Col: 0},
					style: style.Style{},
				},
				{
					data:  []byte("B"),
					pos:   action.Pos{Line: 0, Col: 1},
					style: style.Style{Bold: true},
				},
			},
			lines: []output.Line{
				{
					{
						Data:  []byte("a"),
						Style: style.Style{},
					},
					{
						Data:  []byte("B"),
						Style: style.Style{Bold: true},
					},
					{
						Data:  []byte("c"),
						Style: style.Style{},
					},
				},
			},
		},
		{
			description: "writing inside chunk with differing styles does not keep empty chunks",
			printCalls: []printCall{
				{
					data:  []byte("abc"),
					pos:   action.Pos{Line: 0, Col: 0},
					style: style.Style{},
				},
				{
					data:  []byte("ABC"),
					pos:   action.Pos{Line: 0, Col: 0},
					style: style.Style{Bold: true},
				},
			},
			lines: []output.Line{
				{
					{
						Data:  []byte("ABC"),
						Style: style.Style{Bold: true},
					},
				},
			},
		},
		{
			description: "overlapping write with chunks after",
			printCalls: []printCall{
				{
					data:  []byte("abc"),
					pos:   action.Pos{Line: 0, Col: 0},
					style: style.Style{},
				},
				{
					data:  []byte("def"),
					pos:   action.Pos{Line: 0, Col: 3},
					style: style.Style{Italic: true},
				},
				{
					data:  []byte("ghi"),
					pos:   action.Pos{Line: 0, Col: 6},
					style: style.Style{Underline: true},
				},
				{
					data:  []byte("BCD"),
					pos:   action.Pos{Line: 0, Col: 1},
					style: style.Style{Bold: true},
				},
			},
			lines: []output.Line{
				{
					{
						Data:  []byte("a"),
						Style: style.Style{},
					},
					{
						Data:  []byte("BCD"),
						Style: style.Style{Bold: true},
					},
					{
						Data:  []byte("ef"),
						Style: style.Style{Italic: true},
					},
					{
						Data:  []byte("ghi"),
						Style: style.Style{Underline: true},
					},
				},
			},
		},
		{
			description: "overlapping write that covers a middle chunk",
			printCalls: []printCall{
				{
					data:  []byte("abc"),
					pos:   action.Pos{Line: 0, Col: 0},
					style: style.Style{},
				},
				{
					data:  []byte("def"),
					pos:   action.Pos{Line: 0, Col: 3},
					style: style.Style{Italic: true},
				},
				{
					data:  []byte("ghi"),
					pos:   action.Pos{Line: 0, Col: 6},
					style: style.Style{Underline: true},
				},
				{
					data:  []byte("CDEFG"),
					pos:   action.Pos{Line: 0, Col: 2},
					style: style.Style{Bold: true},
				},
			},
			lines: []output.Line{
				{
					{
						Data:  []byte("ab"),
						Style: style.Style{},
					},
					{
						Data:  []byte("CDEFG"),
						Style: style.Style{Bold: true},
					},
					{
						Data:  []byte("hi"),
						Style: style.Style{Underline: true},
					},
				},
			},
		},
		{
			description: "overlapping write that ends at the end of a chunk",
			printCalls: []printCall{
				{
					data:  []byte("abc"),
					pos:   action.Pos{Line: 0, Col: 0},
					style: style.Style{},
				},
				{
					data:  []byte("def"),
					pos:   action.Pos{Line: 0, Col: 3},
					style: style.Style{Italic: true},
				},
				{
					data:  []byte("ghi"),
					pos:   action.Pos{Line: 0, Col: 6},
					style: style.Style{Underline: true},
				},
				{
					data:  []byte("CDEF"),
					pos:   action.Pos{Line: 0, Col: 2},
					style: style.Style{Bold: true},
				},
			},
			lines: []output.Line{
				{
					{
						Data:  []byte("ab"),
						Style: style.Style{},
					},
					{
						Data:  []byte("CDEF"),
						Style: style.Style{Bold: true},
					},
					{
						Data:  []byte("ghi"),
						Style: style.Style{Underline: true},
					},
				},
			},
		},
	} {
		t.Run(tt.description, func(t *testing.T) {
			g := NewGomegaWithT(t)
			o := &output.InMemory{Lines: tt.initLines}

			for _, pc := range tt.printCalls {
				o.Print(pc.data, pc.style, pc.pos)
			}

			g.Expect(o.Lines).To(Equal(tt.lines))
		})
	}
}


func TestInMemory_ClearRight(t *testing.T) {
	for _, tt := range []struct {
		description string
		initLines   []output.Line
		clearCalls  []clearCall
		lines       []output.Line
	}{
		{
			description: "clears within a chunk",
			initLines: []output.Line{
				{
					{
						Data: []byte("abcdefghi"),
					},
				},
			},
			clearCalls: []clearCall{
				{
					pos: action.Pos{Line: 0, Col: 3},
				},
			},
			lines: []output.Line{
				{
					{
						Data: []byte("abc"),
					},
				},
			},
		},
		{
			description: "clears multiple chunks",
			initLines: []output.Line{
				{
					{
						Data: []byte("abc"),
					},
					{
						Data: []byte("def"),
					},
					{
						Data: []byte("ghi"),
					},
				},
			},
			clearCalls: []clearCall{
				{
					pos: action.Pos{Line: 0, Col: 2},
				},
			},
			lines: []output.Line{
				{
					{
						Data: []byte("ab"),
					},
				},
			},
		},
		{
			description: "clears from the second chunk on",
			initLines: []output.Line{
				{
					{
						Data: []byte("abc"),
					},
					{
						Data: []byte("def"),
					},
					{
						Data: []byte("ghi"),
					},
				},
			},
			clearCalls: []clearCall{
				{
					pos: action.Pos{Line: 0, Col: 5},
				},
			},
			lines: []output.Line{
				{
					{
						Data: []byte("abc"),
					},
					{
						Data: []byte("de"),
					},
				},
			},
		},
		{
			description: "fully clearing a chunk removes it",
			initLines: []output.Line{
				{
					{
						Data: []byte("abc"),
					},
					{
						Data: []byte("def"),
					},
					{
						Data: []byte("ghi"),
					},
				},
			},
			clearCalls: []clearCall{
				{
					pos: action.Pos{Line: 0, Col: 0},
				},
			},
			lines: []output.Line{
				{},
			},
		},
		{
			description: "clearing an out of bounds line is a noop",
			initLines: []output.Line{},
			clearCalls: []clearCall{
				{
					pos: action.Pos{Line: 0, Col: 0},
				},
			},
			lines: []output.Line{},
		},
		{
			description: "clearing a negative line is a noop",
			initLines: []output.Line{},
			clearCalls: []clearCall{
				{
					pos: action.Pos{Line: -1, Col: 0},
				},
			},
			lines: []output.Line{},
		},
		{
			description: "clearing from a negative column is the same as from 0",
			initLines: []output.Line{
				{
					{
						Data: []byte("abc"),
					},
					{
						Data: []byte("def"),
					},
					{
						Data: []byte("ghi"),
					},
				},
			},
			clearCalls: []clearCall{
				{
					pos: action.Pos{Line: 0, Col: -1},
				},
			},
			lines: []output.Line{
				{},
			},
		},
	} {
		t.Run(tt.description, func(t *testing.T) {
			g := NewGomegaWithT(t)
			o := &output.InMemory{Lines: tt.initLines}

			for _, cc := range tt.clearCalls {
				o.ClearRight(cc.pos)
			}

			g.Expect(o.Lines).To(Equal(tt.lines))
		})
	}
}
