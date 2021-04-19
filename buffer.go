package buffer

import (
	"context"
	"log"

	"github.com/KlyuchnikovV/buffer/broadcast"
	"github.com/KlyuchnikovV/buffer/cursor"
	"github.com/KlyuchnikovV/buffer/messages"
	"github.com/KlyuchnikovV/buffer/runes"
	"github.com/KlyuchnikovV/edigode-cli/constants"
)

type Buffer struct {
	Name string

	BufferTree

	cursor.Cursor

	// line   int
	// column int

	Events broadcast.Broadcast
}

func New(name string, lines []rune) (*Buffer, error) {
	var buffer = Buffer{
		Name:       name,
		BufferTree: BufferTree{},
		Events:     *broadcast.New(context.Background()),
	}

	buffer.Cursor = cursor.New(0, 0, buffer.Size, buffer.getLineWidth)

	for i, item := range runes.Split(lines, '\n') {
		if err := buffer.Insert(item, i); err != nil {
			return nil, err
		}
	}
	return &buffer, nil
}

func (buffer *Buffer) CurrentLine() *Line {
	return buffer.GetNode(buffer.Line())
}

func (buffer *Buffer) getLineWidth(line int) int {
	if line < 0 || line >= buffer.size {
		log.Panicf("Wrong line parameter to getLineWidth: line: %d, size: %d", line, buffer.size)
	}
	node := buffer.GetNode(line)
	if node == nil {
		log.Panicf("Node is nil")
	}
	return len(node.data)
}

func (buffer *Buffer) GetLinesData(line, number int) [][]rune {
	var result = make([][]rune, number)
	for i := 0; i < number; i++ {
		result[i] = buffer.GetNode(line + i).Data()
	}
	return result
}

func (buffer *Buffer) InsertRune(r rune) {
	if r == '\n' {
		buffer.NewLine()
		return
	}
	defer func() {
		buffer.SendLineChanged(buffer.Line())
	}()

	buffer.CurrentLine().Insert(r, buffer.Column())
}

func (buffer *Buffer) DeleteRune() {
	if buffer.Column() == 0 {
		buffer.DeleteNewLine()
		return
	}
	defer func() {
		buffer.SendLineChanged(buffer.Line())
	}()
	buffer.CursorLeft()
	buffer.CurrentLine().Remove(buffer.Column())
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

	line.data = temp[:buffer.Column()]

	if err := buffer.Insert(temp[buffer.Column():], buffer.Line()+1); err != nil {
		log.Panic(err)
	}

	buffer.CursorRight()
}

func (buffer *Buffer) DeleteNewLine() {
	if buffer.Line() == 0 {
		return
	}
	defer func() {
		buffer.SendLineChanged(-1)
	}()

	line, _ := buffer.Delete(buffer.Line())

	buffer.CursorLeft()

	currentLine := buffer.CurrentLine()
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
		buffer.SendCursor(buffer.Line(), buffer.Column())
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
		buffer.SendCursor(buffer.Line(), buffer.Column())
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
