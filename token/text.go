package token

import (
	"context"
	"log"
	"strings"

	"github.com/KlyuchnikovV/buffer/broadcast"
	"github.com/KlyuchnikovV/buffer/cursor"
	"github.com/KlyuchnikovV/buffer/messages"
	"github.com/KlyuchnikovV/buffer/runes"
	"github.com/KlyuchnikovV/edigode-cli/constants"
)

type Text struct {
	Name string

	Value []Line

	cursor.Cursor

	Events broadcast.Broadcast
}

func NewText(name string, input []rune) *Text {
	var (
		result = &Text{
			Name:   name,
			Value:  make([]Line, 0),
			Events: *broadcast.New(context.Background()),
		}
		parser      = new(SyntaxParser)
		tokens, err = parser.ParseText(input)
	)
	result.Cursor = cursor.New(0, 0, result.Size, result.getLineWidth)
	if err != nil {
		panic(err)
	}
	var lineTokens = make([]Token, 0)
	for _, token := range tokens {
		if !strings.ContainsRune(string(token.Value), '\n') {
			lineTokens = append(lineTokens, token)
			continue
		}

		token.Value = runes.ReplaceAll(token.Value, []rune{'\n'}, []rune{})
		if len(token.Value) != 0 {
			lineTokens = append(lineTokens, token)
		}

		result.Value = append(result.Value, Line{
			Value: lineTokens,
		})
		lineTokens = make([]Token, 0)
	}

	if len(lineTokens) > 0 {
		result.Value = append(result.Value, Line{
			Value: lineTokens,
		})
	}

	return result
}

func (text *Text) getLineWidth(line int) int {
	if line < 0 || line >= text.Size() {
		log.Panicf("Wrong line parameter to getLineWidth: line: %d, size: %d", line, text.Size())
	}
	return text.Value[line].Length()
}

func (text *Text) CurrentLine() *Line {
	return text.GetLine(text.Line())
}

func (text *Text) GetLine(line int) *Line {
	if line < 0 || line >= text.Size() {
		return nil
	}
	return &text.Value[line]
}

func (text *Text) Size() int {
	return len(text.Value)
}

func (text Text) String() string {
	var result string
	for _, line := range text.Value {
		result += line.String()
	}
	return result
}

func (text *Text) ProcessRune(r rune) error {
	// defer func() {
	// 	buffer.SendCursor(buffer.Line(), buffer.Column())
	// }()
	// var (
	// 	err error
	// )
	// switch r {
	// case constants.BackspaceRune:
	// 	buffer.DeleteRune()
	// case constants.EnterRune:
	// 	buffer.NewLine()
	// default:
	// 	buffer.InsertRune(r)
	// 	buffer.CursorRight()
	// }
	return nil //err
}

func (text *Text) ProcessEscape(sequence []rune) error {
	defer func() {
		text.SendCursor(text.Line(), text.Column())
	}()
	switch {
	case runes.Equal(constants.ArrowUp, sequence):
		text.CursorUp()
	case runes.Equal(constants.ArrowDown, sequence):
		text.CursorDown()
	case runes.Equal(constants.ArrowLeft, sequence):
		text.CursorLeft()
	case runes.Equal(constants.ArrowRight, sequence):
		text.CursorRight()
	}
	return nil
}

func (text *Text) SendLineChanged(line int) {
	text.Events.Receiver <- messages.BufferChange{
		Source: text.Name,
		Row:    line,
	}
}

func (text *Text) SendCursor(line, column int) {
	text.Events.Receiver <- messages.CursorReposition{
		Source: text.Name,
		Row:    line,
		Column: column,
	}
}
