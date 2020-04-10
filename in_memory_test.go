package ansi_test

import (
	"encoding/json"
	"github.com/aoldershaw/ansi"
	"github.com/aoldershaw/ansi/action"
	"github.com/aoldershaw/ansi/style"
	. "github.com/onsi/gomega"
	"testing"
)

func TestInMemory_Print(t *testing.T) {
	for _, tt := range []struct {
		description string
		initLines   []ansi.Line
		printCalls  []printCall
		lines       []ansi.Line
	}{
		{
			description: "printing to a new line creates a new line",
			printCalls: []printCall{
				{
					data: []byte("foo"),
					pos:  action.Pos{Line: 0, Col: 0},
				},
			},
			lines: []ansi.Line{
				{
					{
						Data: ansi.Text("foo"),
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
			lines: []ansi.Line{
				{},
				{
					{
						Data: ansi.Text("foo"),
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
			lines: []ansi.Line{
				{
					{
						Data: ansi.Text("foo"),
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
			lines: []ansi.Line{
				{
					{
						Data: ansi.Text("foobar"),
					},
				},
			},
		},
		{
			description: "printing to an empty existing line works",
			initLines: []ansi.Line{
				{},
			},
			printCalls: []printCall{
				{
					data: []byte("foo"),
					pos:  action.Pos{Line: 0, Col: 3},
				},
				{
					data: []byte("bar"),
					pos:  action.Pos{Line: 1, Col: 6},
				},
			},
			lines: []ansi.Line{
				{
					{
						Data: ansi.Text("   foo"),
					},
				},
				{
					{
						Data: ansi.Text("      bar"),
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
			lines: []ansi.Line{
				{
					{
						Data: ansi.Text("foo  bar"),
					},
				},
			},
		},
		{
			description: "printing to a line not at Column 0 adds whitespace",
			printCalls: []printCall{
				{
					data: []byte("foo"),
					pos:  action.Pos{Line: 0, Col: 5},
				},
				{
					data: []byte("bar"),
					pos:  action.Pos{Line: 1, Col: 7},
				},
			},
			lines: []ansi.Line{
				{
					{
						Data: ansi.Text("     foo"),
					},
				},
				{
					{
						Data: ansi.Text("       bar"),
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
			lines: []ansi.Line{
				{
					{
						Data: ansi.Text("fbar"),
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
			lines: []ansi.Line{
				{
					{
						Data: ansi.Text("fbaro"),
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
			lines: []ansi.Line{
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
			lines: []ansi.Line{
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
			lines: []ansi.Line{
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
			lines: []ansi.Line{
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
			lines: []ansi.Line{
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
			lines: []ansi.Line{
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
			lines: []ansi.Line{
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
			o := &ansi.InMemory{Lines: tt.initLines}

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
		initLines   []ansi.Line
		clearCalls  []clearCall
		lines       []ansi.Line
	}{
		{
			description: "clears within a chunk",
			initLines: []ansi.Line{
				{
					{
						Data: ansi.Text("abcdefghi"),
					},
				},
			},
			clearCalls: []clearCall{
				{
					pos: action.Pos{Line: 0, Col: 3},
				},
			},
			lines: []ansi.Line{
				{
					{
						Data: ansi.Text("abc"),
					},
				},
			},
		},
		{
			description: "clears multiple chunks",
			initLines: []ansi.Line{
				{
					{
						Data: ansi.Text("abc"),
					},
					{
						Data: ansi.Text("def"),
					},
					{
						Data: ansi.Text("ghi"),
					},
				},
			},
			clearCalls: []clearCall{
				{
					pos: action.Pos{Line: 0, Col: 2},
				},
			},
			lines: []ansi.Line{
				{
					{
						Data: ansi.Text("ab"),
					},
				},
			},
		},
		{
			description: "clears from the second chunk on",
			initLines: []ansi.Line{
				{
					{
						Data: ansi.Text("abc"),
					},
					{
						Data: ansi.Text("def"),
					},
					{
						Data: ansi.Text("ghi"),
					},
				},
			},
			clearCalls: []clearCall{
				{
					pos: action.Pos{Line: 0, Col: 5},
				},
			},
			lines: []ansi.Line{
				{
					{
						Data: ansi.Text("abc"),
					},
					{
						Data: ansi.Text("de"),
					},
				},
			},
		},
		{
			description: "fully clearing a chunk removes it",
			initLines: []ansi.Line{
				{
					{
						Data: ansi.Text("abc"),
					},
					{
						Data: ansi.Text("def"),
					},
					{
						Data: ansi.Text("ghi"),
					},
				},
			},
			clearCalls: []clearCall{
				{
					pos: action.Pos{Line: 0, Col: 0},
				},
			},
			lines: []ansi.Line{
				{},
			},
		},
		{
			description: "clearing an out of bounds line is a noop",
			initLines:   []ansi.Line{},
			clearCalls: []clearCall{
				{
					pos: action.Pos{Line: 0, Col: 0},
				},
			},
			lines: []ansi.Line{},
		},
		{
			description: "clearing a negative line is a noop",
			initLines:   []ansi.Line{},
			clearCalls: []clearCall{
				{
					pos: action.Pos{Line: -1, Col: 0},
				},
			},
			lines: []ansi.Line{},
		},
		{
			description: "clearing from a negative column is the same as from 0",
			initLines: []ansi.Line{
				{
					{
						Data: ansi.Text("abc"),
					},
					{
						Data: ansi.Text("def"),
					},
					{
						Data: ansi.Text("ghi"),
					},
				},
			},
			clearCalls: []clearCall{
				{
					pos: action.Pos{Line: 0, Col: -1},
				},
			},
			lines: []ansi.Line{
				{},
			},
		},
	} {
		t.Run(tt.description, func(t *testing.T) {
			g := NewGomegaWithT(t)
			o := &ansi.InMemory{Lines: tt.initLines}

			for _, cc := range tt.clearCalls {
				o.ClearRight(cc.pos)
			}

			g.Expect(o.Lines).To(Equal(tt.lines))
		})
	}
}

func TestText_MarshalJSON(t *testing.T) {
	g := NewGomegaWithT(t)

	text := ansi.Text("hello world\x1b")
	marshalled, err := text.MarshalJSON()
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(marshalled).To(Equal([]byte(`"hello world\u001b"`)))
}

func TestText_UnmarshalJSON(t *testing.T) {
	g := NewGomegaWithT(t)

	marshalled := []byte(`"hello world\u001b"`)
	var text ansi.Text
	err := json.Unmarshal(marshalled, &text)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(text).To(Equal(ansi.Text("hello world\x1b")))
}
