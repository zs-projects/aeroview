package trees

import (
	"errors"
)

var (
	// ErrNoLeftChildNode is returned when we are trying to access the non existant left child of a node.
	ErrNoLeftChildNode = errors.New("no left child for the provided node")
	// ErrNoRightChildNode is returned when we are trying to access the non existant right child of a node.
	ErrNoRightChildNode = errors.New("no right child for the provided node")
)

// NodeMetadata represents whether a node has a left and/or a right child
type NodeMetadata struct {
	HasLeftChild  bool
	HasRightChild bool
}

// LevelOrderer traverses the underlying binary tree in level order a provides for each node
// whether it has a left and/or right child.
type LevelOrderer interface {
	LevelOrder() []NodeMetadata
}

// BinaryTreeStructure is a interface that provides basic operations for trees.
type BinaryTreeStructure interface {
	HasLeftChild(position int) bool
	HasRightChild(position int) bool
	LeftChild(position int) (int, error)
	RightChild(position int) (int, error)
}
