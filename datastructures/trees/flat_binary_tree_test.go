package trees

import (
	"math"
	"reflect"
	"testing"
)

type TreeLeafS struct {
	values []FBValue
	path   []bool
}

func (s TreeLeafS) Values() []FBValue {
	return s.values
}

func (s TreeLeafS) Path() []bool {
	return s.path
}

func TestMakeFBTreeFromLeafs(t *testing.T) {
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
	expectedNodes := []FBValue{{0, 5}, {0, 7}, {0, 8}, {0, 3}, {0, 0}, {0, 6}, {0, 11}, {0, 0}, {0, 9}, {0, 0}, {0, 0}, {0, 7}}
	expectedStructure := []uint64{uint64(0b010010110111)}
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
	fbt := MakeFBTreeFromLeafs(data)
	if !reflect.DeepEqual(expectedNodes, fbt.nodes[:12]) {
		t.Errorf("node layout is not good, expected %v \n got %v", expectedNodes, fbt.nodes)
	}
	if !reflect.DeepEqual(fbt.structure.Data(), expectedStructure) {
		t.Errorf("node structure is not good, expected %v \n got %v", expectedStructure, fbt.structure.Data())
	}
	r := fbt.Root()
	if r.Value.R != 5 {
		t.Errorf("Root of the tree should be 5 got %v", r.Value)
	}
	l := fbt.LeftChild(*r)
	if l.Value.R != 7 {
		t.Errorf("Left Child of Root of the tree should be 7 got %v", l)
	}
	lr := fbt.RightChild(*l)
	if lr != nil {
		t.Errorf("Right Child of Left Child of Root should be nil, got %v", lr.Value)
	}
	rg := fbt.RightChild(*r)
	if rg.Value.R != 8 {
		t.Errorf("Left Child of Root of the tree should be 8 got %v", rg)
	}

}

func TestPreallocateFBTree(t *testing.T) {
	depth := 5
	tr := preallocateFBTree(depth)
	nbNodes := int(math.Pow(2, float64(depth+1))) - 1
	if len(tr.nodes) != nbNodes {
		t.Errorf("preallocateFBTree failed, expecting %v, got %v.", nbNodes, len(tr.nodes))
	}
}

func TestMaxDepth(t *testing.T) {
	tls := []TreeLeafS{
		{
			values: []FBValue{{0, 2}, {0, 4}, {0, 6}},
			path:   []bool{true, false}},
		{
			values: []FBValue{{0, 2}, {0, 3}, {0, 5}, {0, 7}},
			path:   []bool{false, false, true}}}
	data := make([]TreeLeaf, 0)
	for _, v := range tls {
		data = append(data, v)
	}
	if maxDepth(data) != 4 {
		t.Errorf("Max Depth failed, expecting %v, got %v.", 4, maxDepth(data))
	}
}
