package buffer

import (
	"math"
)

type Line struct {
	data                  []rune
	numberOfChildren      int
	leftChild, rightChild *Line
}

func (n Line) LeftChild() *Line {
	return n.leftChild
}

func (n Line) RightChild() *Line {
	return n.rightChild
}

func (n Line) HasLeft() bool {
	return n.leftChild != nil
}

func (n Line) HasRight() bool {
	return n.rightChild != nil
}

func (n *Line) IsLeftOf(parent Line) bool {
	return parent.leftChild == n
}

func (n *Line) IsRightOf(parent Line) bool {
	return parent.rightChild == n
}

func (n Line) IsLeaf() bool {
	return n.leftChild == nil && n.rightChild == nil
}

func (n Line) getPosition(prevPosition int, goingLeft bool) int {
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

func (n Line) height() int {
	return 1 + log2(n.numberOfChildren+1)
}

func (n Line) getBalance() int {
	result := 0
	if n.leftChild != nil {
		result -= n.leftChild.height()
	}
	if n.rightChild != nil {
		result += n.rightChild.height()
	}
	return result
}

func (n *Line) rotateRight() {
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

func (n *Line) rotateLeft() {
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

func (n *Line) fixNumberOfChildren() {
	if n.leftChild == nil {
		(*n).numberOfChildren = 0
	} else {
		(*n).numberOfChildren = n.leftChild.numberOfChildren + 1
	}

	if n.rightChild != nil {
		(*n).numberOfChildren += n.rightChild.numberOfChildren + 1
	}
}

func (n *Line) balance() {
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

func (n *Line) toList() []*Line {
	var result = make([]*Line, 0, n.numberOfChildren+1)

	if n.leftChild != nil {
		result = n.leftChild.toList()
	}
	result = append(result, n)
	if n.rightChild != nil {
		result = append(result, n.rightChild.toList()...)
	}
	return result
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
