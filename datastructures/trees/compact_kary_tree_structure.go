package trees

import (
	"github.com/zs-projects/aeroview/datastructures/bits"
	"github.com/zs-projects/aeroview/rank"
)

// CompactKAryTreeStructure encode the structure of a k-ary tree
// in a bitset.
type CompactKAryTreeStructure struct {
	// The maximum number of children for a node in the tree.
	K NodeNBChildren
	rank.GetRanker
}

// MakeCompactKAryTreeStructure make a bitset k-ary tree structure for a level order traversal of the tree.
func MakeCompactKAryTreeStructure(tr KAryLevelOrder) CompactKAryTreeStructure {
	maxChildren := NodeNBChildren(0)
	for _, v := range tr {
		if maxChildren <= v {
			maxChildren = v
		}
	}

	q := bits.MakeQueue()
	for _, nbChildren := range tr {
		for i := 0; i < int(nbChildren); i++ {
			q.PushBack(1)
		}
		if nbChildren < maxChildren {
			for i := nbChildren; i < maxChildren; i++ {
				q.PushBack(0)
			}
		}
	}
	return CompactKAryTreeStructure{
		K:         maxChildren,
		GetRanker: rank.MakePopCount(q.Vector()),
	}
}

// Children returns the position of the children of the provided node and the number of children the
// provided node has.
// if nbChildren is 0, then startPosition is invalid and should be discarded.
// if nbChildren is greater than 0, then startPosition is the position of the first child.
func (t CompactKAryTreeStructure) Children(nodePos int) (nbChildren, startPosition int) {
	if t.Get(t.K*nodePos) == 0 {
		// No children return nil
		return 0, 0
	}
	bitPos := t.K * nodePos
	startPosition = t.Rank(t.K * nodePos)
	nbChildren = 0
	for i := 0; i < t.K; i++ {
		if t.Get(bitPos+i) == 0 {
			return nbChildren, startPosition
		}
		nbChildren++
	}
	return nbChildren, startPosition
}
