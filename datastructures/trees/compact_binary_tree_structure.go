package trees

import (
	"github.com/zs-projects/aeroview/datastructures/bits"
	"github.com/zs-projects/aeroview/rank"
)

// CompactBinaryTreeStructure encode the structure of a binary tree
// in a bitset.
type CompactBinaryTreeStructure struct {
	rank.GetRanker
}

// MakeCompactBinaryTreeStructure make a bitset binary tree structure for a level order traversal of the tree.
func MakeCompactBinaryTreeStructure(tr BinaryLevelOrder) CompactBinaryTreeStructure {
	q := bits.MakeQueue()
	for _, node := range tr {
		if node.HasLeftChild {
			q.PushBack(1)
		} else {
			q.PushBack(0)
		}
		if node.HasRightChild {
			q.PushBack(1)
		} else {
			q.PushBack(0)
		}
	}
	return CompactBinaryTreeStructure{rank.MakePopCount(q.Vector())}
}

// LeftChild return the position of the left child of the node given the position of the current node.
func (b CompactBinaryTreeStructure) LeftChild(position int) (offset int, exists bool) {
	existsI := b.Get(2 * position)
	return int(existsI) * b.Rank(2*position), existsI == 1
}

// RightChild return the position of the right child of the node given the position of the current node.
func (b CompactBinaryTreeStructure) RightChild(position int) (int, bool) {
	existsI := b.Get(2*position + 1)
	return int(existsI) * b.Rank(2*position+1), existsI == 1
}
