package selecter

import "zs-project.org/aeroview/rank"

type BinarySearchSelect struct {
	rank.PopCount
	Data rank.BitVec
}

func (bss BinarySearchSelect) IdentifySuperBlock(i uint64) uint64 {
	sblocks := bss.SuperBlockRanks
	hi := len(sblocks) - 1
	lo := 0
	pos := (hi - lo) / 2
	for hi > lo {
		if sblocks[pos] < i {
			hi = pos - 1
			pos = (hi - lo) / 2
		} else if sblocks[pos] > i {
			lo = pos + 1
			pos = (hi - lo) / 2
		} else {
			return uint64(pos)
		}
	}
	return uint64(pos)
}
func (bss BinarySearchSelect) Select(i uint64) uint64 {
	sbBlock := bss.IdentifySuperBlock(i)
	if bss.SuperBlockRanks[sbBlock] == i {
		return sbBlock * uint64(bss.SuperBlockSize)
	}
	// Question: do we perform a linear scan on the blocks?.
	// Question: Do we do w linear scan on the blocks.
	return uint64(0)
}
