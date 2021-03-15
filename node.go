package buffer

import (
	"fmt"
	"math"
	"strings"

	"github.com/KlyuchnikovV/buffer/line"
)

type Node struct {
	data                  *line.Line
	numberOfChildren      int
	leftChild, rightChild *Node
}

func newNode(data line.Line) *Node {
	return &Node{
		data: &data,
	}
}

func (n *Node) SetData(data *line.Line) {
	n.data = data
}

func (n Node) Data() *line.Line {
	return n.data
}

func (n Node) LeftChild() *Node {
	return n.leftChild
}

func (n Node) RightChild() *Node {
	return n.rightChild
}

func (n Node) HasLeft() bool {
	return n.leftChild != nil
}

func (n Node) HasRight() bool {
	return n.rightChild != nil
}

func (n *Node) IsLeftOf(parent Node) bool {
	return parent.leftChild == n
}

func (n *Node) IsRightOf(parent Node) bool {
	return parent.rightChild == n
}

func (n Node) IsLeaf() bool {
	return n.leftChild == nil && n.rightChild == nil
}

func (n Node) getPosition(prevPosition int, goingLeft bool) int {
	if goingLeft {
		if n.rightChild != nil {
			return prevPosition - n.rightChild.numberOfChildren - 2
		}
		return prevPosition - 1
	}
	if n.leftChild != nil {
		return prevPosition + n.leftChild.numberOfChildren + 2
	}
	return prevPosition + 1
}

func (n Node) height() int {
	return 1 + log2(n.numberOfChildren+1)
}

func (n Node) getBalance() int {
	result := 0
	if n.leftChild != nil {
		result -= n.leftChild.height()
	}
	if n.rightChild != nil {
		result += n.rightChild.height()
	}
	return result
}

func (n *Node) rotateRight() {
	if n.leftChild == nil {
		return
	}
	temp := *n

	*n = *n.leftChild
	temp.leftChild = (*n).rightChild
	(*n).rightChild = &temp

	n.rightChild.fixNumberOfChildren()
	n.fixNumberOfChildren()
}

func (n *Node) rotateLeft() {
	if n.rightChild == nil {
		return
	}
	temp := *n

	*n = *n.rightChild
	temp.rightChild = (*n).leftChild
	(*n).leftChild = &temp

	n.leftChild.fixNumberOfChildren()
	n.fixNumberOfChildren()
}

func (n *Node) fixNumberOfChildren() {
	if n.leftChild == nil {
		(*n).numberOfChildren = 0
	} else {
		(*n).numberOfChildren = n.leftChild.numberOfChildren + 1
	}

	if n.rightChild != nil {
		(*n).numberOfChildren += n.rightChild.numberOfChildren + 1
	}
}

func (n *Node) balance() {
	switch n.getBalance() {
	case 2:
		if n.rightChild.getBalance() < 0 {
			n.rightChild.rotateRight()
		}
		n.rotateLeft()
	case -2:
		if n.leftChild.getBalance() > 0 {
			n.leftChild.rotateLeft()
		}
		n.rotateRight()
	}
}

func (n Node) toList() []*line.Line {
	var result = make([]*line.Line, 0, n.numberOfChildren+1)

	if n.leftChild != nil {
		result = n.leftChild.toList()
	}
	result = append(result, n.data)
	if n.rightChild != nil {
		result = append(result, n.rightChild.toList()...)
	}
	return result
}

func (n Node) visualizeNodeSubtree(currentLevel, treeHeight int) {
	if n.leftChild != nil {
		n.leftChild.visualizeNodeSubtree(currentLevel+1, treeHeight)
	}

	fmt.Printf("%s%#v%s\n", strings.Repeat("  ", currentLevel), n.data, strings.Repeat("--", treeHeight-currentLevel))

	if n.rightChild != nil {
		n.rightChild.visualizeNodeSubtree(currentLevel+1, treeHeight)
	}
}

func log2(n int) int {
	if n < 1 {
		return math.MinInt64
	}

	count := 0
	for ; n >= 2; n >>= 1 {
		count++
	}
	return count
}
