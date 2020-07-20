package rank

import (
	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, 1, b.Get(int(x)))
	}
}

func TestBitVec_Get8BitRange(t *testing.T) {
	lookup := Make8BitLookup()

	b := NewBitVec(1)
	b.Set(59, 1)
	b.Set(64, 1)
	b.Set(65, 1)

	/*
		from bit set representation
		indexes:	 65, 64, 63, 62, 61, 60, 59, ...
		values:  	 1,  1,   0,  0,  0,  0,  1,
	    block id     1,  1,   0,   0,  0,  0,  0,

		binary   1100001
		decimal  97
		number of ones: 3
	*/

	n := b.Get8BitRange(59, 66)
	assert.Equal(t, 97, n)
	assert.Equal(t, uint8(3), lookup[n])
}
