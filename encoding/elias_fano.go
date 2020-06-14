package encoding

import (
	"math"
)

// Encode64 encodes a slice of sorted ascending list of uint32 int
// TODO take into account the fact that the array is sorted to encode things .
func Encode64(values []uint64) ([]byte, int) {
	b := MakeBitQueue()
	for _, v := range values {
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
			b.PushBack(uint8(N1 >> i))
		}
		for i := N1 - 2; i >= 0; i-- {
			b.PushBack(uint8(v >> i))
		}
	}
	return b.Data(), b.Len()
}

// Decode64 encodes a slice of sorted ascending list of uint32 int
func Decode64(b []byte, size int) ([]uint64, error) {
	out := make([]uint64, 0)
	bq, err := MakeBitQueueFromSlice(b, size)
	if err != nil {
		return nil, err
	}
	for !bq.Empty() {
		num := uint64(1)
		len := 1
		lengthOfLen := 0
		for {
			if !bq.Empty() && bq.Pop() == 0 {
				lengthOfLen++
			} else {
				break
			}
		}
		for i := 0; i < lengthOfLen; i++ {
			len <<= 1
			if !bq.Empty() && bq.Pop() == 0b1 {
				len |= 1
			}
		}
		for i := 0; i < len-1; i++ {
			num <<= 1
			if !bq.Empty() && bq.Pop() == 0b1 {
				num |= 1
			}
		}
		out = append(out, num)
	}
	return out, nil
}
