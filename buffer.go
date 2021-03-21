package buffer

import (
	"context"
	"log"

	"github.com/KlyuchnikovV/buffer/broadcast"
	"github.com/KlyuchnikovV/buffer/runes"
	"github.com/KlyuchnikovV/edigode-cli/constants"
	"github.com/KlyuchnikovV/edigode-cli/types"
)

type Buffer struct {
	BufferTree

	line   int
	column int

	Events broadcast.Broadcast
}

func New(lines []rune) (*Buffer, error) {
	var buffer = Buffer{
		BufferTree: BufferTree{},
		Events:     *broadcast.New(context.Background()),
	}

	defer func() {
		buffer.SendLineChanged(-1)
	}()

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
	log.Printf("BUFFER: cursor down")
	buffer.CursorDown()
	buffer.column = 0
}

func (buffer *Buffer) InsertRune(r rune) {
	if r == '\n' {
		buffer.NewLine()
		return
	}
	defer func() {
		// for i := buffer.line; i < buffer.size; i++ {
		buffer.SendLineChanged(buffer.line)
		// }
	}()

	buffer.CurrentLine().Insert(r, buffer.column)
	// buffer.column++
	// buffer.CursorRight()
}

func (buffer *Buffer) DeleteRune() {
	if buffer.column == 0 {
		buffer.DeleteNewLine()
		return
	}
	defer func() {
		// for i := buffer.line; i < buffer.size; i++ {
		buffer.SendLineChanged(buffer.line)
		// }
	}()
	// buffer.column--
	buffer.CursorLeft()
	log.Printf("DELETE: column is %d", buffer.column)
	buffer.CurrentLine().Remove(buffer.column)
}

func (buffer *Buffer) NewLine() {
	var (
		line = buffer.CurrentLine()
		temp = make([]rune, line.Length())
	)
	defer func() {
		for i := buffer.line - 1; i < buffer.size; i++ {
			buffer.SendLineChanged(i)
		}
	}()

	copy(temp, line.data)

	line.data = temp[:buffer.column]

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
	defer func() {
		for i := buffer.line; i < buffer.size; i++ {
			buffer.SendLineChanged(i)
		}
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

func (buffer *Buffer) String() string {
	var (
		lines  = buffer.ToList()
		result = make([][]rune, len(lines))
	)
	for i, line := range lines {
		result[i] = (*line).Data()
	}

	return string(runes.Join(result, ' '))
}

func (buffer *Buffer) ProcessRune(r rune) error {
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
		line, column := buffer.Cursor()
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
	buffer.Events.Receiver <- types.BufferChangeMessage{
		Line: line,
	}
}

// func (buffer *Buffer) SendRender(line, column int, runes ...rune) {
// 	buffer.Events <- types.RenderMessage{
// 		Line:   line,
// 		Column: column,
// 		Runes:  runes,
// 	}
// }

func (buffer *Buffer) SendCursor(line, column int) {
	buffer.Events.Receiver <- types.CursorReposition{
		Line:   line,
		Column: column,
	}
}
