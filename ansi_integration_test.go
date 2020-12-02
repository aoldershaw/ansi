package ansi_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"

	"github.com/aoldershaw/ansi"
	. "github.com/onsi/gomega"
)

func TestAnsi_Integration_Lines(t *testing.T) {
	for _, tt := range []struct {
		description string
		events      [][]byte
		lines       ansi.Lines
	}{
		{
			description: "basic test",
			events: [][]byte{
				[]byte("hello\nworld"),
			},
			lines: ansi.Lines{
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
			lines: ansi.Lines{
				{
					{
						Data: ansi.Text("hello "),
					},
					{
						Data:  []byte("world"),
						Style: ansi.Style{Modifier: ansi.Bold},
					},
				},
				{
					{
						Data:  []byte("this is red"),
						Style: ansi.Style{Foreground: ansi.Red},
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
			lines: ansi.Lines{
				{
					{
						Data:  []byte("this is red"),
						Style: ansi.Style{Foreground: ansi.Red},
					},
					{
						Data: ansi.Text(" but this is not"),
					},
				},
			},
		},
		{
			description: "runes that are split over multiple events",
			events: [][]byte{
				[]byte("hello \xe3\x81"),
				[]byte("\x93 world"),
			},
			lines: ansi.Lines{
				{
					{
						Data: ansi.Text("hello こ world"),
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
			lines: ansi.Lines{
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
			lines: ansi.Lines{
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
			lines: ansi.Lines{
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

			var lines ansi.Lines
			writer := ansi.NewWriter(&lines)

			initialEvents := make([][]byte, len(tt.events))
			for i, evt := range tt.events {
				initialEvents[i] = make([]byte, len(evt))
				copy(initialEvents[i], evt)
			}

			for _, evt := range tt.events {
				_, err := writer.Write(evt)
				g.Expect(err).ToNot(HaveOccurred())
			}

			g.Expect(lines).To(Equal(tt.lines))
			g.Expect(tt.events).To(Equal(initialEvents), "modified input bytes")
		})
	}
}

func benchmark(b *testing.B, numEvents int, numBytesPerEvent int, probOfControlSequence float64) {
	b.Helper()

	r := rand.New(rand.NewSource(456))
	events := make([][]byte, numEvents)
	for i := 0; i < len(events); i++ {
		events[i] = generateEvent(r, numBytesPerEvent, probOfControlSequence)
	}
	b.SetBytes(int64(numEvents * numBytesPerEvent))
	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		var lines ansi.Lines
		writer := ansi.NewWriter(&lines)

		for _, evt := range events {
			_, err := writer.Write(evt)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func Benchmark_1_4096x80_5(b *testing.B) {
	benchmark(b, 1, 4096*80, 0.05)
}

func Benchmark_4096_80_5(b *testing.B) {
	benchmark(b, 4096, 80, 0.05)
}

func Benchmark_4096_120_5(b *testing.B) {
	benchmark(b, 4096, 120, 0.05)
}

func Benchmark_4096_80_10(b *testing.B) {
	benchmark(b, 4096, 80, 0.10)
}

func Benchmark_8192_80_5(b *testing.B) {
	benchmark(b, 8192, 80, 0.05)
}

const modes = "mABCDEFGHfsuJK"
const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789\t\n\r "

func generateEvent(r *rand.Rand, length int, probOfControlSequence float64) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, length))
	for i := 0; i < length; i++ {
		if r.Float64() < probOfControlSequence {
			mode := modes[r.Int()%len(modes)]
			n, _ := buf.WriteString(fmt.Sprintf("\x1b[%d%c", r.Int()%40, mode))
			i += n
		} else {
			buf.WriteByte(chars[r.Int()%len(chars)])
			i++
		}
	}
	return buf.Bytes()
}

func Example() {
	var lines ansi.Lines
	writer := ansi.NewWriter(&lines)

	writer.Write([]byte("\x1b[1mbold\x1b[m text"))
	writer.Write([]byte("\nline 2"))

	linesJSON, _ := json.Marshal(lines)
	fmt.Println(string(linesJSON))
	// Output: [[{"data":"bold","style":{"bold":true}},{"data":" text","style":{}}],[{"data":"line 2","style":{}}]]
}
