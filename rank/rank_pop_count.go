package rank

import (
	mbits "math/bits"

	"zs-project.org/aeroview/datastructures/bits"
)

const (
	// BLOCKSIZE the size of the block to use in the RankPopCount.
	BLOCKSIZE = 64
)

// PopCount is a data structure that helps compute rank and select on the bit vector BitVec.
type PopCount struct {
	metadata
	SuperBlockRanks []uint64
	Blocks          []uint16
	Data            bits.Vector
}

// MakePopCount creates a RankPopCount instance.
func MakePopCount(b bits.Vector) PopCount {
	// Blocksize is 64 bits for mecanichal sympathy.
	rm := makeRankMetadata(BLOCKSIZE, len(b)*64)

	rk := PopCount{
		metadata:        rm,
		SuperBlockRanks: make([]uint64, rm.NbSuperBlocks),
		Blocks:          make([]uint16, rm.NbBlocks),
		Data:            b,
	}
	cum := uint64(0)
	diff := uint16(0)
	for superBlockIdx := range rk.SuperBlockRanks {
		rk.SuperBlockRanks[superBlockIdx] = cum
		lower, upper := rk.blocksIdxForSuperBlock(superBlockIdx)
		for blockIdx := lower; blockIdx <= upper; blockIdx++ {
			d := rk.Data[blockIdx]
			rk.Blocks[blockIdx] = diff
			diff += uint16(mbits.OnesCount64(d))
			cum += uint64(mbits.OnesCount64(d))
		}
		diff = 0
	}
	return rk
}

// Rank ruturns the number of 1 bits in the bitvector for the first idx bits.
func (r PopCount) Rank(idx int) int {
	spblocIdx := idx / r.metadata.SuperBlockSize
	blockIdx := idx / BLOCKSIZE
	rankSuperBlock := r.SuperBlockRanks[spblocIdx]
	blockRank := uint64(r.Blocks[blockIdx])
	shift := idx % BLOCKSIZE
	mask := uint64(1<<(shift+1) - 1)
	pop := uint64(mbits.OnesCount64(r.Data[blockIdx] & mask))
	return int(rankSuperBlock + blockRank + pop)
}

// Select return the idx of the i'th one in the underlying bit vector.
func (r PopCount) Select(idx uint64) uint64 {
	spBlock := r.identifySuperBlock(idx)
	spBlockRank := r.SuperBlockRanks[spBlock]
	if spBlockRank >= idx && spBlock > 0 {
		spBlock--
	}
	spBlockRank = r.SuperBlockRanks[spBlock]
	lower, upper := r.blocksIdxForSuperBlock(spBlock)
	blockIdx := r.identifyBlock(idx, spBlockRank, lower, upper) + lower
	blockDiffRank := uint64(r.Blocks[blockIdx])
	if spBlockRank+blockDiffRank >= idx && blockIdx > 0 {
		blockIdx--
	}
	blockDiffRank = uint64(r.Blocks[blockIdx])
	d := r.Data[blockIdx]
	bDiffRank := int(blockDiffRank + spBlockRank)
	mask := uint64(1)<<(1) - 1
	for i := 0; i < r.BlockSize; i++ {
		dr := mbits.OnesCount64(d & mask)
		if dr+bDiffRank == int(idx) {
			return uint64(blockIdx*r.BlockSize + i)
		}
		mask = mask<<1 + 0b1
	}
	return uint64(0)
}

func (r PopCount) identifySuperBlock(i uint64) int {
	sblocks := r.SuperBlockRanks
	hi := len(sblocks) - 1
	lo := 0
	pos := (hi - lo) / 2
	for hi > lo {
		if data := sblocks[pos]; data < i {
			lo = pos + 1
			pos = (hi-lo)/2 + lo
		} else if data > i {
			hi = pos - 1
			pos = (hi-lo)/2 + lo
		} else {
			return pos
		}
	}
	return pos
}

func (r PopCount) identifyBlock(i, supBlockValue uint64, lowerBlockIdx, upperBlockIdx int) int {
	diff := i - supBlockValue
	blocks := r.Blocks[lowerBlockIdx:upperBlockIdx]
	hi := upperBlockIdx - lowerBlockIdx
	lo := 0
	pos := (hi - lo) / 2
	for hi > lo {
		if v := uint64(blocks[pos]); v < diff {
			lo = pos + 1
			pos = (hi-lo)/2 + lo
		} else if v > diff {
			hi = pos - 1
			pos = (hi-lo)/2 + lo
		} else {
			for i := pos - 1; i >= 0; i-- {
				if blocks[pos] != blocks[i] {
					return i + 1
				}
			}
		}
	}
	return pos
}

func (r PopCount) Get(idx int) int {
	return r.Data.Get(idx)
}
