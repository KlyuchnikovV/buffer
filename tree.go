package buffer

import (
	"fmt"

	"github.com/KlyuchnikovV/stack"
)

type BufferTree struct {
	root *Line
	size int
}

func NewTree(data []rune) *BufferTree {
	return &BufferTree{
		root: NewLine(data),
		size: 1,
	}
}

func (tree *BufferTree) Insert(data []rune, position int) error {
	if position > tree.size || position < 0 {
		return fmt.Errorf("position \"%d\" is out of tree range: [%d - %d]", position, 0, tree.size+1)
	}
	tree.size++
	if tree.root == nil {
		tree.root = NewLine(data)
		return nil
	}

	node := tree.root
	nodeStack := stack.New(tree.root.height())
	currentPosition := node.getPosition(-1, false)

loop:
	for {
		node.numberOfChildren++
		switch {
		case currentPosition >= position && node.HasLeft():
			nodeStack.Push(node)
			node = node.leftChild
			currentPosition = node.getPosition(currentPosition, true)
		case currentPosition < position && node.HasRight():
			nodeStack.Push(node)
			node = node.rightChild
			currentPosition = node.getPosition(currentPosition, false)
		case currentPosition >= position && !node.HasLeft():
			node.leftChild = NewLine(data)
			break loop
		case currentPosition < position && !node.HasRight():
			node.rightChild = NewLine(data)
			break loop
		}
	}

	for v, ok := nodeStack.Pop(); ok; v, ok = nodeStack.Pop() {
		v.(*Line).balance()
	}

	return nil
}

func (tree *BufferTree) GetNode(position int) *Line {
	if tree.size == 0 || position < 0 || position > tree.size {
		return nil
	}
	node := tree.root
	currentPos := node.getPosition(-1, false)

	for node != nil {
		switch {
		case position < currentPos:
			node = (*node).leftChild
			currentPos = node.getPosition(currentPos, true)
			continue
		case position > currentPos:
			node = (*node).rightChild
			currentPos = node.getPosition(currentPos, false)
			continue
		}
		break
	}
	return node
}

func (tree *BufferTree) Find(position int) (*Line, bool) {
	node := tree.GetNode(position)
	if node == nil {
		return nil, false
	}
	return node, true
}

func (tree *BufferTree) Size() int {
	return tree.size
}

func (tree *BufferTree) ToList() []*Line {
	return tree.root.toList()
}

func (tree *BufferTree) Delete(position int) ([]rune, bool) {
	if position >= tree.size || position < 0 {
		return nil, false
	}

	node := tree.root
	nodeStack := stack.New(tree.root.height())
	currentPosition := node.getPosition(-1, false)

	for currentPosition != position {
		node.numberOfChildren--
		nodeStack.Push(node)
		if currentPosition > position {
			node = node.leftChild
			currentPosition = node.getPosition(currentPosition, true)
		} else if currentPosition < position {
			node = node.rightChild
			currentPosition = node.getPosition(currentPosition, false)
		}
	}

	result := node.data

	if node.HasRight() {

		var (
			parentNode    = node
			replacingNode = node.rightChild
		)
		for replacingNode.HasLeft() {
			parentNode = replacingNode
			replacingNode.numberOfChildren--
			replacingNode = replacingNode.leftChild
		}
		if parentNode != node {
			parentNode.leftChild = nil
		}
		replacingNode.leftChild = node.leftChild

		if node.rightChild != replacingNode {
			replacingNode.rightChild = node.rightChild
		}
		*node = *replacingNode
		node.fixNumberOfChildren()
	} else if !node.HasLeft() {
		v, ok := nodeStack.Peek()
		if !ok {
			tree.root = nil
		} else if node.IsLeftOf(*v.(*Line)) {
			v.(*Line).leftChild = nil
		} else {
			v.(*Line).rightChild = nil
		}
	} else {
		*node = *node.leftChild
	}

	for v, ok := nodeStack.Pop(); ok; v, ok = nodeStack.Pop() {
		v.(*Line).balance()
	}

	tree.size--
	return result, true
}
