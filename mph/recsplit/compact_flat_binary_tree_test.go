package recsplit

import (
	"reflect"
	"testing"
)

func TestMakeCompactFBTreeFromLeafs(t *testing.T) {
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
	expectedNodesNbKeys := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	expectedNodesR := []int{5, 7, 8, 3, 6, 11, 9, 7, 15, 4, 17, 2}
	expectedStructure := []uint64{0b01100010010110110111}

	tls := []TreeLeafS{
		{
			values: []FBValue{{0, 5}, {0, 7}, {0, 3}, {0, 9}},
			path:   []bool{false, false, true}},
		{
			values: []FBValue{{0, 5}, {0, 7}, {0, 3}, {0, 9}, {0, 4}, {0, 2}},
			path:   []bool{false, false, true, true, false}},
		{
			values: []FBValue{{0, 5}, {0, 8}, {0, 6}, {0, 7}},
			path:   []bool{true, false, false}},
		{
			values: []FBValue{{0, 5}, {0, 8}, {0, 11}, {0, 15}, {0, 17}},
			path:   []bool{true, true, false, true}}}
	data := make([]TreeLeaf, 0)
	for _, v := range tls {
		data = append(data, v)
	}
	fbt := FromFBTree(MakeFBTreeFromLeafs(data))
	if !reflect.DeepEqual(expectedNodesR, fbt.nodesR[:len(expectedNodesR)]) {
		t.Errorf("node layout is not good, expected %v \n got %v", expectedNodesR, fbt.nodesR)
	}
	if !reflect.DeepEqual(expectedNodesNbKeys, fbt.nodesNBKeys[:len(expectedNodesNbKeys)]) {
		t.Errorf("node layout is not good, expected %v \n got %v", expectedNodesNbKeys, fbt.nodesNBKeys)
	}
	if !reflect.DeepEqual(fbt.structure.Data[0], expectedStructure[0]) {
		t.Errorf("node structure is not good, expected \n %b \n got \n %b", expectedStructure, fbt.structure.Data)
	}
	ranks := make([]int, len(expectedNodesR))
	for i := 0; i < len(ranks); i++ {
		ranks[i] = fbt.structure.Rank(i)
	}
	expectedRanks := []int{1, 2, 3, 3, 4, 5, 5, 6, 7, 7, 8, 8}
	if !reflect.DeepEqual(ranks, expectedRanks) {
		t.Errorf("ranks are not good, expected \n %v \n got \n %v", expectedRanks, ranks)
	}

	r := fbt.Root()
	VR, _ := fbt.node(r)
	if VR != 5 {
		t.Errorf("Root of the tree should be 5 got %v", VR)
	}
	l := fbt.LeftChild(r)
	VR, _ = fbt.node(l)
	if VR != 7 {
		t.Errorf("Left Child of Root of the tree should be 7 got %v", l)
	}
	if fbt.IsLeaf(l) {
		t.Errorf("Right child of Left Child of Left Child of Root of the tree should be a leaf")
	}
	ll := fbt.LeftChild(l)
	VR, _ = fbt.node(ll)
	if VR != 3 {
		t.Errorf("Left Child of Left Child of Root of the tree should be 3 got %v", ll)
	}
	if fbt.IsLeaf(ll) {
		t.Errorf("Right child of Left Child of Left Child of Root of the tree should be a leaf")
	}
	lr := fbt.RightChild(l)
	if lr != 0 {
		t.Errorf("Right Child of Left Child of Root should be nil, got %v", lr)
	}
	lll := fbt.LeftChild(ll)
	if lll != 0 {
		t.Errorf(" Left child of Left Child of Left Child of Root of the tree should be nil got %v", lll)
	}
	llr := fbt.RightChild(ll)
	VR, _ = fbt.node(llr)
	if VR != 9 {
		t.Errorf("Right child of Left Child of Left Child of Root of the tree should be 9 got %v", llr)
	}
	if fbt.IsLeaf(llr) {
		t.Errorf("Right child of Left Child of Left Child of Root of the tree should not be a leaf")
	}

	rg := fbt.RightChild(r)
	VR, _ = fbt.node(rg)
	if VR != 8 {
		t.Errorf("Right Child of Root of the tree should be 8 got %v", rg)
	}
	rgl := fbt.LeftChild(rg)
	VR, _ = fbt.node(rgl)
	if VR != 6 {
		t.Errorf("Left Child of Right Child of Root of the tree should be 6 got %v", rgl)
	}
	if fbt.IsLeaf(rgl) {
		t.Errorf("Left child of Right Child of Root of the tree should be a leaf")
	}
	rgll := fbt.LeftChild(rgl)
	VR, _ = fbt.node(rgll)
	if VR != 7 {
		t.Errorf("Left Child Left Child of Right Child of Root of the tree should be 7 got %v", rgll)
	}
	if !fbt.IsLeaf(rgll) {
		t.Errorf("Left child of Right Child of Root of the tree should be a leaf")
	}
}
