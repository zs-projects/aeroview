package encoding

import (
	"fmt"
	"math"

	"github.com/zs-projects/aeroview/datastructures/bits"
	"github.com/zs-projects/aeroview/rank"
)

// EliasFanoVector encodes a list of ascending integers using the Elias Fano Code.
type EliasFanoVector struct {
	highBits      bits.Queue
	lowBits       bits.Queue
	rank          rank.PopCount
	nElements     int // the number of elements in the data structure.
	lowBitsCount  int // The number of bits used to encode the low bits.
	highBitsCount int // The number of bits used to encode the low bits.
}

// MakeEliasFanoVector encodes a list of uint64 using EliasFano Code.
func MakeEliasFanoVector(values []uint64) EliasFanoVector {
	lowBitsQ := bits.MakeQueue()
	highBitsQ := bits.MakeQueue()
	lowerBitCount := int64(math.Log2(float64(values[len(values)-1]) / float64(len(values))))
	// Important to ensure mechanical sympathie and that the lower bits of any numbers
	// will never overleap between two bytes
	if lowerBitCount%2 != 0 {
		lowerBitCount = lowerBitCount + 1
	}
	if lowerBitCount < 2 {
		// Trick to handle long arrays of very small 16 bits numbers.
		lowerBitCount = 2
	}
	lowBitsMask := uint64(1)<<lowerBitCount - 1
	prev := uint64(0)
	for i := 0; i < len(values); i++ {
		v := values[i]
		highBits := v >> lowerBitCount
		lowBits := v & lowBitsMask
		highDelta := highBits - prev
		for k := int64(0); k < lowerBitCount; k++ {
			lowBitsQ.PushBack(uint64(lowBits >> k))
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
	r := rank.MakePopCount(highBitsQ.Vector())
	return EliasFanoVector{
		highBits:     highBitsQ,
		lowBits:      lowBitsQ,
		rank:         r,
		nElements:    len(values),
		lowBitsCount: int(lowerBitCount),
	}
}

// Len returns the number of bits set in the vector.
func (e *EliasFanoVector) Len() int {
	return e.highBits.Len() + e.lowBits.Len()
}

// Data returns the raw elias-fano encoded array.
func (e *EliasFanoVector) Data() ([]uint64, int) {
	result, _ := bits.MakeBitQueueFromSlice(e.highBits.Data(), e.highBits.Len())
	result.Append(e.lowBits.Data(), e.lowBits.Len())
	return result.Data(), result.Len()
}

// Get Returns the element at index i
func (e *EliasFanoVector) Get(i int) uint64 {
	if i >= e.nElements {
		panic(fmt.Sprintf("Trying to access element with index %v on EliasFanoVector on length %v", i, e.nElements))
	}
	highBit := e.rank.Select(uint64(i+1)) - uint64(i)
	num := uint64(highBit)
	lowBitsPosition := e.lowBitsCount * i
	lowbits := uint64(e.lowBits.GetN(lowBitsPosition, e.lowBitsCount-1))
	num = (num << e.lowBitsCount) | lowbits
	return num
}
