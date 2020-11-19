package rank

import (
	mbits "math/bits"

	"github.com/zs-projects/aeroview/datastructures/bits"
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

func (r *PopCount) SizeInBytes() int {
	return len(r.SuperBlockRanks)*8 + len(r.Blocks)*2 + len(r.Data)*8 + 6*8
}

// MakePopCount creates a RankPopCount instance.
func MakePopCount(b bits.Vector) *PopCount {
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
	return &rk
}

// Rank ruturns the number of 1 bits in the bitvector for the first idx bits.
func (r *PopCount) Rank(idx int) int {
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
func (r *PopCount) Select(idx uint64) uint64 {
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

func (r *PopCount) identifySuperBlock(i uint64) int {
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

func (r *PopCount) identifyBlock(i, supBlockValue uint64, lowerBlockIdx, upperBlockIdx int) int {
	diff := i - supBlockValue
	for idx := lowerBlockIdx; idx <= upperBlockIdx; idx++ {
		if v := uint64(r.Blocks[idx]); v >= diff {
			return idx - lowerBlockIdx
		}
	}
	return upperBlockIdx - lowerBlockIdx
}

func (r *PopCount) Get(idx int) uint64 {
	return r.Data.Get(idx)
}
