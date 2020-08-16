package trees

import (
	"testing"
)

type MockLevelOrderer struct {
	nodes []NodeMetadata
}

func (ml MockLevelOrderer) LevelOrder() []NodeMetadata {
	return ml.nodes
}

func TestMakeBitsetBinaryTreeStruct(t *testing.T) {
	/* We are testing this tree structure :
					5
				   / \
				  /   \
				 /     \
				7       8
			   / \     / \
	          /   \   /   \
			 3     x 6    11
			/ \     / \   / \
		   x   9   7   x 15  x
			  / \       / \
			 x   4     x  17
			    / \
			   2   x
	*/
	lvo := MockLevelOrderer{
		nodes: []NodeMetadata{
			{true, true},   // 5
			{true, false},  // 7
			{true, true},   // 8
			{false, true},  // 3
			{true, false},  // 6
			{true, false},  // 11
			{false, true},  // 9
			{false, false}, // 7
			{false, true},  // 15
			{true, false},  // 4
			{false, false}, // 17
			{false, false}, // 2
		},
	}
	expectedStructure := uint64(0b01100010010110110111)
	bbts := MakeBitsetBinaryTreeStructure(lvo)
	// Checking for structure
	for i := 0; i < 24; i++ {
		bit := bbts.Get(i)
		expected := int((expectedStructure >> i) & 0b1)
		if bit != expected {
			t.Errorf("Structure check failed at index %v, expected %v, got %v", i, expected, bit)
		}
	}
	// Let's naviguate through the tree.
	testTableHasLeft := map[int]bool{
		0:  true,  // 5
		1:  true,  // 7
		2:  true,  // 8
		3:  false, // 3
		4:  true,  // 6
		5:  true,  // 11
		6:  false, // 9
		7:  false, // 7
		8:  false, // 15
		9:  true,  // 4
		10: false, // 17
		11: false, // 2
	}
	for i := 0; i < len(lvo.LevelOrder()); i++ {
		if bbts.HasLeftChild(i) != testTableHasLeft[i] {
			t.Errorf("Has left check failed at index %v, expected %v, got %v", i, testTableHasLeft[i], bbts.HasLeftChild(i))
		}
	}
	// Let's naviguate through the tree.
	testTableHasRight := map[int]bool{
		0:  true,  // 5
		1:  false, // 7
		2:  true,  // 8
		3:  true,  // 3
		4:  false, // 6
		5:  false, // 11
		6:  true,  // 9
		7:  false, // 7
		8:  true,  // 15
		9:  false, // 4
		10: false, // 17
		11: false, // 2
	}
	for i := 0; i < len(lvo.LevelOrder()); i++ {
		if bbts.HasRightChild(i) != testTableHasRight[i] {
			t.Errorf("Has right check failed at index %v, expected %v, got %v", i, testTableHasLeft[i], bbts.HasLeftChild(i))
		}
	}
	// Let's naviguate through the tree.
	testTableLeft := map[int]int{
		0: 1,  // 5
		1: 3,  // 7
		2: 4,  // 8
		4: 7,  // 6
		5: 8,  // 11
		9: 11, // 4
	}
	for i, expected := range testTableLeft {
		if off, ok := bbts.LeftChild(i); ok != nil || (off != expected) {
			t.Errorf("Right child method failed at index %v, expected %v, got %v", i, expected, off)
		}
	}
	// Let's naviguate through the tree.
	testTableRight := map[int]int{
		0: 2,  // 5
		2: 5,  // 8
		3: 6,  // 3
		6: 9,  // 9
		8: 10, // 15
	}
	for i, expected := range testTableRight {
		if off, ok := bbts.RightChild(i); ok != nil || (off != expected) {
			t.Errorf("Right child method failed at index %v, expected %v, got %v", i, expected, off)
		}
	}
}
