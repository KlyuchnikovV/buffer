package buffer

import (
	"log"

	"github.com/KlyuchnikovV/buffer/line"
)

type Buffer struct {
	BufferTree

	line   int
	column int
}

func New(lines []line.Line) (*Buffer, error) {
	var buffer = Buffer{
		BufferTree: BufferTree{},
	}

	for i, item := range lines {
		if err := buffer.Insert(item, i); err != nil {
			return nil, err
		}
	}
	return &buffer, nil
}

func (buffer *Buffer) CurrentLine() *line.Line {
	line, _ := buffer.Find(buffer.line)
	return line
}

func (buffer *Buffer) NewLine() {
	var (
		line = buffer.CurrentLine()
		temp = line.Data()
	)

	(*line) = temp[:buffer.column]

	buffer.line++

	if err := buffer.Insert(temp[buffer.column:], buffer.line); err != nil {
		log.Print(err)
		panic(err)
	}

	buffer.column = 0
}

func (buffer *Buffer) DeleteNewLine() {
	if buffer.line == 0 {
		return
	}

	line, _ := buffer.Delete(buffer.line)
	buffer.line--
	currentLine := buffer.CurrentLine()
	buffer.column = line.Length()

	(*currentLine) = append((*currentLine), line.Data()...)
}

func (buffer *Buffer) Cursor() (int, int) {
	return buffer.line, buffer.column
}

func (buffer *Buffer) CursorUp() {
	if buffer.line == 0 {
		buffer.column = 0
		return
	}

	buffer.line--

	var lenOfLine = buffer.CurrentLine().Length()
	if lenOfLine < buffer.column {
		buffer.column = lenOfLine
	}
}

func (buffer *Buffer) CursorDown() {
	line, ok := buffer.Find(buffer.Size() - 1)
	if !ok {
		panic("line not found")
	}
	if buffer.line == buffer.Size()-1 {
		buffer.column = line.Length()
		return
	}

	buffer.line++

	var lenOfLine = buffer.CurrentLine().Length()
	if lenOfLine < buffer.column {
		buffer.column = lenOfLine
	}
}

func (buffer *Buffer) CursorLeft() {
	if buffer.column > 0 {
		buffer.column--
		return
	}

	if buffer.line == 0 {
		return
	}
	buffer.CursorUp()
	buffer.column = buffer.CurrentLine().Length()
}

func (buffer *Buffer) CursorRight() {
	if buffer.column < buffer.CurrentLine().Length() {
		buffer.column++
		return
	}

	if buffer.line == buffer.Size()-1 {
		return
	}
	buffer.CursorDown()
	buffer.column = 0
}

func (buffer *Buffer) FindAll(runes []rune) [][]int {
	var result = make([][]int, 0)
	for _, line := range buffer.ToList() {
		result = append(result, line.FindAll(runes))
	}
	return result
}

func (buffer *Buffer) FindNext(prevLine, prevColumn int, runes []rune) (int, int) {
	var temp = buffer.ToList()
	if prevLine != -1 {
		temp = temp[prevLine:]
	}

	if len(temp) > 0 && prevColumn != -1 {
		column, ok := temp[0].FindNext(prevColumn, runes)
		if ok {
			return prevLine, column
		}
		temp = temp[1:]
		prevLine++
	}

	for _, line := range temp {
		column, ok := line.FindNext(0, runes)
		if ok {
			return prevLine, column
		}
		prevLine++
	}
	return -1, -1
}
