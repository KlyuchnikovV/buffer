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
	"github.com/KlyuchnikovV/stack"
)

type Text struct {
	Name string

	Value []Line

	cursor.Cursor

	Events broadcast.Broadcast
	*SyntaxParser
}

func NewText(name string, input []rune) *Text {
	var (
		result = &Text{
			Name:         name,
			Value:        make([]Line, 0),
			Events:       *broadcast.New(context.Background()),
			SyntaxParser: new(SyntaxParser),
		}
		tokens, err = result.SyntaxParser.ParseText(input)
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
	defer func() {
		text.SendLineChanged(text.Line())
	}()
	var (
		lineNum   = text.Line()
		line      = text.Value[lineNum]
		err       error
		prevToken *Token
	)

	if len(line.Value) > 0 {
		prevToken = &line.Value[text.tokenByPosition(lineNum, text.Column())]
	}

	switch r {
	case constants.BackspaceRune:
		if text.Column() == 0 {
			defer text.SendLineChanged(-1)
			text.DeleteNewLine()
			lineNum = text.Line()
		} else {
			text.DeleteRune()
		}
		text.SendCursor(text.Line(), text.Column())
	case constants.EnterRune:
		text.NewLine()
		text.SendCursor(text.Line(), text.Column())
	default:
		text.InsertRune(r)
		text.CursorRight()
		text.SendCursor(text.Line(), text.Column())
	}

	text.reParseText(lineNum, prevToken)

	return err
}

func (text *Text) NewLine() {
	var (
		lineNum            = text.Line()
		prevLine, nextLine = text.Value[lineNum].Split(text.Column())
		temp               = make([]Line, len(text.Value))
	)

	copy(temp, text.Value)

	text.Value = append(text.Value[:lineNum], prevLine, nextLine)
	text.Value = append(text.Value, temp[lineNum+1:]...)

	text.SendLineChanged(-1)
	text.CursorRight()
}

func (text *Text) DeleteNewLine() {
	var lineNum = text.Line()
	if lineNum == 0 {
		return
	}

	var prevLine = text.CurrentLine()

	text.CursorLeft()

	text.Value[lineNum-1].Append(prevLine.Value...)
	text.Value = append(text.Value[:lineNum], text.Value[lineNum+1:]...)
}

func (text *Text) InsertRune(r rune) {
	var (
		line   = text.CurrentLine()
		column = text.Column()
	)

	log.Printf("INSERT: inserting into line %d position %d rune %c", text.Line(), text.Column(), r)
	line.Insert(column, r)

	log.Printf("INSERT: reparsing...")
}

func (text Text) tokenByPosition(line, position int) int {
	index, _ := text.Value[line].TokenByPosition(position)
	return index
}

func (text *Text) DeleteRune() {
	var column = text.Column()
	if column == 0 {
		text.DeleteNewLine()
		return
	}

	text.CurrentLine().Remove(column)
	text.CursorLeft()
}

func (text *Text) ProcessEscape(sequence []rune) error {
	switch {
	case runes.Equal(constants.ArrowUp, sequence):
		text.CursorUp()
		text.SendCursor(text.Line(), text.Column())
	case runes.Equal(constants.ArrowDown, sequence):
		text.CursorDown()
		text.SendCursor(text.Line(), text.Column())
	case runes.Equal(constants.ArrowLeft, sequence):
		text.CursorLeft()
		text.SendCursor(text.Line(), text.Column())
	case runes.Equal(constants.ArrowRight, sequence):
		text.CursorRight()
		text.SendCursor(text.Line(), text.Column())
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

func (text *Text) reParseText(line int, prevToken *Token) {

	var (
		column        = text.Column()
		err           error
		context       = stack.New(0)
		startingToken = text.tokenByPosition(line, column)
		prevClass     Class
	)

	if prevToken != nil {
		prevClass = Class(prevToken.Value)
	}

	if startingToken > 0 {
		startingToken -= 1
	}

	if line != 0 && startingToken != 0 {
		var classes = text.Value[line].Value[0].Classes
		for i := 0; i < len(classes)-1; i++ {
			if prevClass != "" && classes[i] == prevClass {
				break
			}
			context.Push(classes[i])
		}
	}

	for i := line; i < text.Size(); i++ {
		var tokens = text.Value[i].Value
		if len(tokens) == 0 {
			continue
		}

		// log.Printf("REPARSE: line %d was %s", i, text.Value[i].debugString())
		text.Value[i].Value, err = text.ProcessTokens(tokens, *context)
		// log.Printf("REPARSE: line %d now %s", i, text.Value[i].debugString())
		if err != nil {
			panic(err)
		}
	}
}
