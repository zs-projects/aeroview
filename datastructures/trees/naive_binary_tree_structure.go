package trees

import (
	"fmt"
	"math"
)

// NaiveBinaryTreeStruct encode the structure of a binary tree
// in a bitset.
type NaiveBinaryTreeStruct struct {
	positions []uint8
}

// MakeNaiveBinaryTreeStructure make a bitset binary tree structure for a level order traversal of the tree.
func MakeNaiveBinaryTreeStructure(tr BinaryLevelOrderer) NaiveBinaryTreeStruct {
	u := make([]uint8, int(math.Pow(2, float64(len(tr.LevelOrder())))))
	for k, node := range tr.LevelOrder() {
		position := int(math.Pow(2, float64(k))) - 1
		if node.HasLeftChild {
			u[position] = 1
		}
		if node.HasRightChild {
			u[position+1] = 1
		}
	}
	return NaiveBinaryTreeStruct{u}
}

// LeftChild return the position of the left child of the node given the position of the current node.
func (b NaiveBinaryTreeStruct) LeftChild(position int) (int, error) {
	child := int(math.Pow(2, float64(position))) - 1
	if b.positions[child] == 1 {
		return child, nil
	}
	return -1, fmt.Errorf("%w at position %v", ErrNoLeftChildNode, position)
}

// RightChild return the position of the right child of the node given the position of the current node.
func (b NaiveBinaryTreeStruct) RightChild(position int) (int, error) {
	child := int(math.Pow(2, float64(position)))
	if b.positions[child] == 1 {
		return child, nil
	}
	return -1, fmt.Errorf("%w at position %v", ErrNoRightChildNode, position)
}

// HasLeftChild return whether the current node has a left child.
func (b NaiveBinaryTreeStruct) HasLeftChild(position int) bool {
	child := int(math.Pow(2, float64(position))) - 1
	return b.positions[child] == 1
}

// HasRightChild return whether the current node has a right child.
func (b NaiveBinaryTreeStruct) HasRightChild(position int) bool {
	child := int(math.Pow(2, float64(position)))
	return b.positions[child] == 1
}
