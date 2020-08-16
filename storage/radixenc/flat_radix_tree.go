package radixenc

import (
	"strings"

	"github.com/zs-projects/aeroview/datastructures/bits"
	"github.com/zs-projects/aeroview/datastructures/trees"
	"github.com/zs-projects/aeroview/encoding"
)

type FlatRadixTree struct {
	data         []rune
	offsetsStart encoding.EliasFanoVector
	structure    trees.KAryTreeStructure
	leafs        bits.Queue
	maxChildren  int
}
type RadixLevelOrder []int

func (u RadixLevelOrder) LevelOrder() []int {
	return []int(u)
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
	structure := make([]int, 0)
	bitQueueLeafs := bits.MakeQueue()
	maxChildren := 0
	for cur < len(queueCh) {
		curNode := queueCh[cur]
		curOffset := len(data)
		offsetsStart = append(offsetsStart, uint64(curOffset))
		structure = append(structure, len(curNode.children))
		if curNode.isLeaf {
			bitQueueLeafs.PushBack(1)
		} else {
			bitQueueLeafs.PushBack(0)
		}
		for k, v := range curNode.children {
			queueCh = append(queueCh, v)
			queuePr = append(queuePr, k)
		}
		if maxChildren < len(curNode.children) {
			maxChildren = len(curNode.children)
		}
		cur++
	}

	return FlatRadixTree{
		data:         data,
		offsetsStart: encoding.MakeEliasFanoVector(offsetsStart),
		structure:    trees.MakeCompactKAryTreeStructure(RadixLevelOrder(structure)),
		leafs:        bitQueueLeafs,
		maxChildren:  maxChildren,
	}
}

func (f FlatRadixTree) children(nodeIdx int) (nbChildren, startPos int) {
	// TODO: Implement the encoding function.
	return f.structure.Children(nodeIdx)
}

func (f FlatRadixTree) Encode(data []string) [][]int {
	// TODO: Implement the encoding function.
	ret := make([][]int, 0, len(data))
	node := -1 // For root node
	for _, s := range data {
		cur := 0
		encoding := make([]int, 0)
		for cur < len(s) {
			if nbChildren, startPos := f.children(node); nbChildren > 0 {
				for i := 0; i < nbChildren; i++ {
					childIdx := i + startPos
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
		}
		ret = append(ret, encoding)
	}
	return nil
}
