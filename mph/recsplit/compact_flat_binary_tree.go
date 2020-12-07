package recsplit

import (
	"github.com/zs-projects/aeroview/datastructures/bits"
	"github.com/zs-projects/aeroview/rank"
)

type CompactFBTree struct {
	nodesR      []int
	nodesNBKeys []int
	structure   *rank.PopCount
}

func (c CompactFBTree) SizeInBytes() int {
	return len(c.nodesR)*8 + len(c.nodesNBKeys)*8 + c.structure.SizeInBytes()
}

func MakeFbTreeFromRecSplitSubTree(rst recsplitSubTree) CompactFBTree {
	queue := []*Node{rst.Node}
	nodesNBKeys := []int{rst.nbKeys}
	nodesR := []int{rst.R}
	cur := 0
	q := bits.MakeQueue()
	for cur < len(queue) {
		current := queue[cur]
		if current.Left != nil {
			q.PushBack(1)
			queue = append(queue, current.Left)
			nodesNBKeys = append(nodesNBKeys, current.Left.nbKeys)
			nodesR = append(nodesR, current.Left.R)
		} else {
			q.PushBack(0)
		}
		if current.Right != nil {
			q.PushBack(1)
			queue = append(queue, current.Right)
			nodesNBKeys = append(nodesNBKeys, current.Right.nbKeys)
			nodesR = append(nodesR, current.Right.R)
		} else {
			q.PushBack(0)
		}
		cur++
	}
	structure := rank.MakePopCount(q.Vector())
	return CompactFBTree{
		nodesNBKeys: nodesNBKeys,
		structure:   structure,
		nodesR:      nodesR,
	}
}

func (c CompactFBTree) Root() int {
	return 0
}

func (c CompactFBTree) LeftChild(offset int) int {
	return int(c.structure.Get(2*offset)) * c.structure.Rank(2*offset)
}

func (c CompactFBTree) RightChild(offset int) int {
	return int(c.structure.Get(2*offset+1)) * c.structure.Rank(2*offset+1)
}

func (c CompactFBTree) nodeHasLeftChild(offset int) bool {
	return c.structure.Get(2*offset) == 1
}

func (c CompactFBTree) nodeHasRightChild(offset int) bool {
	return c.structure.Get(2*offset+1) == 1
}

func (c CompactFBTree) node(offset int) (r, nbKeys int) {
	return c.nodesR[offset], c.nodesNBKeys[offset]
}

func (c CompactFBTree) IsLeaf(offset int) bool {
	return !(c.structure.Get(2*offset) == 1 || c.structure.Get(2*offset+1) == 1)
}
