package cursor

import (
	"github.com/KlyuchnikovV/gapbuf"
	"github.com/KlyuchnikovV/linetree"
)

type Cursors struct {
	Lines []int
	Tree  linetree.LineTree
}

func New(tree *linetree.LineTree, lines ...int) *Cursors {
	return &Cursors{
		Lines: lines,
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
	for i, line := range cursors.Lines {
		line, err := cursors.Tree.GetLine(line)
		if err != nil {
			panic(err)
		}

		f(&line.Line, &cursors.Lines, i)
	}
}
