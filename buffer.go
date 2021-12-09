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

func New(str ...byte) *Buffer {
	var (
		lines  = strings.Split(string(str), "\n")
		buffer = &Buffer{
			Cursors: *cursor.New(
				linetree.New(
					node.New(*gapbuf.New()),
				),
				1,
			),
		}
	)

	if len(str) == 0 {
		return buffer
	}

	l, err := buffer.Tree.GetLine(1)
	if err != nil {
		panic(err)
	}
	l.Line.Insert([]byte(lines[0])...)

	for i := 1; i < len(lines); i++ {
		if err := buffer.Tree.Insert(i+1, *node.New(
			*gapbuf.New([]byte(lines[i])...),
		)); err != nil {
			panic(err)
		}
	}

	return buffer
}

func (buffer *Buffer) NewLine(bytes ...byte) {
	buffer.ApplyToAll(func(gapBuffer *gapbuf.GapBuffer, lines *[]int, index int) {
		(*lines)[index]++

		if err := buffer.Tree.Insert((*lines)[index],
			*node.New(
				*gapbuf.New(append(gapBuffer.Split(), bytes...)...),
			),
		); err != nil {
			panic(err)
		}
	})
}

func (buffer *Buffer) SetCursor(line, offset int) {
	if line != -1 {
		buffer.Cursors.Lines[0] = line
	}

	if offset != -1 {
		buffer.ApplyToAll(func(gapBuffer *gapbuf.GapBuffer, lines *[]int, index int) {
			gapBuffer.SetCursor(offset)
		})
	}
}

func (buffer *Buffer) Insert(bytes ...byte) {
	buffer.ApplyToAll(func(gapBuffer *gapbuf.GapBuffer, _ *[]int, _ int) {
		gapBuffer.Insert(bytes...)
	})
}

func (buffer *Buffer) Delete() {
	buffer.ApplyToAll(func(gapBuffer *gapbuf.GapBuffer, lines *[]int, index int) {
		if gapBuffer.Offset() != 0 {
			gapBuffer.Delete()
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

		node, err := buffer.Tree.GetLine((*lines)[index])
		if err != nil {
			panic(err)
		}

		gapBuffer = &node.Line

		gapBuffer.Gap.SetCursor(gapBuffer.Size())

		gapBuffer.Insert([]byte(oldLine.String())...)
	})
}

func (buffer *Buffer) TranslateOffset(offset int) (int, int) {
	var (
		data      = buffer.Tree.String()[:offset]
		line      = strings.Count(data, "\n") + 1
		lastIndex = strings.LastIndex(data, "\n")
	)

	if lastIndex == -1 {
		return line, offset
	}

	return line, len(data[lastIndex+1:])
}
