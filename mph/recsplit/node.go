package recsplit

import (
	"github.com/zs-projects/aeroview/datastructures/trees"
)

type Node struct {
	Left   *Node
	Right  *Node
	R      int
	nbKeys int
}

func FromCompactFBTree(fbt RecsplitStree) *Node {
	return fromCompactFBTree(fbt, 0)
}

func (n *Node) LevelOrder() []*Node {
	queue := []*Node{n}
	cur := 0
	for cur < len(queue) {
		current := queue[cur]
		curNode := trees.BinaryNode{}
		if current.Left != nil {
			curNode.HasLeftChild = true
			queue = append(queue, current.Left)
		}
		if current.Right != nil {
			curNode.HasRightChild = true
			queue = append(queue, current.Right)
		}
		cur++
	}
	return queue
}

func fromCompactFBTree(fbt RecsplitStree, currOffset int) *Node {
	var (
		left  *Node = nil
		right *Node = nil
	)
	if fbt.nodeHasLeftChild(currOffset) {
		lOffset := fbt.LeftChild(currOffset)
		left = fromCompactFBTree(fbt, lOffset)
	}
	if fbt.nodeHasRightChild(currOffset) {
		rOffset := fbt.RightChild(currOffset)
		right = fromCompactFBTree(fbt, rOffset)
	}
	R, nbkeys := fbt.node(currOffset)
	return &Node{
		Left:   left,
		Right:  right,
		R:      R,
		nbKeys: nbkeys,
	}
}
