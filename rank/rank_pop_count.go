package rank

import (
	"math/bits"
)

const (
	BLOCKSIZE = 64
)

type rankMetadata struct {
	SuperBlockSize      int
	NbSuperBlocks       int
	BlocksPerSuperBlock int
	NbBlocks            int
	BlockSize           int
	NbBits              int
}

func makeRankMetadata(blockSize int, nbBits int) rankMetadata {
	nbBlocks := ceil_devide(nbBits, blockSize)
	nbSuperBlocks := floor_devide(nbBits, blockSize*nbBits_floor(nbBits))
	superBlockSize := (nbBits / nbSuperBlocks) - (nbBits/nbSuperBlocks)%blockSize
	blocksPerSuperBlock := nbBlocks / nbSuperBlocks
	return rankMetadata{
		SuperBlockSize:      superBlockSize,
		NbSuperBlocks:       nbSuperBlocks + 1,
		BlocksPerSuperBlock: blocksPerSuperBlock,
		NbBlocks:            nbBlocks,
		BlockSize:           blockSize,
		NbBits:              nbBits,
	}
}

func (r rankMetadata) Overhead() float64 {
	sizeOfSuperBlocks := r.NbSuperBlocks * 64
	sizeOfBlocks := r.NbBlocks * 8
	sizeOfMetadata := 6 * 64
	sizeOfData := r.NbBits
	return float64(sizeOfMetadata+sizeOfData+sizeOfSuperBlocks+sizeOfBlocks) / float64(sizeOfData)
}

type RankPopCount struct {
	rankMetadata
	SuperBlockRanks []uint64
	Blocks          []uint8
	data            BitVec
}

func (r RankPopCount) BlocksIdxForSuperBlock(i int) []int {
	if i < 0 {
		return nil
	}
	lower := i * r.rankMetadata.BlocksPerSuperBlock
	upper := lower + r.rankMetadata.BlocksPerSuperBlock
	if upper > len(r.Blocks) {
		upper = len(r.Blocks)
	}
	ret := make([]int, upper-lower)
	for i := range ret {
		ret[i] = lower + i
	}
	return ret
}
func (r RankPopCount) Rank(idx int) int {
	spblocIdx := idx / r.rankMetadata.SuperBlockSize
	blockIdx := idx / BLOCKSIZE
	shift := BLOCKSIZE - idx%BLOCKSIZE
	rankSuperBlock := r.SuperBlockRanks[spblocIdx]
	blockRank := uint64(r.Blocks[blockIdx])
	dataIdx := idx / BLOCKSIZE
	pop := uint64(bits.OnesCount64(r.data[dataIdx] >> shift))
	return int(rankSuperBlock + blockRank + pop)
}

func MakeRankPopCount(b BitVec) RankPopCount {
	// Blocksize is 64 bits for mecanichal sympathy.
	rm := makeRankMetadata(blockSize, len(b)*64)

	rk := RankPopCount{
		rankMetadata:    rm,
		SuperBlockRanks: make([]uint64, rm.NbSuperBlocks),
		Blocks:          make([]uint8, rm.NbBlocks),
		data:            b,
	}
	cum := uint64(0)
	diff := uint8(0)
	for superBlockIdx := range rk.SuperBlockRanks {
		rk.SuperBlockRanks[superBlockIdx] = cum
		for _, blockIdx := range rk.BlocksIdxForSuperBlock(superBlockIdx) {
			d := rk.data[blockIdx]
			rk.Blocks[blockIdx] = diff
			diff += uint8(bits.OnesCount64(d))
			cum += uint64(bits.OnesCount64(d))
		}
		diff = 0
	}
	return rk
}
