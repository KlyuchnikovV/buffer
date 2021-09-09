package token

import (
	"fmt"
	"log"
)

type Line struct {
	Value []Token
}

func NewLine(tokens []Token) Line {
	return Line{
		Value: tokens,
	}
}

func (l Line) String() string {
	return string(l.Data())
}

func (l Line) Data() []rune {
	var result = make([]rune, 0)
	for _, token := range l.Value {
		result = append(result, token.Value...)
	}
	return result
}

func (l Line) debugString() string {
	var result string
	for _, token := range l.Value {
		result += fmt.Sprintf("<%s[%s]>", string(token.Value), token.Classes)
	}
	return result
}

func (l Line) Length() int {
	var result int
	for _, token := range l.Value {
		result += len(token.Value)
	}
	return result
}

func (l *Line) appendTokens(tokenIndex, position int, tokens ...Token) {
	var (
		result = make([]Token, 0, len(l.Value)+len(tokens)+1)
	)

	result = append(result, l.Value[:tokenIndex]...)

	var prev, next = l.Value[tokenIndex].Split(position)
	result = append(result, prev)
	result = append(result, tokens...)
	result = append(result, next)

	result = append(result, l.Value[tokenIndex+1:]...)

	l.Value = result
}

func (l *Line) Insert(position int, r rune) {
	defer func() {

	}()
	var (
		index    int
		runeType = defineClass(r)
	)

	index, position = l.TokenByPosition(position)

	log.Printf("INSERT: chosen token %d with value '%s'", index, string(l.Value[index].Value))
	if runeType == l.Value[index].Classes[0] {
		l.Value[index].Insert(position, r)
		return
	}
	var classes = append([]Class{runeType}, l.Value[index].Classes[1:]...)
	l.appendTokens(index, position, New([]rune{r}, classes...))
}

func (l *Line) Remove(position int) {
	index, position := l.TokenByPosition(position)
	if index == 0 && position == 0 {
		return
	}
	if position == 0 {
		index--
		position = l.Value[index].Length()
	}

	position--

	l.Value[index].Remove(position)

}

func (l *Line) TokenByPosition(position int) (int, int) {
	for i, token := range l.Value {
		if position >= 0 && position <= token.Length() {
			return i, position
		}
		position -= token.Length()
	}
	return 0, 0
}

func (l *Line) Append(tokens ...Token) {
	if l.Value[len(l.Value)-1].Classes[0] == tokens[0].Classes[0] {
		l.Value[len(l.Value)-1].Value = append(l.Value[len(l.Value)-1].Value, tokens[0].Value...)
		tokens = tokens[1:]
	}
	l.Value = append(l.Value, tokens...)
}

func (l Line) Split(position int) (Line, Line) {
	if len(l.Value) == 0 {
		return NewLine(nil), NewLine(nil)
	}

	index, position := l.TokenByPosition(position)
	prevToken, nextToken := l.Value[index].Split(position)

	var temp = make([]Token, len(l.Value))
	copy(temp, l.Value)

	return NewLine(append(temp[:index], prevToken)), NewLine(append([]Token{nextToken}, temp[index+1:]...))
}
