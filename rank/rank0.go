package rank

import "github.com/zs-projects/aeroview/datastructures/bits"

const (
	smallBlockSize = 8
	bigBlockSize   = 64
)

type Rank0 struct {
	bits   bits.Vector
	data   []uint64
	blocks []int
	lookup [256]uint8
}

func (r *Rank0) Get(idx int) uint64 {
	if r.bits.Get(idx) == 0 {
		return 0
	}
	var sum int
	bigBlockIndex := idx / bigBlockSize
	blockOffset := idx % bigBlockSize

	// 1. get big block
	if bigBlockIndex > 0 {
		sum += r.blocks[bigBlockIndex-1]
	}

	// 2. get sum from lookup
	sum += r.rank(idx-blockOffset, idx)

	return r.data[sum-1]
}

func (r *Rank0) rank(low, idx int) int {
	runningRank := 0
	for i := low; i <= idx; {
		var n uint8
		if i+7 <= idx {
			n = r.bits.Get8BitRange(i, i+7)
		} else {
			n = r.bits.Get8BitRange(i, idx)
		}
		runningRank += int(r.lookup[n])
		i += 8
	}
	return runningRank
}

func NewRank0(xs []uint64) *Rank0 {
	lookup := bits.Make8BitLookup()

	nBigBlocks := len(xs) / 64
	blocks := make([]int, nBigBlocks+1)
	bits := bits.NewVector(nBigBlocks)

	var data []uint64
	var s int
	for i, x := range xs {

		if x != 0 {
			data = append(data, x)
			s++
			blocks[i/bigBlockSize] = s
			bits.Set(i, 1)
		}
	}

	return &Rank0{
		bits:   *bits,
		data:   data,
		blocks: blocks,
		lookup: lookup,
	}
}
