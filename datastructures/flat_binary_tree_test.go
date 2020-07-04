package datastructures

import (
	"math"
	"reflect"
	"testing"
)

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
	expectedNodes := []int{5, 7, 8, 3, 0, 6, 11, 0, 9, 0, 0, 7}
	expectedStructure := []byte{0b11101101, 0b00100000, 0, 0}
	tls := []TreeLeaf{
		{
			Values: []int{5, 7, 3, 9},
			Path:   []bool{false, false, true}},
		{
			Values: []int{5, 7, 3, 9},
			Path:   []bool{false, false, true}},
		{
			Values: []int{5, 8, 6, 7},
			Path:   []bool{true, false, false}},
		{
			Values: []int{5, 8, 11},
			Path:   []bool{true, true}}}
	fbt := MakeFBTreeFromLeafs(tls)
	if !reflect.DeepEqual(expectedNodes, fbt.nodes[:12]) {
		t.Errorf("node layout is not good, expected %v \n got %v", expectedNodes, fbt.nodes)
	}
	if !reflect.DeepEqual(fbt.structure.bits, expectedStructure) {
		t.Errorf("node structure is not good, expected %v \n got %v", expectedStructure, fbt.structure.bits)
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
	tls := []TreeLeaf{
		{
			Values: []int{2, 4, 6},
			Path:   []bool{true, false}},
		{
			Values: []int{2, 3, 5, 7},
			Path:   []bool{false, false, true}}}
	if maxDepth(tls) != 4 {
		t.Errorf("Max Depth failed, expecting %v, got %v.", 4, maxDepth(tls))
	}
}
