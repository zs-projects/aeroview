package trees

import (
	"zs-project.org/aeroview/datastructures/bits"
	"zs-project.org/aeroview/rank"
)

type CompactFBTree struct {
	nodes     []FBValue
	structure rank.PopCount
}

func FromFBTree(u FBTree) CompactFBTree {
	q := bits.MakeQueue()
	nodes := compressStructure(&q, &u)
	b := q.Vector()
	pc := rank.MakePopCount(b)
	return CompactFBTree{
		nodes:     nodes,
		structure: pc,
	}
}
func (c CompactFBTree) Root() *FBNode {
	return c.node(0)
}
func (c CompactFBTree) LeftChild(u *FBNode) *FBNode {
	if c.nodeHasLeftChild(*u) {
		rk := c.structure.Rank(u.offset)
		position := 2*rk - 1
		return c.node(position)
	}
	return nil
}

func (c CompactFBTree) RightChild(u *FBNode) *FBNode {
	if c.nodeHasRightChild(*u) {
		rk := c.structure.Rank(u.offset)
		position := 2 * rk
		return c.node(position)
	}
	return nil
}

func (c CompactFBTree) nodeHasRightChild(node FBNode) bool {
	rk := c.structure.Rank(node.offset)
	position := 2 * rk
	return c.structure.Get(64-position) == 1
}

func (c CompactFBTree) nodeHasLeftChild(node FBNode) bool {
	rk := c.structure.Rank(node.offset)
	position := 2*rk - 1
	exists := c.structure.Get(64 - position)
	return exists == 1
}

func compressStructure(q *bits.Queue, u *FBTree) []FBValue {
	nodesQueue := make([]*FBNode, 0)
	nodesQueue = append(nodesQueue, u.Root())
	nodes := make([]FBValue, 0)
	nodes = append(nodes, u.nodes[0])
	i := 0
	for len(nodesQueue) > i {
		node := nodesQueue[i]
		if u.IsLeaf(*node) {
			i++
			continue
		}
		if u.nodeHasLeftChild(*node) {
			q.PushBack(1)
			left := u.LeftChild(*node)
			nodesQueue = append(nodesQueue, left)
			nodes = append(nodes, u.nodes[left.offset])
		} else {
			q.PushBack(0)
		}
		if u.nodeHasRightChild(*node) {
			q.PushBack(1)
			right := u.RightChild(*node)
			nodesQueue = append(nodesQueue, right)
			nodes = append(nodes, u.nodes[right.offset])
		} else {
			q.PushBack(0)
		}
		i++
	}
	return nodes
}

func (f CompactFBTree) node(offset int) *FBNode {
	return &FBNode{
		offset: offset,
		Value:  f.nodes[offset],
	}
}
