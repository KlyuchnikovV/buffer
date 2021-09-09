package token

import (
	"log"
	"regexp"
	"strings"
)

type Class string

type Classes []Class

func (c Classes) String() string {
	var result = make([]string, len(c))
	for i, class := range c {
		result[i] = string(class)
	}
	return strings.Join(result, "|")
}

func (c Classes) Contains(class Class) bool {
	for _, cl := range c {
		if cl == class {
			return true
		}
	}
	return false
}

const (
	Undefined    = "undefined"
	Delimeter    = "delimeter"
	Symbols      = "symbols"
	OpeningBrace = "opening-brace"
	ClosingBrace = "closing-brace"
	Quote        = "quote"
	newLine      = "new-line"
	Keyword      = "keyword"
	Type         = "built-in-type"
	BuiltInFunc  = "built-in-func"
	Comment      = "comment"
	String       = "string"
	Whitespace   = "whitespace"
)

func defineClass(char rune) Class {
	if char == '\n' {
		return newLine
	}
	if isWhitespace(char) {
		return Whitespace
	}
	if isQuote(char) {
		return Quote
	}
	if isDelimeter(char) {
		return Delimeter
	}
	if isOpeningBrace(char) {
		return OpeningBrace
	}
	if isClosingBrace(char) {
		return ClosingBrace
	}
	if isSymbol(char) {
		return Symbols
	}
	return Undefined
}

func isWhitespace(char rune) bool {
	result, err := regexp.MatchString("\\s", string(char))
	if err != nil {
		log.Panic(err)
	}
	return result
}

func isDelimeter(char rune) bool {
	result, err := regexp.MatchString("\\.|,|!|=|\\+|-|:|\\*|>|<", string(char))
	if err != nil {
		log.Panic(err)
	}
	return result
}

func isSymbol(char rune) bool {
	result, err := regexp.MatchString("[a-zA-Z/\\][a-zA-Z0-9\\/]*", string(char))
	if err != nil {
		log.Panic(err)
	}
	return result
}

func isQuote(char rune) bool {
	result, err := regexp.MatchString("\"|'|`", string(char))
	if err != nil {
		log.Panic(err)
	}
	return result
}

func isOpeningBrace(char rune) bool {
	result, err := regexp.MatchString("\\(|{|\\[", string(char))
	if err != nil {
		log.Panic(err)
	}
	return result
}

func isClosingBrace(char rune) bool {
	result, err := regexp.MatchString("\\)|}|\\]", string(char))
	if err != nil {
		log.Panic(err)
	}
	return result
}

func bracesHasSameType(left, right Token) bool {
	if left.Classes[0] != OpeningBrace {
		return false
	}
	if right.Classes[0] != ClosingBrace {
		return false
	}
	switch {
	case left.Value[0] == '(' && right.Value[0] == ')',
		left.Value[0] == '[' && right.Value[0] == ']',
		left.Value[0] == '{' && right.Value[0] == '}':
		return true
	}
	return false
}
