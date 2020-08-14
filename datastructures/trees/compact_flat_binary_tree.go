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
		nodePosition := c.structure.Rank(2 * u.offset)
		return c.node(nodePosition)
	}
	return nil
}

func (c CompactFBTree) RightChild(u *FBNode) *FBNode {
	if c.nodeHasRightChild(*u) {
		nodePosition := c.structure.Rank(2*u.offset + 1)
		return c.node(nodePosition)
	}
	return nil
}

func (c CompactFBTree) nodeHasLeftChild(node FBNode) bool {
	offset := node.offset
	position := 2 * offset
	exists := c.structure.Get(position)
	return exists == 1
}

func (c CompactFBTree) nodeHasRightChild(node FBNode) bool {
	offset := node.offset
	position := 2*offset + 1
	return c.structure.Get(position) == 1
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

func (f CompactFBTree) IsLeaf(node FBNode) bool {
	return !f.nodeHasLeftChild(node) && !f.nodeHasRightChild(node)
}
