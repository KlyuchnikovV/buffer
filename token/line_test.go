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
		fmt.Printf("%d - `%s`, classes: %v\n", i+1, strings.ReplaceAll(string(token.value), "\n", "<new-line>"), token.classes)
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

func init() {

}

func main() {
	fmt.Printf("Hello, world!")
}

	`
	text := NewText([]rune(input))
	for i, line := range text.value {
		fmt.Printf("%d - `%s`\n", i, line)
	}
	t.FailNow()
}
