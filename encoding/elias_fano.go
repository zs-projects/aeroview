package encoding

import (
	"math"

	"zs-project.org/aeroview/datastructures"
	"zs-project.org/aeroview/rank"
)

// EliasFanoVector encodes a list of ascending integers using the Elias Fano Code.
type EliasFanoVector struct {
	highBits      datastructures.BitQueue
	lowBits       datastructures.BitQueue
	rank          rank.Ranker
	nElements     int // the number of elements in the data structure.
	lowBitsCount  int // The number of bits used to encode the low bits.
	highBitsCount int // The number of bits used to encode the low bits.
}

// MakeEliasFanoVector encodes a list of uint64 using EliasFano Code.
func MakeEliasFanoVector(values []uint64) EliasFanoVector {
	lowBitsQ := datastructures.MakeBitQueue()
	highBitsQ := datastructures.MakeBitQueue()
	lowerBitCount := uint64(math.Log2(float64(len(values))))
	lowBitsMask := uint64(1)<<lowerBitCount - 1
	prev := uint64(0)
	for i := 0; i < len(values); i++ {
		v := values[i]
		highBits := v >> lowerBitCount
		lowBits := v & lowBitsMask
		highDelta := highBits - prev
		for k := lowerBitCount; k > 0; k-- {
			lowBitsQ.PushBack(uint8(lowBits >> (k - 1)))
		}
		if highDelta == 0 {
			highBitsQ.PushBack(1)
		} else {
			for j := uint64(0); j < highDelta; j++ {
				highBitsQ.PushBack(0)
			}
			highBitsQ.PushBack(1)
		}
		prev = highBits
	}

	return EliasFanoVector{
		highBits:     highBitsQ,
		lowBits:      lowBitsQ,
		nElements:    len(values),
		lowBitsCount: int(lowerBitCount),
	}
}

// Len returns the number of bits set in the vector.
func (e EliasFanoVector) Len() int {
	return e.highBits.Len() + e.lowBits.Len()
}

// Data returns the raw elias-fano encoded array.
func (e EliasFanoVector) Data() ([]byte, int) {
	result, _ := datastructures.MakeBitQueueFromSlice(e.highBits.Data(), e.highBits.Len())
	result.Append(e.lowBits.Data(), e.lowBits.Len())
	return result.Data(), result.Len()
}

// Get Returns the element at index i
func (e EliasFanoVector) Get(i int) (uint64, bool) {
	if i >= e.nElements {
		return 0, false
	}
	cursor := 0
	highBit := 0
	for j := 0; j < e.highBits.Len(); j++ {
		v := e.highBits.Get(j)
		highBit += int((v + 1) & 0b1) // if v == 0
		cursor += int(v & 0b1)        // if v == 1
		if cursor == i+1 {
			break
		}
	}
	num := uint64(highBit)
	lowBitsPosition := e.lowBitsCount * i
	for k := 0; k < e.lowBitsCount; k++ {
		num = (num << 1) | uint64(e.lowBits.Get(lowBitsPosition+k))
	}
	return num, true
}
