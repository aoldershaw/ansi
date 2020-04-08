package output_test

import (
	"github.com/aoldershaw/ansi"
	"github.com/aoldershaw/ansi/action"
	"github.com/aoldershaw/ansi/output"
	. "github.com/onsi/gomega"
	"testing"
)

type printCall struct {
	data  []byte
	style ansi.Style
	pos   action.Pos
}

func TestOutput_Print_InMemory(t *testing.T) {
	for _, tt := range []struct {
		description string
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
					style: ansi.Style{},
				},
				{
					data:  []byte("bar"),
					pos:   action.Pos{Line: 0, Col: 4},
					style: ansi.Style{Bold: true},
				},
			},
			lines: []output.Line{
				{
					{
						Data:  []byte("foo "),
						Style: ansi.Style{},
					},
					{
						Data:  []byte("bar"),
						Style: ansi.Style{Bold: true},
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
					style: ansi.Style{},
				},
				{
					data:  []byte("bar"),
					pos:   action.Pos{Line: 0, Col: 3},
					style: ansi.Style{Bold: true},
				},
				{
					data:  []byte("baz"),
					pos:   action.Pos{Line: 0, Col: 0},
					style: ansi.Style{},
				},
			},
			lines: []output.Line{
				{
					{
						Data:  []byte("baz"),
						Style: ansi.Style{},
					},
					{
						Data:  []byte("bar"),
						Style: ansi.Style{Bold: true},
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
					style: ansi.Style{},
				},
				{
					data:  []byte("B"),
					pos:   action.Pos{Line: 0, Col: 1},
					style: ansi.Style{Bold: true},
				},
			},
			lines: []output.Line{
				{
					{
						Data:  []byte("a"),
						Style: ansi.Style{},
					},
					{
						Data:  []byte("B"),
						Style: ansi.Style{Bold: true},
					},
					{
						Data:  []byte("c"),
						Style: ansi.Style{},
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
					style: ansi.Style{},
				},
				{
					data:  []byte("ABC"),
					pos:   action.Pos{Line: 0, Col: 0},
					style: ansi.Style{Bold: true},
				},
			},
			lines: []output.Line{
				{
					{
						Data:  []byte("ABC"),
						Style: ansi.Style{Bold: true},
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
					style: ansi.Style{},
				},
				{
					data:  []byte("def"),
					pos:   action.Pos{Line: 0, Col: 3},
					style: ansi.Style{Italic: true},
				},
				{
					data:  []byte("ghi"),
					pos:   action.Pos{Line: 0, Col: 6},
					style: ansi.Style{Underline: true},
				},
				{
					data:  []byte("BCD"),
					pos:   action.Pos{Line: 0, Col: 1},
					style: ansi.Style{Bold: true},
				},
			},
			lines: []output.Line{
				{
					{
						Data:  []byte("a"),
						Style: ansi.Style{},
					},
					{
						Data:  []byte("BCD"),
						Style: ansi.Style{Bold: true},
					},
					{
						Data:  []byte("ef"),
						Style: ansi.Style{Italic: true},
					},
					{
						Data:  []byte("ghi"),
						Style: ansi.Style{Underline: true},
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
					style: ansi.Style{},
				},
				{
					data:  []byte("def"),
					pos:   action.Pos{Line: 0, Col: 3},
					style: ansi.Style{Italic: true},
				},
				{
					data:  []byte("ghi"),
					pos:   action.Pos{Line: 0, Col: 6},
					style: ansi.Style{Underline: true},
				},
				{
					data:  []byte("CDEFG"),
					pos:   action.Pos{Line: 0, Col: 2},
					style: ansi.Style{Bold: true},
				},
			},
			lines: []output.Line{
				{
					{
						Data:  []byte("ab"),
						Style: ansi.Style{},
					},
					{
						Data:  []byte("CDEFG"),
						Style: ansi.Style{Bold: true},
					},
					{
						Data:  []byte("hi"),
						Style: ansi.Style{Underline: true},
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
					style: ansi.Style{},
				},
				{
					data:  []byte("def"),
					pos:   action.Pos{Line: 0, Col: 3},
					style: ansi.Style{Italic: true},
				},
				{
					data:  []byte("ghi"),
					pos:   action.Pos{Line: 0, Col: 6},
					style: ansi.Style{Underline: true},
				},
				{
					data:  []byte("CDEF"),
					pos:   action.Pos{Line: 0, Col: 2},
					style: ansi.Style{Bold: true},
				},
			},
			lines: []output.Line{
				{
					{
						Data:  []byte("ab"),
						Style: ansi.Style{},
					},
					{
						Data:  []byte("CDEF"),
						Style: ansi.Style{Bold: true},
					},
					{
						Data:  []byte("ghi"),
						Style: ansi.Style{Underline: true},
					},
				},
			},
		},
	} {
		t.Run(tt.description, func(t *testing.T) {
			g := NewGomegaWithT(t)
			b := output.NewBuffered(64, nil)

			for _, pc := range tt.printCalls {
				b.Print(pc.data, pc.style, pc.pos)
			}

			g.Expect(b.Lines).To(Equal(tt.lines))
		})
	}
}
