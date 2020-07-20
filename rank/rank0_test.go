package rank

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestB(t *testing.T) {
	xs := make([]uint64, 1000)
	for i := 0; i < 1000; i++ {
		if i % 2 == 0 {
			xs[i] = uint64(i)
		}
	}
	r0 := NewRank0(xs)
	for i := 0; i < 1000; i++ {
		assert.Equal(t, xs[i], r0.Get(i))
	}
}

