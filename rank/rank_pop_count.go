package rank

import (
	"math/bits"
)

const (
	// BLOCKSIZE the size of the block to use in the RankPopCount.
	BLOCKSIZE = 64
)

// PopCount is a data structure that helps compute rank and select on the bit vector BitVec.
type PopCount struct {
	metadata
	SuperBlockRanks []uint64
	Blocks          []uint8
	data            BitVec
}

// MakePopCount creates a RankPopCount instance.
func MakePopCount(b BitVec) PopCount {
	// Blocksize is 64 bits for mecanichal sympathy.
	rm := makeRankMetadata(blockSize, len(b)*64)

	rk := PopCount{
		metadata:        rm,
		SuperBlockRanks: make([]uint64, rm.NbSuperBlocks),
		Blocks:          make([]uint8, rm.NbBlocks),
		data:            b,
	}
	cum := uint64(0)
	diff := uint8(0)
	for superBlockIdx := range rk.SuperBlockRanks {
		rk.SuperBlockRanks[superBlockIdx] = cum
		lower, upper := rk.blocksIdxForSuperBlock(superBlockIdx)
		for blockIdx := lower; blockIdx <= upper; blockIdx++ {
			d := rk.data[blockIdx]
			rk.Blocks[blockIdx] = diff
			diff += uint8(bits.OnesCount64(d))
			cum += uint64(bits.OnesCount64(d))
		}
		diff = 0
	}
	return rk
}

// Rank ruturns the number of 1 bits in the bitvector for the first idx bits.
func (r PopCount) Rank(idx int) int {
	spblocIdx := idx / r.metadata.SuperBlockSize
	blockIdx := idx / BLOCKSIZE
	shift := BLOCKSIZE - idx%BLOCKSIZE - 1
	rankSuperBlock := r.SuperBlockRanks[spblocIdx]
	blockRank := uint64(r.Blocks[blockIdx])
	dataIdx := idx / BLOCKSIZE
	pop := uint64(bits.OnesCount64(r.data[dataIdx] >> shift))
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
	d := r.data[blockIdx]
	bDiffRank := int(blockDiffRank + spBlockRank)
	for i := 63; i >= 0; i-- {
		if bits.OnesCount64(d>>i)+bDiffRank == int(idx) {
			return uint64(blockIdx*r.BlockSize+(r.BlockSize-i)) - 1
		}
	}
	return uint64(0)
}

func (r PopCount) identifySuperBlock(i uint64) int {
	sblocks := r.SuperBlockRanks
	hi := len(sblocks) - 1
	lo := 0
	pos := (hi - lo) / 2
	for hi > lo {
		if sblocks[pos] < i {
			lo = pos + 1
			pos = (hi-lo)/2 + lo
		} else if sblocks[pos] > i {
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
