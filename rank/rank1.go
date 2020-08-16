package rank

import "github.com/zs-projects/aeroview/datastructures/bits"

type Rank1 struct {
	bits        bits.Vector
	data        []uint64
	bigBlocks   []int
	smallBlocks []uint8
	lookup      [256]uint8
}

func (r *Rank1) Get(idx int) uint64 {
	if r.bits.Get(idx) == 0 {
		return 0
	}
	var sum int

	// 1. get big block
	bigBlockIndex := idx / bigBlockSize
	if bigBlockIndex > 0 {
		sum += r.bigBlocks[bigBlockIndex-1]
	}

	// 2. get small block sum
	smallBlockIndex := idx / smallBlockSize
	smallBlockOffset := idx % smallBlockSize
	sum += int(r.smallBlocks[smallBlockIndex])

	// 3. get sum from lookup
	lookupIndex := r.bits.Get8BitRange(idx-smallBlockOffset, idx)
	sum += int(r.lookup[lookupIndex])

	return r.data[sum-1]
}

func NewRank1(xs []uint64) *Rank1 {
	lookup := bits.Make8BitLookup()

	nBigBlocks := len(xs) / bigBlockSize
	bigBlocks := make([]int, nBigBlocks+1)

	nSmallBlocks := len(xs) / smallBlockSize
	// +2 because front and back padding
	smallBlocks := make([]uint8, nSmallBlocks+2)

	bits := bits.NewVector(nBigBlocks)

	var data []uint64
	var bigBlockSum int
	var smallBlockSum uint8

	for i, x := range xs {

		// keep running sum only within a big block
		if i%bigBlockSize == 0 {
			smallBlockSum = 0
		}

		if x != 0 {
			// 1. data: keep non zeroed data
			data = append(data, x)
			bits.Set(i, 1)

			// 2. big block sum
			blockIndex := i / bigBlockSize
			bigBlockSum++
			bigBlocks[blockIndex] = bigBlockSum

			// 3. small block
			smallBlockSum++
		}

		smallBlockIndex := i / smallBlockSize
		if smallBlockIndex%8 != 7 {
			smallBlocks[smallBlockIndex+1] = smallBlockSum
		}
	}

	return &Rank1{
		bits:        *bits,
		data:        data,
		bigBlocks:   bigBlocks,
		smallBlocks: smallBlocks,
		lookup:      lookup,
	}
}
