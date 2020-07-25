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
	superBlockRanks []uint64
	blocks          []uint8
	data            BitVec
}

func (r RankPopCount) blocksIdxForSuperBlock(i int) []int {
	if i < 0 {
		return nil
	}
	lower := i * r.rankMetadata.BlocksPerSuperBlock
	upper := lower + r.rankMetadata.BlocksPerSuperBlock
	if upper > len(r.blocks) {
		upper = len(r.blocks)
	}
	ret := make([]int, upper-lower)
	for i := range ret {
		ret[i] = lower + i
	}
	return ret
}
func (r RankPopCount) Rank(idx int) uint64 {
	spblocIdx := idx / r.rankMetadata.SuperBlockSize
	blockIdx := idx / BLOCKSIZE
	shift := BLOCKSIZE - idx%BLOCKSIZE
	rankSuperBlock := r.superBlockRanks[spblocIdx]
	blockRank := uint64(r.blocks[blockIdx])
	dataIdx := idx / BLOCKSIZE
	pop := uint64(bits.OnesCount64(r.data[dataIdx] >> shift))
	return rankSuperBlock + blockRank + pop
}

func MakeRankPopCount(b BitVec) RankPopCount {
	// Blocksize is 64 bits for mecanichal sympathy.
	rm := makeRankMetadata(blockSize, len(b)*64)

	rk := RankPopCount{
		rankMetadata:    rm,
		superBlockRanks: make([]uint64, rm.NbSuperBlocks),
		blocks:          make([]uint8, rm.NbBlocks),
		data:            b,
	}
	cum := uint64(0)
	diff := uint8(0)
	for superBlockIdx := range rk.superBlockRanks {
		rk.superBlockRanks[superBlockIdx] = cum
		for _, blockIdx := range rk.blocksIdxForSuperBlock(superBlockIdx) {
			d := rk.data[blockIdx]
			rk.blocks[blockIdx] = diff
			diff += uint8(bits.OnesCount64(d))
			cum += uint64(bits.OnesCount64(d))
		}
		diff = 0
	}
	return rk
}
