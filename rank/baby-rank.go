package rank

import "zs-project.org/aeroview/datastructures/bits"

type BabyRank struct {
	bits   *bits.Vector
	data   []uint64
	lookup [256]uint8
}

func NewBabyRank(xs []uint64) *BabyRank {

	lookup := bits.Make8BitLookup()

	nBigBlocks := len(xs) / 64
	bits := bits.NewVector(nBigBlocks)

	var data []uint64
	for i, x := range xs {

		if x != 0 {
			data = append(data, x)
			bits.Set(i, 1)
		}
	}

	return &BabyRank{
		bits:   bits,
		data:   data,
		lookup: lookup,
	}
}

func (r *BabyRank) Get(idx int) uint64 {
	if r.bits.Get(idx) == 0 {
		return 0
	}
	return r.data[r.rank(idx)-1]
}

func (r *BabyRank) rank(idx int) int {
	runningRank := 0
	for i := 0; i <= idx; {
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
