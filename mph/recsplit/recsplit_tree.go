package recsplit

import (
	"github.com/zs-projects/aeroview/datastructures/trees"
)

type RecsplitStree struct {
	nodesR      []int
	nodesNBKeys []int
	structure   trees.CompactBinaryTreeStructure
}

func (c RecsplitStree) SizeInBytes() int {
	return len(c.nodesR)*8 + len(c.nodesNBKeys)*8
}

func MakeFbTreeFromRecSplitSubTree(rst recsplitSubTree) RecsplitStree {
	levelOrder := rst.LevelOrder()
	nodesNBKeys := []int{levelOrder[0].nbKeys}
	nodesR := []int{levelOrder[0].R}
	nodesStructure := []trees.BinaryNode{}
	for _, node := range levelOrder {
		currentStruc := trees.BinaryNode{}
		if node.Left != nil {
			currentStruc.HasLeftChild = true
			nodesNBKeys = append(nodesNBKeys, node.Left.nbKeys)
			nodesR = append(nodesR, node.Left.R)
		}
		if node.Right != nil {
			currentStruc.HasRightChild = true
			nodesNBKeys = append(nodesNBKeys, node.Right.nbKeys)
			nodesR = append(nodesR, node.Right.R)
		}
		nodesStructure = append(nodesStructure, currentStruc)
	}
	structure := trees.MakeCompactBinaryTreeStructure(nodesStructure)
	return RecsplitStree{
		nodesNBKeys: nodesNBKeys,
		structure:   structure,
		nodesR:      nodesR,
	}
}

func (c RecsplitStree) Root() int {
	return 0
}

func (c RecsplitStree) LeftChild(offset int) int {
	position, _ := c.structure.LeftChild(offset)
	return position
}

func (c RecsplitStree) RightChild(offset int) int {
	position, _ := c.structure.RightChild(offset)
	return position
}

func (c RecsplitStree) nodeHasLeftChild(offset int) bool {
	_, ok := c.structure.LeftChild(offset)
	return ok
}

func (c RecsplitStree) nodeHasRightChild(offset int) bool {
	_, ok := c.structure.RightChild(offset)
	return ok
}

func (c RecsplitStree) node(offset int) (r, nbKeys int) {
	return c.nodesR[offset], c.nodesNBKeys[offset]
}

func (c RecsplitStree) IsLeaf(offset int) bool {
	return !(c.nodeHasLeftChild(offset) || c.nodeHasRightChild(offset))
}
