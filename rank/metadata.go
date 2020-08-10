package rank

type metadata struct {
	SuperBlockSize      int
	NbSuperBlocks       int
	BlocksPerSuperBlock int
	NbBlocks            int
	BlockSize           int
	NbBits              int
}

func makeRankMetadata(blockSize int, nbBits int) metadata {
	nbBlocks := ceil_devide(nbBits, blockSize)
	nbSuperBlocks := valueOrOne(floor_devide(nbBits, blockSize*nbBits_floor(nbBits)))
	superBlockSize := (nbBits / nbSuperBlocks) - (nbBits/nbSuperBlocks)%blockSize
	blocksPerSuperBlock := nbBlocks / nbSuperBlocks
	return metadata{
		SuperBlockSize:      superBlockSize,
		NbSuperBlocks:       nbSuperBlocks + 1,
		BlocksPerSuperBlock: blocksPerSuperBlock,
		NbBlocks:            nbBlocks,
		BlockSize:           blockSize,
		NbBits:              nbBits,
	}
}

func (r metadata) Overhead() float64 {
	sizeOfSuperBlocks := r.NbSuperBlocks * 64
	sizeOfBlocks := r.NbBlocks * 8
	sizeOfMetadata := 6 * 64
	sizeOfData := r.NbBits
	return float64(sizeOfMetadata+sizeOfData+sizeOfSuperBlocks+sizeOfBlocks) / float64(sizeOfData)
}

// TODO refactor to only return min and max ( saving the allocation )
func (r metadata) blocksIdxForSuperBlock(i int) (int, int) {
	if i < 0 {
		return 0, 0
	}
	lower := i * r.BlocksPerSuperBlock
	upper := lower + r.BlocksPerSuperBlock
	if upper > r.NbBlocks {
		upper = r.NbBlocks
	}
	return lower, upper - 1
}

func valueOrOne(i int) int {
	if i <= 0 {
		return 1
	}
	return i
}
