package recsplit

type Node struct {
	Left   *Node
	Right  *Node
	R      int
	nbKeys int
}

func FromCompactFBTree(fbt CompactFBTree) *Node {
	return fromCompactFBTree(fbt, 0)
}

func fromCompactFBTree(fbt CompactFBTree, currOffset int) *Node {
	var (
		left  *Node = nil
		right *Node = nil
	)
	if fbt.nodeHasLeftChild(currOffset) {
		lOffset := fbt.LeftChild(currOffset)
		left = fromCompactFBTree(fbt, lOffset)
	}
	if fbt.nodeHasRightChild(currOffset) {
		rOffset := fbt.RightChild(currOffset)
		right = fromCompactFBTree(fbt, rOffset)
	}
	R, nbkeys := fbt.node(currOffset)
	return &Node{
		Left:   left,
		Right:  right,
		R:      R,
		nbKeys: nbkeys,
	}
}
