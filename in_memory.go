package ansi

import (
	"bytes"
	"encoding/json"
	"github.com/aoldershaw/ansi/action"
	"github.com/aoldershaw/ansi/style"
)

// Used as an optimization to avoid heap allocations
const spacerBytesSize = 512

var spacerPreallocBytes [spacerBytesSize]byte

func init() {
	for i := 0; i < spacerBytesSize; i++ {
		spacerPreallocBytes[i] = ' '
	}
}

type Text []byte

func (t Text) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(t))
}

func (t *Text) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*t = Text(s)
	return nil
}

type Chunk struct {
	Data  Text        `json:"data"`
	Style style.Style `json:"style"`
}

type Line = []Chunk

type InMemory struct {
	Lines []Line
}

func (b *InMemory) Print(data []byte, style style.Style, pos action.Pos) {
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
	if pos.Line >= len(b.Lines) {
		spacerLen := pos.Col
		if spacerLen > 0 {
			newData := make([]byte, spacerLen+len(data))
			copy(newData, spacer(spacerLen))
			copy(newData[spacerLen:], data)
			data = newData
		}
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

func (b *InMemory) appendToLine(data []byte, style style.Style, pos action.Pos) {
	line := b.Lines[pos.Line]

	lineLen := b.lineLength(pos.Line)
	spacerLen := pos.Col - lineLen

	if len(line) == 0 {
		b.addFirstChunk(data, style, pos)
		return
	}

	lastChunk := &line[len(line)-1]
	lastChunk.Data = append(lastChunk.Data, spacer(spacerLen)...)
	if lastChunk.Style == style {
		lastChunk.Data = append(lastChunk.Data, data...)
		return
	}
	b.Lines[pos.Line] = append(line, Chunk{Data: data, Style: style})
}

func (b *InMemory) addFirstChunk(data []byte, style style.Style, pos action.Pos) {
	if pos.Col > 0 {
		newData := make([]byte, pos.Col+len(data))
		copy(newData, spacer(pos.Col))
		copy(newData[pos.Col:], data)
		data = newData
	}
	b.Lines[pos.Line] = Line{{Data: data, Style: style}}
}

func (b *InMemory) insertWithinLine(data []byte, style style.Style, pos action.Pos) {
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

func (b InMemory) insertInsideChunk(data []byte, style style.Style, lineNum, relCol int, chunkIndex int) {
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

func (b InMemory) removeBytesInLine(bytesToRemove int, chunkIndex int, lineNum int, newLine *Line) {
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

func (b InMemory) lineLength(i int) int {
	l := 0
	for _, chunk := range b.Lines[i] {
		l += len(chunk.Data)
	}
	return l
}

func (b *InMemory) ClearRight(pos action.Pos) {
	if pos.Line < 0 || pos.Line >= len(b.Lines) {
		return
	}
	if pos.Col < 0 {
		pos.Col = 0
	}
	line := b.Lines[pos.Line]
	chunkEnd := 0
	for i := 0; i < len(line); i++ {
		chunk := &line[i]
		chunkStart := chunkEnd
		chunkEnd += len(chunk.Data)
		if chunkEnd < pos.Col {
			continue
		}
		chunk.Data = chunk.Data[:pos.Col-chunkStart]
		keepUpToChunk := i
		if len(chunk.Data) == 0 {
			keepUpToChunk--
		}
		b.Lines[pos.Line] = line[:keepUpToChunk+1]
		return
	}
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
