package ansi_test

import (
	"github.com/aoldershaw/ansi"
	"github.com/aoldershaw/ansi/style"
	. "github.com/onsi/gomega"
	"testing"
)

func TestAnsi_Integration_InMemory(t *testing.T) {
	for _, tt := range []struct {
		description string
		events      [][]byte
		lines       []ansi.Line
	}{
		{
			description: "basic test",
			events: [][]byte{
				[]byte("hello\nworld"),
			},
			lines: []ansi.Line{
				{
					{
						Data: ansi.Text("hello"),
					},
				},
				{
					{
						Data: ansi.Text("world"),
					},
				},
			},
		},
		{
			description: "styling",
			events: [][]byte{
				[]byte("hello \x1b[1mworld\x1b[m\n"),
				[]byte("\x1b[31mthis is red\x1b[m\n"),
			},
			lines: []ansi.Line{
				{
					{
						Data: ansi.Text("hello "),
					},
					{
						Data:  []byte("world"),
						Style: style.Style{Bold: true},
					},
				},
				{
					{
						Data:  []byte("this is red"),
						Style: style.Style{Foreground: style.Red},
					},
				},
			},
		},
		{
			description: "styling",
			events: [][]byte{
				[]byte("hello \x1b[1mworld\x1b[m\n"),
				[]byte("\x1b[31mthis is red\x1b[m\n"),
			},
			lines: []ansi.Line{
				{
					{
						Data: ansi.Text("hello "),
					},
					{
						Data:  []byte("world"),
						Style: style.Style{Bold: true},
					},
				},
				{
					{
						Data:  []byte("this is red"),
						Style: style.Style{Foreground: style.Red},
					},
				},
			},
		},
		{
			description: "control sequences split over multiple events",
			events: [][]byte{
				[]byte("\x1b[31mthis is red\x1b"),
				[]byte("[0m but this is not"),
			},
			lines: []ansi.Line{
				{
					{
						Data:  []byte("this is red"),
						Style: style.Style{Foreground: style.Red},
					},
					{
						Data: ansi.Text(" but this is not"),
					},
				},
			},
		},
		{
			description: "moving the cursor",
			events: [][]byte{
				[]byte("hello\x1b[3Cworld"),
				[]byte("\x1b[Ggoodbye"),
			},
			lines: []ansi.Line{
				{
					{
						Data: ansi.Text("goodbye world"),
					},
				},
			},
		},
		{
			description: "save and restore cursor",
			events: [][]byte{
				[]byte("\x1b[shello   world"),
				[]byte("\x1b[ugoodbye"),
			},
			lines: []ansi.Line{
				{
					{
						Data: ansi.Text("goodbye world"),
					},
				},
			},
		},
		{
			description: "erase line",
			events: [][]byte{
				[]byte("this text is very important and will never be removed!\n"),
				[]byte("\x1b[1A\x1b[2Knevermind"),
			},
			lines: []ansi.Line{
				{
					{
						Data: ansi.Text("nevermind"),
					},
				},
			},
		},
	} {
		t.Run(tt.description, func(t *testing.T) {
			g := NewGomegaWithT(t)

			out := &ansi.InMemory{}
			log := ansi.New(out)

			for _, evt := range tt.events {
				log.Parse(evt)
			}

			g.Expect(out.Lines).To(Equal(tt.lines))
		})
	}
}
