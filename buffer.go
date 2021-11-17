package buffer

import (
	"strings"

	"github.com/KlyuchnikovV/buffer/cursor"
	"github.com/KlyuchnikovV/gapbuf"
	"github.com/KlyuchnikovV/linetree"
	"github.com/KlyuchnikovV/linetree/node"
)

type Buffer struct {
	cursor.Cursors
}

func New(str string) *Buffer {
	var (
		lines  = strings.Split(str, "\n")
		buffer = &Buffer{
			Cursors: *cursor.New(linetree.New(nil), 0),
		}
	)

	for _, line := range lines {
		buffer.NewLine(line)
	}

	return buffer
}

func (buffer *Buffer) NewLine(str string) {
	buffer.ApplyToAll(func(gapBuffer *gapbuf.GapBuffer, lines *[]int, index int) {
		(*lines)[index]++

		buffer.Tree.Insert((*lines)[index],
			*node.New(
				gapbuf.NewFromString(string(gapBuffer.Split()) + str),
			),
		)
	})
}

func (buffer *Buffer) Insert(str string) {
	buffer.ApplyToAll(func(gapBuffer *gapbuf.GapBuffer, _ *[]int, _ int) {
		gapBuffer.Insert([]byte(str)...)
	})
}

func (buffer *Buffer) Delete() {
	buffer.ApplyToAll(func(gapBuffer *gapbuf.GapBuffer, lines *[]int, index int) {
		if gapBuffer.GetCursor() != 0 {
			gapBuffer.Delete(1)
			return
		}
		if (*lines)[index] == 0 {
			return
		}

		oldLine, err := buffer.Tree.Remove((*lines)[index])
		if err != nil {
			panic(err)
		}
		(*lines)[index]--

		line, err := buffer.Tree.GetLine((*lines)[index])
		if err != nil {
			panic(err)
		}

		gapBuf, ok := line.Line.(*gapbuf.GapBuffer)
		if !ok {
			panic("!ok")
		}

		gapBuf.MoveGap(gapBuf.Size() - 1)

		gapBuf.Insert([]byte(oldLine.String())...)
	})
}
