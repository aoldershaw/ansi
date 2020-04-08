package output

import (
	"bytes"
	"github.com/aoldershaw/ansi"
	"github.com/aoldershaw/ansi/action"
)

// Used as an optimization to avoid heap allocations
const spacerBytesSize = 512

var spacerPreallocBytes [spacerBytesSize]byte

func init() {
	for i := 0; i < spacerBytesSize; i++ {
		spacerPreallocBytes[i] = ' '
	}
}

type Chunk struct {
	Data  []byte     `json:"data"`
	Style ansi.Style `json:"style,omitempty"`
}

type Line = []Chunk

type Storage interface {
	Put(firstLine int, lines []Line) error
	Get(firstLine int, count int) ([]Line, error)
}

type Buffered struct {
	Lines []Line

	storage Storage
}

func NewBuffered(numLines int, storage Storage) *Buffered {
	return &Buffered{
		Lines: make([]Line, 0, numLines),
	}
}

func (b *Buffered) Print(data []byte, style ansi.Style, pos action.Pos) {
	if pos.Line < 0 {
		pos.Line = 0
	}
	if pos.Col < 0 {
		pos.Col = 0
	}
	numEmpty := pos.Line - len(b.Lines)
	for numEmpty > 0 {
		b.Lines = append(b.Lines, Line{})
		numEmpty--
	}
	// TODO: need to shift by current start line
	if pos.Line >= len(b.Lines) {
		b.Lines = append(b.Lines, Line{{Data: data, Style: style}})
		return
	}

	lineLen := b.lineLength(pos.Line)

	if pos.Col >= lineLen {
		b.appendToLine(data, style, pos)
	} else {
		b.insertWithinLine(data, style, pos)
	}
}

func (b *Buffered) appendToLine(data []byte, style ansi.Style, pos action.Pos) {
	line := b.Lines[pos.Line]

	lineLen := b.lineLength(pos.Line)
	spacerLen := pos.Col - lineLen

	lastChunk := &line[len(line)-1]
	if spacerLen > 0 {
		lastChunk.Data = append(lastChunk.Data, spacer(spacerLen)...)
	}
	if lastChunk.Style == style {
		lastChunk.Data = append(lastChunk.Data, data...)
		return
	}
	b.Lines[pos.Line] = append(line, Chunk{Data: data, Style: style})
}

type interval struct {
	L int
	R int
}

func intervalWithWidth(start int, width int) interval {
	return interval{L: start, R: start + width - 1}
}

func (i interval) contains(i2 interval) bool {
	return i.L <= i2.L && i2.R <= i.R
}

func (i interval) overlaps(i2 interval) bool {
	return i.R >= i2.L && i.L <= i2.R
}

func (b *Buffered) insertWithinLine(data []byte, style ansi.Style, pos action.Pos) {
	line := b.Lines[pos.Line]
	chunkStart := 0
	for i := 0; i < len(line); i++ {
		chunk := line[i]
		chunkInterval := intervalWithWidth(chunkStart, len(chunk.Data))
		printInterval := intervalWithWidth(pos.Col, len(data))

		relCol := pos.Col - chunkStart
		chunkStart += len(chunk.Data)

		if !chunkInterval.overlaps(printInterval) {
			continue
		}

		if chunkInterval.contains(printInterval) {
			b.insertInsideChunk(data, style, pos.Line, relCol, i)
			return
		}
		newLine := append(make(Line, 0, len(line)+1), line[:i]...)
		originalChunkLength := len(chunk.Data)

		chunk.Data = chunk.Data[:relCol]
		if chunk.Style == style {
			chunk.Data = append(chunk.Data, data...)
			newLine = append(newLine, chunk)
		} else {
			if len(chunk.Data) > 0 {
				newLine = append(newLine, chunk)
			}
			newLine = append(newLine, Chunk{Data: data, Style: style})
		}

		bytesToRemove := len(data) - originalChunkLength + relCol
		b.removeBytesInLine(bytesToRemove, i, pos.Line, &newLine)
		b.Lines[pos.Line] = newLine
		return
	}
}

func (b Buffered) insertInsideChunk(data []byte, style ansi.Style, lineNum, relCol int, chunkIndex int) {
	chunk := b.Lines[lineNum][chunkIndex]
	if chunk.Style == style {
		copy(b.Lines[lineNum][chunkIndex].Data[relCol:], data)
		return
	}
	line := b.Lines[lineNum]
	newLine := make(Line, 0, len(line)+2)
	newLine = append(newLine, line[:chunkIndex]...)
	if relCol > 0 {
		leftChunk := chunk
		leftChunk.Data = leftChunk.Data[:relCol]
		newLine = append(newLine, leftChunk)
	}
	newLine = append(newLine, Chunk{Data: data, Style: style})
	if relCol+len(data) < len(chunk.Data) {
		rightChunk := chunk
		rightChunk.Data = rightChunk.Data[relCol+len(data):]
		newLine = append(newLine, rightChunk)
	}
	newLine = append(newLine, line[chunkIndex+1:]...)
	b.Lines[lineNum] = newLine
	return
}

func (b Buffered) removeBytesInLine(bytesToRemove int, chunkIndex int, lineNum int, newLine *Line) {
	line := b.Lines[lineNum]
	for i := chunkIndex + 1; i < len(line); i++ {
		chunk := line[i]
		if bytesToRemove >= len(chunk.Data) {
			bytesToRemove -= len(chunk.Data)
			continue
		}
		if bytesToRemove > 0 {
			line[i].Data = chunk.Data[bytesToRemove:]
		}
		*newLine = append(*newLine, line[i:]...)
		return
	}
}

func (b Buffered) lineLength(i int) int {
	l := 0
	for _, chunk := range b.Lines[i] {
		l += len(chunk.Data)
	}
	return l
}

func spacer(length int) []byte {
	if length <= 0 {
		return nil
	}
	// Minor optimization: if spacer is small enough, don't need to perform a heap alloc
	if length <= spacerBytesSize {
		return spacerPreallocBytes[:length]
	}
	return bytes.Repeat([]byte{' '}, length)
}
