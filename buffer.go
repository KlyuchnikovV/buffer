package buffer

import (
	"context"
	"log"

	"github.com/KlyuchnikovV/buffer/broadcast"
	"github.com/KlyuchnikovV/buffer/messages"
	"github.com/KlyuchnikovV/buffer/runes"
	"github.com/KlyuchnikovV/edigode-cli/constants"
)

type Buffer struct {
	Name string

	BufferTree

	line   int
	column int

	Events broadcast.Broadcast
}

func New(name string, lines []rune) (*Buffer, error) {
	var buffer = Buffer{
		Name:       name,
		BufferTree: BufferTree{},
		Events:     *broadcast.New(context.Background()),
	}

	for i, item := range runes.Split(lines, '\n') {
		if err := buffer.Insert(item, i); err != nil {
			return nil, err
		}
	}
	return &buffer, nil
}

func (buffer *Buffer) CurrentLine() *Line {
	return buffer.GetNode(buffer.line)
}

func (buffer *Buffer) GetLinesData(line, number int) [][]rune {
	var result = make([][]rune, number)
	for i := 0; i < number; i++ {
		result[i] = buffer.GetNode(line + i).Data()
	}
	return result
}

func (buffer *Buffer) Line() int {
	return buffer.line
}

func (buffer *Buffer) Column() int {
	return buffer.column
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

func (buffer *Buffer) InsertRune(r rune) {
	if r == '\n' {
		buffer.NewLine()
		return
	}
	defer func() {
		buffer.SendLineChanged(buffer.line)
	}()

	buffer.CurrentLine().Insert(r, buffer.column)
}

func (buffer *Buffer) DeleteRune() {
	if buffer.column == 0 {
		buffer.DeleteNewLine()
		return
	}
	defer func() {
		buffer.SendLineChanged(buffer.line)
	}()
	buffer.CursorLeft()
	buffer.CurrentLine().Remove(buffer.column)
}

func (buffer *Buffer) NewLine() {
	var (
		line = buffer.CurrentLine()
		temp = make([]rune, line.Length())
	)
	defer func() {
		buffer.SendLineChanged(-1)
	}()

	copy(temp, line.data)

	line.data = temp[:buffer.column]

	buffer.line++

	if err := buffer.Insert(temp[buffer.column:], buffer.line); err != nil {
		log.Panic(err)
	}

	buffer.column = 0
}

func (buffer *Buffer) DeleteNewLine() {
	if buffer.line == 0 {
		return
	}
	defer func() {
		buffer.SendLineChanged(-1)
	}()

	line, _ := buffer.Delete(buffer.line)
	buffer.line--
	currentLine := buffer.CurrentLine()
	buffer.column = currentLine.Length()

	currentLine.AppendData(line, currentLine.Length())
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
	return -0, -1
}

func (buffer Buffer) String() string {
	var (
		lines  = buffer.ToList()
		result = make([][]rune, len(lines))
	)
	for i, line := range lines {
		result[i] = (*line).Data()
	}

	return string(runes.Join(result, '\n'))
}

func (buffer *Buffer) ProcessRune(r rune) error {
	defer func() {
		line, column := buffer.line, buffer.column
		buffer.SendCursor(line, column)
	}()
	var (
		err error
	)
	switch r {
	case constants.BackspaceRune:
		buffer.DeleteRune()
	case constants.EnterRune:
		buffer.NewLine()
	default:
		buffer.InsertRune(r)
		buffer.CursorRight()
	}
	return err
}

func (buffer *Buffer) ProcessEscape(sequence []rune) error {
	defer func() {
		line, column := buffer.line, buffer.column
		buffer.SendCursor(line, column)
	}()
	switch {
	case runes.Equal(constants.ArrowUp, sequence):
		buffer.CursorUp()
	case runes.Equal(constants.ArrowDown, sequence):
		buffer.CursorDown()
	case runes.Equal(constants.ArrowLeft, sequence):
		buffer.CursorLeft()
	case runes.Equal(constants.ArrowRight, sequence):
		buffer.CursorRight()
	}
	return nil
}

func (buffer *Buffer) SendLineChanged(line int) {
	buffer.Events.Receiver <- messages.BufferChange{
		Source: buffer.Name,
		Row:    line,
	}
}

func (buffer *Buffer) SendCursor(line, column int) {
	buffer.Events.Receiver <- messages.CursorReposition{
		Source: buffer.Name,
		Row:    line,
		Column: column,
	}
}
