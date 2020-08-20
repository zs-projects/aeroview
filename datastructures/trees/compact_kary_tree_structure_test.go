package trees

import (
	"testing"
)

type MockKAryLevelOrderer struct {
	structure []int
}

func (m MockKAryLevelOrderer) LevelOrder() []int {
	return m.structure
}

func TestMakeCompactKAryTreeStructure(t *testing.T) {
	/* We are testing this tree structure :
				   _______
				  /|  |  |\
				 / |  |  | \
				/  |  |  |  \
	           /   |  |  |   \
			  7    5  1  6    8
	        ____     / \     /|\
	       /|  |\   4   7   8 4 8
	      / |  | \     / \
	     /  |  |  \   3   1
		1   3  5   9
	*/
	lvo := []int{5, 4, 0, 2, 0, 3, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0}
	rdxt := MakeCompactKAryTreeStructure(lvo)
	structureCheck := map[int]int{
		0:  1,
		1:  6,
		3:  10,
		5:  12,
		11: 15}
	for i, nbChildrenExpected := range lvo {
		if nbChildren, _ := rdxt.Children(i); nbChildren != nbChildrenExpected {
			t.Errorf("number of children is wrong for node %v expected %v got  %v", i, nbChildrenExpected, nbChildren)
		}
	}
	for ndIdx, childrendPos := range structureCheck {

		if _, pos := rdxt.Children(ndIdx); pos != childrendPos {
			t.Errorf("children position is wrong for node %v expected position %v got position %v", ndIdx, childrendPos, pos)
		}
	}
}
