package trees

import (
	"bytes"
	"fmt"
	"math"

	"zs-project.org/aeroview/datastructures/bits"
)

const (
	rootNodeOffset = 1
)

type FBValue struct {
	NbKeys int
	R      int
}

// FBNode reprensents a node in the flat binary tree.
type FBNode struct {
	offset int
	Value  FBValue
}

// FBTree stands from Flat Binary Tree.
// It is called flat becauses it does not used a pointer base data representation.
// We use a BitQueue to represent the structure of the tree.
// And in the future Rank and Select to have constant time access to parent/ left child and right child.
// TODO: Add Rank and Select.
type FBTree struct {
	nodes     []FBValue
	structure bits.Queue
}

func (f FBTree) String() string {
	w := bytes.NewBuffer(nil)
	for r := range f.structure.Data() {
		fmt.Fprintf(w, "%b", r)
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "%v", f.nodes)
	return w.String()
}

// Root Returns the root value of the tree.
func (f FBTree) Root() *FBNode {
	return f.node(0)
}

// LeftChild returns the left child of the node.
func (f FBTree) LeftChild(node FBNode) *FBNode {
	if f.nodeHasLeftChild(node) {
		nodePos := 2*node.offset + rootNodeOffset
		return f.node(nodePos)
	}
	return nil
}

// RightChild returns the right child of the node.
func (f FBTree) RightChild(node FBNode) *FBNode {
	if f.nodeHasRightChild(node) {
		nodePos := 2*node.offset + 1 + rootNodeOffset
		return f.node(nodePos)
	}
	return nil
}

func (f FBTree) nodeHasRightChild(node FBNode) bool {
	nodePos := 2*node.offset + 1
	return f.structure.Get(nodePos) == 1
}

func (f FBTree) nodeHasLeftChild(node FBNode) bool {
	nodePos := 2 * node.offset
	return f.structure.Get(nodePos) == 1
}

// IsLeaf returns tree if the node is a leaf node.
func (f FBTree) IsLeaf(node FBNode) bool {
	return f.structure.Get(2*node.offset) == 0 && f.structure.Get(2*node.offset+1) == 0
}

func (f FBTree) node(offset int) *FBNode {
	return &FBNode{
		offset: offset,
		Value:  f.nodes[offset],
	}
}

// TreeLeaf represent a leaf of a binary tree.
// It knows about it's current value and the value of all it's parents and
// the path from the root of the tree.
// Invariant : Len(TreeLeaf.Values) == Len(TreeLeaf.Path) + 1
type TreeLeaf interface {
	// Values return the values in the tree from the root to the leaf.
	Values() []FBValue
	// The path to the root of the tree, with false for left and true for right.
	Path() []bool
}

func maxDepth(tls []TreeLeaf) int {
	maxDepth := 0
	for _, tl := range tls {
		// We make sure that we have enough room to store this one.
		if len(tl.Values()) > maxDepth {
			maxDepth = len(tl.Values())
		}
	}
	return maxDepth
}

// MakeFBTreeFromLeafs creates a tree from a slice of TreeLeafs.
// order in the slice is not important.
func MakeFBTreeFromLeafs(tls []TreeLeaf) FBTree {
	depth := maxDepth(tls)
	tr := preallocateFBTree(depth)
	for _, v := range tls {
		parents := v.Values()
		tr.nodes[0] = parents[0]
		path := v.Path()
		position := 0
		for k, n := range parents {
			// 1. Set the node in it's correct position
			position = 2*position + int(math.Min(1, float64(k)))
			if k > 0 && path[k-1] {
				position++
			}
			tr.nodes[position] = n
			// 2. Set the structure in it's correct position
			if k < len(path) {
				// If there are still elements in the path, we then know for sure
				// that the current node has at least one additional child.
				// let's put a high bit at that position.
				if path[k] {
					tr.structure.High(2*position + 1)
				} else {
					tr.structure.High(2 * position)
				}
			}
		}
	}
	return tr
}

func preallocateFBTree(depth int) FBTree {
	tr := FBTree{
		nodes:     make([]FBValue, 0),
		structure: bits.Queue{},
	}
	tr.ensureCapacity(depth)
	tr.structure.Append(make([]uint64, int(math.Ceil(float64(len(tr.nodes))/8))), len(tr.nodes))
	return tr
}

// ensureCapacity ensure that there is enough capacity in FBTree to handle a binary tree of depth depth.
// it will allocate a new slice and copy the values otherwise.
func (f *FBTree) ensureCapacity(depth int) {
	nbNodes := int(math.Pow(2, float64(depth+1))) - 1
	if nbNodes > len(f.nodes) {
		nNodes := make([]FBValue, nbNodes)
		copy(nNodes, f.nodes)
		f.nodes = nNodes
	}
}
