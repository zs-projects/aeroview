package rank

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestBitVec_Get(t *testing.T) {
	b := NewBitVec(1)
	indexes := make(map[int32]struct{})

	for i := 0; i < 100000; i++ {
		n := rand.Int31() >> 16
		indexes[n] = struct{}{}
		b.Set(int(n), 1)
	}
	for x := range indexes {
		if b.Get(int(x)) != 1 {

			fmt.Println("ooooo")
		}
	}
}

func TestBitVec_GetRange(t *testing.T) {
	b := NewBitVec(1)
	b.Set(2, 1)
	b.Set(4, 1)
	b.Set(66, 1)
	n := b.GetRange(50, 65)
	fmt.Println(count(int(n)))
	if count(int(n)) != 1 {
		fmt.Println("ooo")
	}
}

func TestBitVec_Get8BitRange(t *testing.T) {
	b := NewBitVec(1)
	b.Set(59, 1)
	b.Set(64, 1)
	b.Set(65, 1)

	n := b.Get8BitRange(59, 66)
	lookup := Make8BitLookup()
	fmt.Println(lookup[n])
}
