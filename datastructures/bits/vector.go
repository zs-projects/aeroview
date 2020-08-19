package bits

// BLOCKSIZE is the size in bits of the type used as an underlying container for our BitVec.
const BLOCKSIZE = 64

// Vector represent a vector of bits.
// bits are represented from right to left ( the low bit of the first uint64 is the first bit of the Vector )
type Vector []uint64

func NewVector(size int) *Vector {
	vec := make(Vector, size)
	return &vec
}

func (b *Vector) Set(idx int, val uint8) {
	index := idx / BLOCKSIZE
	offset := idx % BLOCKSIZE
	if index >= len(*b) {
		*b = append(*b, make([]uint64, index-len(*b)+1)...)
	}
	(*b)[index] |= uint64(val&1) << offset
}

func (b Vector) Get(idx int) uint64 {
	index := idx / BLOCKSIZE
	offset := idx % BLOCKSIZE
	if index >= len(b) {
		return 0
	}
	return (b)[index] >> offset & 1
}

func (b *Vector) Get8BitRange(low, high int) uint8 {
	if high-low >= 8 {
		panic("cannot do more than 64 bit range")
	}
	if low == high {
		return uint8(b.Get(low))
	}

	lowIndex := low / BLOCKSIZE
	hiIndex := high / BLOCKSIZE
	lowOffset := low % BLOCKSIZE
	hiOffset := high % BLOCKSIZE

	if lowIndex == hiIndex {
		selectionMask := uint64((1 << (hiOffset - lowOffset + 1)) - 1)
		return uint8(((*b)[lowIndex] >> lowOffset) & selectionMask)
	}

	// lsb
	lsbHi := BLOCKSIZE - lowOffset - 1
	lsb := b.Get8BitRange(low, low+lsbHi)

	// msb
	msbLow := high - hiOffset
	msb := b.Get8BitRange(msbLow, high)

	// shift msb by the size of lsb
	return lsb | (msb << (BLOCKSIZE - lowOffset))
}
