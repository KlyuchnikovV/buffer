package token

import (
	"fmt"

	"github.com/KlyuchnikovV/stack"
)

type Line struct {
	Value []Token
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

func NewLine(tokens []Token) Line {
	return Line{
		Value: tokens,
	}
}

type SyntaxParser struct {
}

func (p *SyntaxParser) ParseText(text []rune) ([]Token, error) {
	tokens, err := p.Tokenize(text)
	if err != nil {
		return nil, err
	}

	st := stack.New(5)
	return p.ProcessTokens(tokens, *st)
}

func (p *SyntaxParser) ProcessTokens(tokens []Token, context stack.Stack) ([]Token, error) {
	var immutableContextPart = context.Size()
	for i := 0; i < len(tokens); i++ {
		switch tokens[i].Classes[0] {
		case Symbols:
			switch string(tokens[i].Value) {
			case "//":
				context.Push(Class(Comment))
				to := p.FindElementPosition(tokens[i:], "\n")
				for j := 0; j < to-1; j++ {
					for _, item := range context.ToSlice() {
						tokens[i].Classes = append(tokens[i].Classes, item.(Class))
					}
					i++
				}
				// continue
			case "package", "func", "var", "import", "for", "switch", "type", "return", "if", "else", "case":
				tokens[i].Classes = append(tokens[i].Classes, Keyword)
				context.Push(Class(tokens[i].Value))
			case "int", "bool", "string", "struct", "interface", "error", "nil":
				tokens[i].Classes = append(tokens[i].Classes, Type)
			case "append", "len", "cap", "panic", "go", "defer":
				tokens[i].Classes = append(tokens[i].Classes, BuiltInFunc)
				context.Push(Class(tokens[i].Value))
			}
		case newLine:
			if context.Size() > immutableContextPart {
				context.Pop()
			}
		case Quote:
			context.Push(Class(String))
			t, err := p.ProcessTokens(p.PairQuote(tokens[i:]), context)
			if err != nil {
				return nil, err
			}
			for j := range t {
				tokens[i] = t[j]
				i++
			}
			context.Pop()
			continue
		case OpeningBrace:
			t, err := p.ProcessTokens(p.PairBrace(tokens[i:]), context)
			if err != nil {
				return nil, err
			}
			for j := range t {
				tokens[i] = t[j]
				i++
			}
			continue
		}

		// append classes from stack
		for _, item := range context.ToSlice() {
			tokens[i].Classes = append(tokens[i].Classes, item.(Class))
		}
	}

	return tokens, nil
}

func (p *SyntaxParser) FindElementPosition(tokens []Token, value string) int {
	for i, token := range tokens {
		if string(token.Value) == value {
			return i
		}
	}
	return -1
}

func (p *SyntaxParser) PairBrace(tokens []Token) []Token {
	var braceStack = stack.New(1)
	for i, token := range tokens {
		switch token.Classes[0] {
		case OpeningBrace:
			braceStack.Push(token)
		case ClosingBrace:
			t, ok := braceStack.Peek()
			if !ok {
				panic("brace err")
			}
			if bracesHasSameType(t.(Token), token) {
				braceStack.Pop()
			}
		}
		if braceStack.Size() == 0 {
			return tokens[:i]
		}
	}
	return nil
}

func (p *SyntaxParser) PairQuote(tokens []Token) []Token {
	var quoteStack = stack.New(1)
	for i, token := range tokens {
		if token.Classes[0] != Quote {
			continue
		}
		if i == 0 {
			quoteStack.Push(token)
			continue
		}

		t, ok := quoteStack.Peek()
		if !ok {
			panic("quote err")
		}
		if string(t.(Token).Value) == string(token.Value) {
			quoteStack.Pop()
		} else {
			quoteStack.Push(token)
		}

		if quoteStack.Size() == 0 {
			return tokens[:i]
		}
	}
	return nil
}

func (p *SyntaxParser) Tokenize(text []rune) ([]Token, error) {
	var (
		token = Token{
			Classes: []Class{Symbols},
		}
		result = make([]Token, 0)
	)
	for _, char := range text {
		var class = defineClass(char)
		if token.Classes[0] == class && class != newLine {
			token.Value = append(token.Value, char)
			continue
		}
		if len(token.Value) > 0 {
			result = append(result, token)
		}
		token = New([]rune{char}, class)
	}

	return result, nil
}
