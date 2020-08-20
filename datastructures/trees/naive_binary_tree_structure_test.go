package trees

import (
	"testing"
)

func TestMakeNaiveBinaryTreeStruct(t *testing.T) {
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
	*/
	lvo := []BinaryNode{
		{true, true},   // 5
		{true, false},  // 7
		{true, true},   // 8
		{false, true},  // 3
		{true, false},  // 6
		{true, false},  // 11
		{false, true},  // 9
		{false, false}, // 7
		{false, true},  // 15
	}
	bbts := MakeNaiveBinaryTreeStructure(lvo)
	// Let's naviguate through the tree.
	testTableHasLeft := map[int]bool{
		0: true,  // 5
		1: true,  // 7
		2: true,  // 8
		3: false, // 3
		4: true,  // 6
		5: true,  // 11
		6: false, // 9
		7: false, // 7
		8: false, // 15
	}
	for i := 0; i < len(lvo); i++ {
		if bbts.HasLeftChild(i) != testTableHasLeft[i] {
			t.Errorf("Has left check failed at index %v, expected %v, got %v", i, testTableHasLeft[i], bbts.HasLeftChild(i))
		}
	}
	// Let's naviguate through the tree.
	testTableHasRight := map[int]bool{
		0: true,  // 5
		1: false, // 7
		2: true,  // 8
		3: true,  // 3
		4: false, // 6
		5: false, // 11
		6: true,  // 9
		7: false, // 7
		8: true,  // 15
	}
	for i := 0; i < len(lvo); i++ {
		if bbts.HasRightChild(i) != testTableHasRight[i] {
			t.Errorf("Has right check failed at index %v, expected %v, got %v", i, testTableHasLeft[i], bbts.HasLeftChild(i))
		}
	}
	// Let's naviguate through the tree.
	testTableLeft := map[int]int{
		0: 0,  // 5
		1: 1,  // 7
		2: 3,  // 8
		4: 15, // 6
		5: 31, // 11
	}
	for i, expected := range testTableLeft {
		if off, ok := bbts.LeftChild(i); ok != nil || (off != expected) {
			t.Errorf("Left child method failed at index %v, expected %v, got %v", i, expected, off)
		}
	}
	// Let's naviguate through the tree.
	testTableRight := map[int]int{
		0: 1, // 5
		2: 4, // 8
		3: 8, // 3
	}
	for i, expected := range testTableRight {
		if off, ok := bbts.RightChild(i); ok != nil || (off != expected) {
			t.Errorf("Right child method failed at index %v, expected %v, got %v", i, expected, off)
		}
	}
}
