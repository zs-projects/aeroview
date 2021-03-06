package encoding

import (
	"math"

	"github.com/zs-projects/aeroview/datastructures/bits"
)

// EliasDeltaEncoding encode positive integers using elias code.
type EliasDeltaEncoding struct{}

// Encode64 encodes a slice of sorted ascending list of uint32 int
// TODO take into account the fact that the array is sorted to encode things .
func (EliasDeltaEncoding) Encode64(v uint64) ([]uint64, int) {
	b := bits.MakeQueue()
	// N such that such that X between 2^N and 2^(N+1)
	N := int64(math.Floor(math.Log2(float64(v))))
	N1 := 1 + N
	nbBitsOfN1 := int(math.Floor(math.Log2(float64(N1))))
	// Unary Coding .
	for i := 0; i < nbBitsOfN1; i++ {
		b.PushBack(0)
	}
	//
	for i := nbBitsOfN1; i >= 0; i-- {
		b.PushBack(uint64(N1 >> i))
	}
	for i := N1 - 2; i >= 0; i-- {
		b.PushBack(v >> i)
	}
	return b.Data(), b.Len()
}

// Decode64 encodes a slice of sorted ascending list of uint32 int
func (EliasDeltaEncoding) Decode64(b []uint64, size int) (uint64, error) {
	bq, err := bits.MakeBitQueueFromSlice(b, size)
	if err != nil {
		return 0, err
	}
	num := uint64(1)
	length := 1
	lengthOfLen := 0
	for {
		if !bq.Empty() && bq.Pop() == 0 {
			lengthOfLen++
		} else {
			break
		}
	}
	for i := 0; i < lengthOfLen; i++ {
		length <<= 1
		if !bq.Empty() && bq.Pop() == 0b1 {
			length |= 1
		}
	}
	for i := 0; i < length-1; i++ {
		num <<= 1
		if !bq.Empty() && bq.Pop() == 0b1 {
			num |= 1
		}
	}
	return num, nil
}
