package smalllist

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixedSized_GetSet(t *testing.T) {

	for nBit := 2; nBit < 16; nBit++ {
		max := int(math.Pow(2, float64(nBit)))
		f := &FixedSized{size:uint64(nBit), smalllist:make([]uint64, 100000)}

		for i := 0; i < max * 4; i++ {
			f.Set(uint64(i), i)
		}
		for i := 0; i < max; i++ {
			assert.Equal(t, uint64(i), f.Get(i))
		}
	}
}

func TestSelectKBits(t *testing.T) {
	for i := 0; i < 9; i++ {
		expected := uint64(math.Pow(2, float64(i))) - 1
		n := selectKBits(1023, uint64(i))
		assert.Equal(t, expected, n)
	}
}
