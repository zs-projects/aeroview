package radixenc

import (
	"encoding/binary"
	"strings"

	"github.com/zs-projects/aeroview/datastructures/bits"
	"github.com/zs-projects/aeroview/datastructures/trees"
)

type FlatRadixTree struct {
	data         []rune
	offsetsStart []uint64
	structure    trees.KAryTreeStructure
	leafs        bits.Vector
	maxChildren  int
}

func MakeFlatRadixTree(r RadixTree) FlatRadixTree {
	queueCh := make([]*RadixTree, 0)
	queuePr := make([]string, 0)
	queueCh = append(queueCh, &r)
	queuePr = append(queuePr, "")
	cur := 0

	data := make([]rune, 0)
	offsetsStart := make([]uint64, 1) // We want the offsets to start with 0 for regularity.
	structure := make([]int, 0)
	bitQueueLeafs := bits.MakeQueue()
	maxChildren := 0
	for cur < len(queueCh) {
		curNode := queueCh[cur]
		curData := queuePr[cur]
		data = append(data, []rune(curData)...)
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
		offsetsStart: offsetsStart,
		structure:    trees.MakeCompactKAryTreeStructure(structure),
		leafs:        bitQueueLeafs.Vector(),
		maxChildren:  maxChildren,
	}
}

func (f FlatRadixTree) Children(nodeIdx int) (nbChildren, startPos int) {
	// TODO: Implement the encoding function.
	return f.structure.Children(nodeIdx)
}

func (f FlatRadixTree) Size() int {
	return (len(f.data) * 4) + len(f.leafs)/8 + 4 + len(f.offsetsStart)/8 + (len(f.leafs) * f.maxChildren / 8)
}

func (f FlatRadixTree) Overhead() float64 {
	return float64(f.Size()) / float64(len(f.data)*4)
}

func (f FlatRadixTree) Encode(data []string) [][]int {
	// TODO: Implement the encoding function.
	ret := make([][]int, 0, len(data))
	node := 0 // For root node
	for _, s := range data {
		cur := 0
		encoding := make([]int, 0)
		node = 0
		for cur < len(s) {
			if nbChildren, startPos := f.Children(node); nbChildren > 0 {
				for i := 0; i < nbChildren; i++ {
					childIdx := i + startPos
					start := f.offsetsStart[childIdx]
					stop := f.offsetsStart[childIdx+1]
					prefix := string(f.data[start:stop])
					if strings.HasPrefix(s[cur:], prefix) {
						encoding = append(encoding, childIdx)
						node = childIdx
						cur += len(prefix)
						break
					}
				}
			} else if f.leafs.Get(node) == 1 {
				break
			}
		}
		ret = append(ret, encoding)
	}
	return ret
}

func (f FlatRadixTree) Decode(encodedData [][]int) []string {
	out := make([]string, 0, len(encodedData))
	for _, encodedStr := range encodedData {
		str := make([]rune, 0, 2048)
		for _, token := range encodedStr {
			str = append(str, f.data[f.offsetsStart[token]:f.offsetsStart[token+1]]...)
		}
		out = append(out, string(str))
	}
	return out
}

func (f FlatRadixTree) DecodeFast(encodedData []byte) []rune {
	out := make([]rune, 0, len(encodedData))
	for i := 0; i < len(encodedData); i += 4 {
		token := binary.LittleEndian.Uint32(encodedData[i : i+4])
		out = append(out, f.data[f.offsetsStart[token]:f.offsetsStart[token+1]]...)
	}
	return out
}
