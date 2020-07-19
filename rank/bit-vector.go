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

//func (b *BitVec) GetRange(low, high int) uint64 {
//	if high - low >= blockSize {
//		panic("cannot do more than 64 bit range")
//	}
//
//	lowIndex := low / blockSize
//	hiIndex := high / blockSize
//	lowOffset := low % blockSize
//	hiOffset := high % blockSize
//
//	if lowIndex == hiIndex {
//		return (*b)[lowIndex] >> lowOffset
//	}
//
//	lsbMask := (1 << hiOffset) - 1
//	lsb := (*b)[hiIndex] & (uint64(lsbMask) << (blockSize - hiOffset))
//
//	msb := (*b)[lowIndex] & ((1 << (blockSize - lowOffset)) - 1)
//	return uint64(msb << hiOffset) | lsb
//}



