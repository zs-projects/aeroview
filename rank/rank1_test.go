package rank

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRank1(t *testing.T) {
	xs := make([]uint64, 1000)
	for i := 0; i < 1000; i++ {
		if i%2 == 0 {
			xs[i] = uint64(i)
		}
	}
	r1 := NewRank1(xs)
	for i := 0; i < 1000; i++ {
		assert.Equal(t, xs[i], r1.Get(i), "%d %d", xs[i], r1.Get(i))
	}
}

func TestRank1_Get(t *testing.T) {
	collectionSize := 2000
	nTries := 100
	rand.Seed(1)

	inc := uint64(0)
	for i := 0; i < nTries; i++ {

		// 1. build data
		xs := make([]uint64, collectionSize)
		for i := range xs {
			if rand.Intn(100) < 20 {
				xs[i] = inc
				inc++
			}
		}

		b := NewRank1(xs)
		for i := 0; i < len(xs); i++ {
			assert.Equal(t, xs[i], b.Get(i), "%dth example, expected %d but was %d", i, xs[i], b.Get(i))
		}
	}
}
