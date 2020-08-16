package bits

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitVec_Get(t *testing.T) {
	b := NewVector(1)
	indexes := make(map[int32]struct{})

	for i := 0; i < 100000; i++ {
		n := rand.Int31() >> 16
		indexes[n] = struct{}{}
		b.Set(int(n), 1)
	}
	for x := range indexes {
		assert.Equal(t, 1, b.Get(int(x)))
	}
}

func TestBitVec_Get8BitRange(t *testing.T) {
	lookup := Make8BitLookup()

	// 1. from same block

	b1 := NewVector(1)
	b1.Set(7, 1)
	b1.Set(2, 1)
	n := b1.Get8BitRange(0, 7)
	assert.Equal(t, uint8(132), n)
	assert.Equal(t, uint8(2), lookup[n])

	// 2. from diff blocks

	b2 := NewVector(1)
	b2.Set(59, 1)
	b2.Set(64, 1)
	b2.Set(65, 1)

	/*
			from bit set representation
			indexes:	 65, 64, 63, 62, 61, 60, 59, ...
			values:  	 1,  1,   0,  0,  0,  0,  1,
		    block id     1,  1,   0,   0,  0,  0,  0,

			binary   1100001
			decimal  97
			number of ones: 3
	*/

	n = b2.Get8BitRange(59, 66)
	assert.Equal(t, uint8(97), n)
	assert.Equal(t, uint8(3), lookup[n])
}
