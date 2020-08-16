package radixenc

import (
	"strings"

	"github.com/zs-projects/aeroview/datastructures/bits"
	"github.com/zs-projects/aeroview/encoding"
	"github.com/zs-projects/aeroview/rank"
)

type FlatRadixTree struct {
	data         []rune
	offsetsStart encoding.EliasFanoVector
	structure    rank.PopCount
	leafs        bits.Queue
	maxChildren  int
}

func MakeFlatRadixTree(r RadixTree) FlatRadixTree {
	queueCh := make([]*RadixTree, 0)
	queuePr := make([]string, 0)
	for k, v := range r.children {
		queueCh = append(queueCh, v)
		queuePr = append(queuePr, k)
	}
	cur := 0

	data := make([]rune, 0)
	offsetsStart := make([]uint64, 0)
	offsetsStop := make([]uint64, 0)
	structure := make([]int, 0)
	leafs := make([]bool, 0)

	for cur < len(queueCh) {
		curNode := queueCh[cur]
		curOffset := len(data)
		curPrefix := queuePr[cur]
		offsetsStart = append(offsetsStart, uint64(curOffset))
		offsetsStop = append(offsetsStop, uint64(curOffset+len(curPrefix)))
		structure = append(structure, len(curNode.children))
		leafs = append(leafs, curNode.isLeaf)
		for k, v := range curNode.children {
			queueCh = append(queueCh, v)
			queuePr = append(queuePr, k)
		}
		cur++
	}

	maxChildren := MaxNonEmptyIntSlice(structure)
	bitQueueStructure := bits.MakeQueue()
	bitQueueLeafs := bits.MakeQueue()
	for idx, nbChilds := range structure {
		for i := 0; i < nbChilds; i++ {
			bitQueueStructure.PushBack(1)
			if leafs[idx] {
				bitQueueLeafs.PushBack(1)
			} else {
				bitQueueLeafs.PushBack(0)
			}
		}
		for i := nbChilds; i < maxChildren-nbChilds; i++ {
			bitQueueStructure.PushBack(0)
			bitQueueLeafs.PushBack(0)
		}
	}
	return FlatRadixTree{
		data:         data,
		offsetsStart: encoding.MakeEliasFanoVector(offsetsStart),
		structure:    rank.MakePopCount(bitQueueStructure.Vector()),
		leafs:        bitQueueLeafs,
		maxChildren:  maxChildren,
	}
}

func (f FlatRadixTree) children(nodeIdx int) []int {
	// TODO: Implement the encoding function.
	var start int
	if nodeIdx < 0 {
		start = 0
	} else {
		start = f.maxChildren * f.structure.Rank(nodeIdx)
	}
	children := make([]int, 0, f.maxChildren)
	for i := start; i < start+f.maxChildren; i++ {
		if f.structure.Get(i) == 0 {
			return children
		}
		children = append(children, i)
	}
	return children
}

func (f FlatRadixTree) Encode(data []string) [][]int {
	// TODO: Implement the encoding function.
	ret := make([][]int, 0, len(data))
	node := -1 // For root node
	for _, s := range data {
		cur := 0
		encoding := make([]int, 0)
		for cur < len(s) {
			for childIdx := range f.children(node) {
				start := f.offsetsStart.Get(childIdx)
				stop := f.offsetsStart.Get(childIdx+1) - 1
				prefix := string(f.data[start:stop])
				if strings.HasPrefix(s[cur:], prefix) {
					encoding = append(encoding, childIdx)
					node = childIdx
					cur += len(prefix)
					break
				}
			}
		}
		ret = append(ret, encoding)
	}
	return nil
}

func MaxNonEmptyIntSlice(s []int) int {
	if len(s) != 0 {
		m := s[0]
		for _, v := range s {
			if v > m {
				m = v
			}
		}
		return m
	}
	panic("empty array")
}
