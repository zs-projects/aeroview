package trees

import (
	"fmt"

	"github.com/zs-projects/aeroview/datastructures/bits"
	"github.com/zs-projects/aeroview/rank"
)

// BitsetBinaryTreeStruct encode the structure of a binary tree
// in a bitset.
type BitsetBinaryTreeStruct struct {
	rank.GetRanker
}

// MakeBitsetBinaryTreeStructure make a bitset binary tree structure for a level order traversal of the tree.
func MakeBitsetBinaryTreeStructure(tr LevelOrderer) BitsetBinaryTreeStruct {
	q := bits.MakeQueue()
	for _, node := range tr.LevelOrder() {
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
	return BitsetBinaryTreeStruct{rank.MakePopCount(q.Vector())}
}

// LeftChild return the position of the left child of the node given the position of the current node.
func (b BitsetBinaryTreeStruct) LeftChild(position int) (int, error) {
	if b.HasLeftChild(position) {
		return b.Rank(2 * position), nil
	}
	return -1, fmt.Errorf("%w at position %v", ErrNoLeftChildNode, position)
}

// RightChild return the position of the right child of the node given the position of the current node.
func (b BitsetBinaryTreeStruct) RightChild(position int) (int, error) {
	if b.HasRightChild(position) {
		return b.Rank(2*position + 1), nil
	}
	return -1, fmt.Errorf("%w at position %v", ErrNoRightChildNode, position)
}

// HasLeftChild return whether the current node has a left child.
func (b BitsetBinaryTreeStruct) HasLeftChild(position int) bool {
	return b.Get(2*position) == 1

}

// HasRightChild return whether the current node has a right child.
func (b BitsetBinaryTreeStruct) HasRightChild(position int) bool {
	return b.Get(2*position+1) == 1
}
