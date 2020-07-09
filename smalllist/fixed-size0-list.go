package smalllist

type IntSize int32

const (
	length8  IntSize = 8
	length16 IntSize = 16
	length32 IntSize = 32
)

type IntMask int64

const (
	bit8Mask  IntMask = 0xFF
	bit16Mask IntMask = 0xFFFF
	bit32Mask IntMask = 0xFFFFFFFF
)

type FixedSizedList struct {
	list              []int64
	nElementsPerBlock int
	intSize           IntSize
	mask              IntMask
}

func From(xs []int, intSize IntSize) FixedSizedList {

	var mask IntMask
	switch intSize {
	case length8:
		mask = bit8Mask
	case length16:
		mask = bit16Mask
	case length32:
		mask = bit32Mask

	}

	nElementsPerBlock := 64 / int(intSize)
	newSize := ((len(xs) * int(intSize)) + 63) / 64
	fixedSize := make([]int64, newSize)
	for i, x := range xs {
		block := i / nElementsPerBlock
		positionInBlock := i % nElementsPerBlock
		newVal := int64(x << (positionInBlock * int(intSize)))
		fixedSize[block] |= newVal
	}
	return FixedSizedList{
		list:              fixedSize,
		nElementsPerBlock: nElementsPerBlock,
		intSize:           intSize,
		mask:              mask,
	}
}

func (f *FixedSizedList) Get(index int) int {
	newIndex := index / (f.nElementsPerBlock)
	positionInBlock := index % f.nElementsPerBlock
	element := f.list[newIndex] >> (positionInBlock * int(f.intSize))
	return int(element & int64(f.mask))
}

func (f *FixedSizedList) Set(index int, value int) {
	newIndex := index / (f.nElementsPerBlock)
	positionInBlock := index % f.nElementsPerBlock
	newVal := int64(value << (positionInBlock * int(f.intSize)))
	f.list[newIndex] |= newVal
}

