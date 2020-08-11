package datastructures

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
			/ \     / \
		   x   9   7   x
	*/
	expectedNodes := []FBValue{{0, 5}, {0, 7}, {0, 8}, {0, 3}, {0, 6}, {0, 11}, {0, 9}, {0, 7}}
	expectedStructure := []uint64{0b1110110110 << 54}
	tls := []TreeLeafS{
		{
			values: []FBValue{{0, 5}, {0, 7}, {0, 3}, {0, 9}},
			path:   []bool{false, false, true}},
		{
			values: []FBValue{{0, 5}, {0, 7}, {0, 3}, {0, 9}},
			path:   []bool{false, false, true}},
		{
			values: []FBValue{{0, 5}, {0, 8}, {0, 6}, {0, 7}},
			path:   []bool{true, false, false}},
		{
			values: []FBValue{{0, 5}, {0, 8}, {0, 11}},
			path:   []bool{true, true}}}
	data := make([]TreeLeaf, 0)
	for _, v := range tls {
		data = append(data, v)
	}
	fbt := FromFBTree(MakeFBTreeFromLeafs(data))
	if !reflect.DeepEqual(expectedNodes, fbt.nodes[:len(expectedNodes)]) {
		t.Errorf("node layout is not good, expected %v \n got %v", expectedNodes, fbt.nodes)
	}
	if !reflect.DeepEqual(fbt.structure.Data[0], expectedStructure[0]) {
		t.Errorf("node structure is not good, expected \n %b \n got \n %b", expectedStructure, fbt.structure.Data)
	}
	r := fbt.Root()
	if r.Value.R != 5 {
		t.Errorf("Root of the tree should be 5 got %v", r.Value)
	}
	l := fbt.LeftChild(r)
	if l.Value.R != 7 {
		t.Errorf("Left Child of Root of the tree should be 7 got %v", l)
	}
	lr := fbt.RightChild(l)
	if lr != nil {
		t.Errorf("Right Child of Left Child of Root should be nil, got %v", lr.Value)
	}
	rg := fbt.RightChild(r)
	if rg.Value.R != 8 {
		t.Errorf("Left Child of Root of the tree should be 8 got %v", rg)
	}

}
