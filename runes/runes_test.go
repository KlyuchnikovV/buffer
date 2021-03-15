package runes

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	file, err := os.Open("../../buffer/buffer.go")
	assert.NoError(t, err)

	bytes, err := ioutil.ReadAll(file)
	assert.NoError(t, err)

	r := Split([]rune(string(bytes)), '\n')

	result := make([]string, len(r))
	for i := range r {
		result[i] = string(r[i])
	}

	assert.Equal(t, strings.Split(string(bytes), "\n"), result)
}
