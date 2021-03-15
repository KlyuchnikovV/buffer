package buffer

import (
	"fmt"

	"github.com/KlyuchnikovV/buffer/line"
	"github.com/KlyuchnikovV/stack"
)

type BufferTree struct {
	root *Node
	size int
}

func NewTree(data line.Line) *BufferTree {
	return &BufferTree{
		root: newNode(data),
		size: 1,
	}
}

func (tree *BufferTree) Insert(data line.Line, position int) error {
	if position > tree.size || position < 0 {
		return fmt.Errorf("position \"%d\" is out of tree range: [%d - %d]", position, 0, tree.size+1)
	}
	tree.size++
	if tree.root == nil {
		tree.root = newNode(data)
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
			node.leftChild = newNode(data)
			break loop
		case currentPosition < position && !node.HasRight():
			node.rightChild = newNode(data)
			break loop
		}
	}

	for v, ok := nodeStack.Pop(); ok; v, ok = nodeStack.Pop() {
		v.(*Node).balance()
	}

	return nil
}

func (tree *BufferTree) GetNode(position int) *Node {
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

func (tree *BufferTree) Find(position int) (*line.Line, bool) {
	node := tree.GetNode(position)
	if node == nil {
		return nil, false
	}
	return node.Data(), true
}

func (tree *BufferTree) Size() int {
	return tree.size
}

func (tree *BufferTree) ToList() []*line.Line {
	return tree.root.toList()
}

func (tree *BufferTree) Delete(position int) (*line.Line, bool) {
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
		} else if node.IsLeftOf(*v.(*Node)) {
			v.(*Node).leftChild = nil
		} else {
			v.(*Node).rightChild = nil
		}
	} else {
		*node = *node.leftChild
	}

	for v, ok := nodeStack.Pop(); ok; v, ok = nodeStack.Pop() {
		v.(*Node).balance()
	}

	tree.size--
	return result, true
}

func (tree *BufferTree) Visualize() {
	if tree.root == nil {
		return
	}
	tree.root.visualizeNodeSubtree(0, tree.root.height())
}
