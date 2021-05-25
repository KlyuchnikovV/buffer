package token

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	var text = `package main

import (
	"fmt"
)

func init() {

}

func main() {
	fmt.Printf("Hello, world!")
}

	`
	parser := SyntaxParser{}
	tokenized, err := parser.ParseText([]rune(text))
	assert.NoError(t, err)
	for i, token := range tokenized {
		fmt.Printf("%d - `%s`, classes: %v\n", i+1, strings.ReplaceAll(string(token.Value), "\n", "<new-line>"), token.Classes)
	}
	t.FailNow()
}

// func TestLine(t *testing.T) {
// 	line := NewLine([]rune("fmt.Printf(\"Hello, world!\")"))
// 	for i, token := range line.value {
// 		fmt.Printf("%d - `%s`, classes: %v\n", i+1, strings.ReplaceAll(string(token.value), "\n", "<new-line>"), token.classes)
// 	}
// 	t.FailNow()
// }

func TestText(t *testing.T) {
	var input = `package main

import (
	"fmt"
)

type Editor struct {
	field int
}

func init() {
	// TODO: fallback if file not found etc
	if b, err := editor.LoadFile(path); err != nil {
		return nil, err
	} else {
		textareas = append(textareas, textarea.New(resizing.NewPercents(1), b))
	}
}

// hello variable
var hello = "Hello, world!"

func main() {
	fmt.Printf(hello)
}

	`
	text := NewText("", []rune(input))
	for i, line := range text.Value {
		fmt.Printf("%d - `%s`\n", i+1, line.debugString())
	}
	t.FailNow()
}
