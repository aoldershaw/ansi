package ansi

import (
	"bytes"
	"encoding/json"
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
	Data  Text  `json:"data"`
	Style Style `json:"style"`
}

type Line = []Chunk

type Lines []Line

func (l *Lines) Print(data []byte, style Style, pos Pos) error {
	if pos.Line < 0 {
		pos.Line = 0
	}
	if pos.Col < 0 {
		pos.Col = 0
	}
	numEmpty := pos.Line - len(*l)
	for numEmpty > 0 {
		*l = append(*l, Line{})
		numEmpty--
	}
	if pos.Line >= len(*l) {
		spacerLen := pos.Col
		newData := make([]byte, spacerLen+len(data))
		copy(newData, spacer(spacerLen))
		copy(newData[spacerLen:], data)
		*l = append(*l, Line{{Data: newData, Style: style}})
		return nil
	}

	lineLen := l.lineLength(pos.Line)

	if pos.Col >= lineLen {
		l.appendToLine(data, style, pos)
	} else {
		l.insertWithinLine(data, style, pos)
	}
	return nil
}

func (l Lines) appendToLine(data []byte, style Style, pos Pos) {
	line := l[pos.Line]

	lineLen := l.lineLength(pos.Line)
	spacerLen := pos.Col - lineLen

	if len(line) == 0 {
		l.addFirstChunk(data, style, pos)
		return
	}

	lastChunk := &line[len(line)-1]
	lastChunk.Data = append(lastChunk.Data, spacer(spacerLen)...)
	if lastChunk.Style == style {
		lastChunk.Data = append(lastChunk.Data, data...)
		return
	}
	newData := make([]byte, len(data))
	copy(newData, data)
	l[pos.Line] = append(line, Chunk{Data: newData, Style: style})
}

func (l Lines) addFirstChunk(data []byte, style Style, pos Pos) {
	newData := make([]byte, pos.Col+len(data))
	copy(newData, spacer(pos.Col))
	copy(newData[pos.Col:], data)
	l[pos.Line] = Line{{Data: newData, Style: style}}
}

func (l Lines) insertWithinLine(data []byte, style Style, pos Pos) {
	line := l[pos.Line]
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
			l.insertInsideChunk(data, style, pos.Line, relCol, i)
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
			newData := make([]byte, len(data))
			copy(newData, data)
			newLine = append(newLine, Chunk{Data: newData, Style: style})
		}

		bytesToRemove := len(data) - originalChunkLength + relCol
		l.removeBytesInLine(bytesToRemove, i, pos.Line, &newLine)
		l[pos.Line] = newLine
		return
	}
}

func (l Lines) insertInsideChunk(data []byte, style Style, lineNum, relCol int, chunkIndex int) {
	chunk := l[lineNum][chunkIndex]
	if chunk.Style == style {
		copy(l[lineNum][chunkIndex].Data[relCol:], data)
		return
	}
	line := l[lineNum]
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
	l[lineNum] = newLine
	return
}

func (l Lines) removeBytesInLine(bytesToRemove int, chunkIndex int, lineNum int, newLine *Line) {
	line := l[lineNum]
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

func (l Lines) lineLength(i int) int {
	length := 0
	for _, chunk := range l[i] {
		length += len(chunk.Data)
	}
	return length
}

func (l Lines) ClearRight(pos Pos) error {
	if pos.Line < 0 || pos.Line >= len(l) {
		return nil
	}
	if pos.Col < 0 {
		pos.Col = 0
	}
	line := l[pos.Line]
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
		l[pos.Line] = line[:keepUpToChunk+1]
		return nil
	}
	return nil
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
