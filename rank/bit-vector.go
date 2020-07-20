package rank

type BitVec []uint64

const (
	blockSize = 64
)

func NewBitVec(size int) *BitVec {
	vec := make(BitVec, size)
	return &vec
}

func (b *BitVec) Set(idx int, val uint8) {
	index := idx / blockSize
	offset := idx % blockSize
	if index >= len(*b) {
		*b = append(*b, make([]uint64, index - len(*b) + 1)...)
	}
	(*b)[index] |= uint64(val & 1) << offset
}

func (b *BitVec) Get(idx int) int {
	index := idx / blockSize
	offset := idx % blockSize
	if index >= len(*b) {
		return 0
	}
	return int((*b)[index] >> offset) & 1
}

func (b *BitVec) Get8BitRange(low, high int) uint8 {
	if high - low >= 8 {
		panic("cannot do more than 64 bit range")
	}
	if low == high {
		return uint8(b.Get(low))
	}

	lowIndex := low / blockSize
	hiIndex := high / blockSize
	lowOffset := low % blockSize
	hiOffset := high % blockSize

	if lowIndex == hiIndex {
		selectionMask := uint64((1 << (hiOffset - lowOffset + 1)) - 1)
		return uint8(((*b)[lowIndex] >> lowOffset) & selectionMask)
	}

	// lsb
	lsbHi := blockSize - lowOffset - 1
	lsb := b.Get8BitRange(low, low + lsbHi)

	// msb
	msbLow := high - hiOffset
	msb := b.Get8BitRange(msbLow, high)

	// shift msb by the size of lsb
	return lsb | (msb << (blockSize - lowOffset))
}
