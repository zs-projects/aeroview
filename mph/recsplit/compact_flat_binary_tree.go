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

func FromFBTree(u FBTree) CompactFBTree {
	q := bits.MakeQueue()
	nodesR, nodesNBKeys := compressStructure(&q, &u)
	b := q.Vector()
	pc := rank.MakePopCount(b)
	return CompactFBTree{
		nodesNBKeys: nodesNBKeys,
		nodesR:      nodesR,
		structure:   pc,
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

func compressStructure(q *bits.Queue, u *FBTree) ([]int, []int) {
	nodesQueue := make([]*FBNode, 0)
	nodesQueue = append(nodesQueue, u.Root())
	nodesR := make([]int, 0)
	nodesNBKeys := make([]int, 0)
	nodesR = append(nodesR, u.nodes[0].R)
	nodesNBKeys = append(nodesNBKeys, u.nodes[0].NbKeys)
	i := 0
	for len(nodesQueue) > i {
		node := nodesQueue[i]
		i++
		if u.nodeHasLeftChild(*node) {
			q.PushBack(1)
			left := u.LeftChild(*node)
			nodesQueue = append(nodesQueue, left)
			nodesNBKeys = append(nodesNBKeys, u.nodes[left.offset].NbKeys)
			nodesR = append(nodesR, u.nodes[left.offset].R)
		} else {
			q.PushBack(0)
		}
		if u.nodeHasRightChild(*node) {
			q.PushBack(1)
			right := u.RightChild(*node)
			nodesQueue = append(nodesQueue, right)
			nodesR = append(nodesR, u.nodes[right.offset].R)
			nodesNBKeys = append(nodesNBKeys, u.nodes[right.offset].NbKeys)
		} else {
			q.PushBack(0)
		}
	}
	return nodesR, nodesNBKeys
}

func (f CompactFBTree) node(offset int) (R int, nbKeys int) {
	return f.nodesR[offset], f.nodesNBKeys[offset]
}

func (f CompactFBTree) IsLeaf(offset int) bool {
	return !(f.structure.Get(2*offset) == 1 || f.structure.Get(2*offset+1) == 1)
}
