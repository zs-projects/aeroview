package rank

import (
	"fmt"
	"math/bits"
	"math/rand"
	"testing"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
)

func TestPopCount(t *testing.T) {
	size := 2500
	xs := make([]uint64, size)
	nbOnes := 0
	for i := 0; i < size; i++ {
		u := rand.Int63()
		xs[i] = uint64(u)
		nbOnes += bits.OnesCount64(uint64(u))
	}
	rpc := MakePopCount(xs)
	if len(rpc.Blocks)*8 < len(xs) {
		t.Errorf("%# v", pretty.Formatter(rpc))
	}
	fmt.Println(rpc.Overhead())
	for i := 0; i < size*64; i++ {
		expected := countBits(i, rpc.Data)
		actual := rpc.Rank(i)
		if expected != actual {
			fmt.Println("Shit.")
			rpc.Rank(i)
		}
		assert.EqualValues(t, expected, actual)
	}
	selectCnt := make(map[int]int)
	prev := 0
	for i := 0; i < nbOnes; i++ {
		actual := rpc.Rank(i)
		if actual != prev {
			selectCnt[actual] = i
			prev = actual
		}
	}
	// Some of them are not working as expected.
	for i, expected := range selectCnt {
		actual := rpc.Select(uint64(i))
		if actual != uint64(expected) {
			fmt.Println(actual, expected)
			rpc.Select(uint64(i))
		}
		assert.EqualValues(t, expected, actual, fmt.Sprintf("for position %v", i))
	}
}

func countBits(idx int, data []uint64) int {
	acc := 0
	j := idx / BLOCKSIZE
	for k := 0; k < j; k++ {
		acc += bits.OnesCount64(data[k])
	}
	shift := idx % BLOCKSIZE
	mask := uint64(1<<(shift+1) - 1)
	acc += bits.OnesCount64(data[j] & mask)
	return acc
}
