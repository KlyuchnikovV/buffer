package token

import (
	"fmt"
	"strings"

	"github.com/KlyuchnikovV/buffer/runes"
	"github.com/KlyuchnikovV/stack"
)

type Text struct {
	value []Line
}

func NewText(input []rune) Text {
	var (
		result = Text{
			value: make([]Line, 0),
		}
		parser      = new(SyntaxParser)
		tokens, err = parser.ParseText(input)
	)
	if err != nil {
		panic(err)
	}
	var lineTokens = make([]Token, 0)
	for _, token := range tokens {

		if !strings.ContainsRune(string(token.value), '\n') {
			lineTokens = append(lineTokens, token)
			continue
		}

		token.value = runes.ReplaceAll(token.value, []rune{'\n'}, []rune{})
		if len(token.value) != 0 {
			lineTokens = append(lineTokens, token)
		}

		result.value = append(result.value, Line{
			value: lineTokens,
		})
		lineTokens = make([]Token, 0)
	}
	return result
}

type Line struct {
	value []Token
}

func (l Line) String() string {
	var result string
	for _, token := range l.value {
		result += fmt.Sprintf("<%s[%s]>", string(token.value), token.classes)
	}
	return result
}

func NewLine(tokens []Token) Line {
	return Line{
		value: tokens,
	}
}

type SyntaxParser struct {
	contextStack *stack.Stack
}

func (p *SyntaxParser) ParseText(text []rune) ([]Token, error) {
	tokens, err := p.Tokenize(text)
	if err != nil {
		return nil, err
	}
	p.contextStack = stack.New(5)
	return p.ProcessTokens(tokens)
}

func (p *SyntaxParser) ProcessTokens(tokens []Token) ([]Token, error) {
	var (
		pushedInlineElements int
		contextIsLinear      = true
	)

	for i := 0; i < len(tokens); i++ {

		// Если до конца строки нет скобки, то контекст линейный (не пакет)

		switch tokens[i].classes[0] {
		case openingBrace:
			if i == 0 {
				continue
			}
			// FIXME: infinite recursion
			t, err := p.ProcessTokens(tokens[i:])
			if err != nil {
				return nil, err
			}
			for j := range t {
				tokens[i+j] = t[j]
			}
			i += len(t) - 1
			contextIsLinear = false
		case closingBrace:
			if !contextIsLinear {
				p.contextStack.PopN(pushedInlineElements)
				pushedInlineElements = 0
			}
			contextIsLinear = true
			// FIXME: first element must be opening brace
			if bracesHasSameType(tokens[0], tokens[i]) {
				return tokens[:i+1], nil
			}
			panic("brace")
		case symbols:
			switch string(tokens[i].value) {
			case "var", "func", "import", "for", "switch":
				p.contextStack.Push(Class(tokens[i].value))
				pushedInlineElements++
			}
		case delimeter:
		case newLine:
			p.contextStack.PopN(pushedInlineElements)
			pushedInlineElements = 0
		default:
		}

		// append classes from stack
		for _, item := range p.contextStack.ToSlice() {
			tokens[i].classes = append(tokens[i].classes, item.(Class))
		}
	}

	return tokens, nil
}

func (p *SyntaxParser) Tokenize(text []rune) ([]Token, error) {
	var (
		token = Token{
			classes: []Class{symbols},
		}
		result = make([]Token, 0)
	)
	for _, char := range text {
		var class = defineClass(char)
		if token.classes[0] == class && class != newLine {
			token.value = append(token.value, char)
			continue
		}
		if len(token.value) > 0 {
			result = append(result, token)
		}
		token = New([]rune{char}, class)
	}

	return result, nil
}
