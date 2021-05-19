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

const (
	undefined    = "undefined"
	delimeter    = "delimeter"
	symbols      = "symbols"
	openingBrace = "opening-brace"
	closingBrace = "closing-brace"
	quote        = "quote"
	newLine      = "new-line"
)

func defineClass(char rune) Class {
	if char == '\n' {
		return newLine
	}
	if isQuote(char) {
		return quote
	}
	if isDelimeter(char) {
		return delimeter
	}
	if isOpeningBrace(char) {
		return openingBrace
	}
	if isClosingBrace(char) {
		return closingBrace
	}
	if isSymbol(char) {
		return symbols
	}
	return undefined
}

func isDelimeter(char rune) bool {
	result, err := regexp.MatchString("\\.|,|\\s|!", string(char))
	if err != nil {
		log.Panic(err)
	}
	return result
}

func isSymbol(char rune) bool {
	result, err := regexp.MatchString("[a-zA-Z][a-zA-Z0-9]*", string(char))
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
	if left.classes[0] != openingBrace {
		return false
	}
	if right.classes[0] != closingBrace {
		return false
	}
	switch {
	case left.value[0] == '(' && right.value[0] == ')',
		left.value[0] == '[' && right.value[0] == ']',
		left.value[0] == '{' && right.value[0] == '}':
		return true
	}
	return false
}
