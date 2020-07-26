package rank

import (
	"fmt"
	"math/bits"
	"testing"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
)

func TestRankPpc(t *testing.T) {
	size := 1000
	xs := make([]uint64, size)
	for i := 0; i < size; i++ {
		if i%2 == 0 {
			xs[i] = uint64(i)
		}
	}
	rpc := MakeRankPopCount(xs)
	if len(rpc.Blocks)*8 < len(xs) {
		t.Errorf("%# v", pretty.Formatter(rpc))
	}
	fmt.Println(rpc.Overhead())
	for i := 0; i < size*64; i++ {
		expected := countBits(i, rpc.data)
		actual := rpc.Rank(i)
		if actual != expected {
			fmt.Println("Coucou")
		}
		assert.EqualValues(t, expected, actual)
	}
}

func countBits(idx int, data []uint64) int {
	acc := 0
	j := idx / BLOCKSIZE
	shift := BLOCKSIZE - idx%BLOCKSIZE
	for k := 0; k < j; k++ {
		acc += bits.OnesCount64(data[k])
	}
	acc += bits.OnesCount64(data[j] >> shift)
	return int(acc)
}
