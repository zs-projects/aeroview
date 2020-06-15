package encoding

import (
	"math"
)

// EliasFanoEncoding encodes a list of ascending integers using the Elias Fano Code.
type EliasFanoEncoding struct{}

// Encode64 encodes a list of uint64 using EliasFano Code.
func (EliasFanoEncoding) Encode64(values []uint64) ([]byte, int) {
	lowBitsQ := MakeBitQueue()
	highBitsQ := MakeBitQueue()
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
	for !lowBitsQ.Empty() {
		highBitsQ.PushBack(lowBitsQ.Pop())
	}
	return highBitsQ.Data(), highBitsQ.Len()
}
