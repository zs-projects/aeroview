package rank

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBabyRank_Get(t *testing.T) {
	collectionSize := 2000
	nTries := 100
	rand.Seed(1)

	for i := 0; i < nTries; i++ {

		// 1. build data
		xs := make([]uint64, collectionSize)
		for i := range xs {
			if rand.Intn(100) < 20 {
				xs[i] = rand.Uint64()
			}
		}

		b := NewBabyRank(xs)
		for i := 0; i < len(xs); i++ {
			assert.Equal(t, xs[i], b.Get(i))
		}
	}
}
