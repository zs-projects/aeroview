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

// BinaryNode represents whether a node has a left and/or a right child
type BinaryNode struct {
	HasLeftChild  bool
	HasRightChild bool
}

// BinaryLevelOrder provides a traversal of the underlying binary tree in level order and
// tells for each node whether it has a left and/or right child.
type BinaryLevelOrder []BinaryNode

// NodeNBChildren is the number of children o a node in a KAry Tree
type NodeNBChildren = int

// KAryLevelOrder traverses the underlying k-ary tree in level order a provides for each node
// how many children he has.
type KAryLevelOrder []NodeNBChildren

// BinaryTreeStructure is an interface that provides basic operations for binary trees.
type BinaryTreeStructure interface {
	LeftChild(position int) (childPosition int, exists bool)
	RightChild(position int) (childPosition int, exists bool)
}

//KAryTreeStructure is an interface that provides basic operations for k-ary trees.
type KAryTreeStructure interface {
	Children(nodPos int) (nbChildren, startPosition int)
}
