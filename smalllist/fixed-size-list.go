package smalllist

const blockSize = 64

type FixedSized struct {
	size      uint64
	smalllist []uint64
}

func FromSlice(xs []int) *FixedSized {
	return nil
}

func (f *FixedSized) Get(idx int) uint64 {
	// find block and offset.
	block, offset := f.blockAndOffset(idx)

	// if overflow
	var x uint64
	if offset+f.size > blockSize {
		// size of the overflow
		overflow := (offset + f.size) % blockSize
		// select of msb of size  (blockSize - offset)
		msb := (f.smalllist[block] >> offset) << overflow
		// select of lsb of size of overflow in next block
		lsb := selectLastKBits(f.smalllist[block+1], overflow)
		x |= msb
		x |= lsb
	} else {
		x = selectLastKBits(f.smalllist[block]>>offset, f.size)
	}
	return x
}

func (f *FixedSized) Set(val uint64, idx int) {
	// find block and offset.
	block, offset := f.blockAndOffset(idx)

	// if overflow
	if offset+f.size > blockSize {
		// size of the overflow
		overflow := (offset + f.size) % blockSize
		// most significant bits only, remove overflow bits
		msb := val >> (overflow)
		// least significant bits only
		lsb := selectLastKBits(val, overflow)
		f.smalllist[block] |= msb << offset
		f.smalllist[block+1] |= lsb
	} else {
		f.smalllist[block] |= val << offset
	}
}

func (f *FixedSized) blockAndOffset(idx int) (uint64, uint64) {
	block := (uint64(idx) * f.size) / blockSize
	offset := (uint64(idx) * f.size) % blockSize
	return block, offset
}

func selectLastKBits(val uint64, k uint64) uint64 {
	return val & ((1 << k) - 1)
}

func selectFirstKBits(val uint64, k uint64) uint64 {
	mask := (1 << k) - 1
	return val & (uint64(mask) << (uint64(64) - k))
}
