package cursor

import (
	"github.com/KlyuchnikovV/gapbuf"
	"github.com/KlyuchnikovV/linetree"
)

type Cursors struct {
	lines []int
	Tree  linetree.LineTree
}

func New(tree *linetree.LineTree, lines ...int) *Cursors {
	return &Cursors{
		lines: lines,
		Tree:  *tree,
	}
}

// func NewCursors(line, column int, height func() int, lineWidth func(int) int) Cursors {

// }

func (cursors *Cursors) AddCursor(line, column int) {

}

func (cursors *Cursors) RemoveCursor(index int) {

}

func (cursors *Cursors) ApplyToAll(f func(gapBuffer *gapbuf.GapBuffer, lines *[]int, index int)) {
	for i, line := range cursors.lines {
		line, err := cursors.Tree.GetLine(line)
		if err != nil {
			panic(err)
		}

		gapBuf, ok := line.Line.(*gapbuf.GapBuffer)
		if !ok {
			panic("!ok")
		}

		f(gapBuf, &cursors.lines, i)
	}
}
