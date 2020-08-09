package rank

import (
	"fmt"
	"math/bits"
	"testing"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
)

func TestPopCount(t *testing.T) {
	size := 1000
	xs := make([]uint64, size)
	nbOnes := 0
	for i := 0; i < size; i++ {
		if i%2 == 0 {
			xs[i] = uint64(i)
			nbOnes += bits.OnesCount64(uint64(i))
		}
	}
	rpc := MakePopCount(xs)
	if len(rpc.Blocks)*8 < len(xs) {
		t.Errorf("%# v", pretty.Formatter(rpc))
	}
	fmt.Println(rpc.Overhead())
	for i := 0; i < size*64; i++ {
		expected := countBits(i, rpc.data)
		actual := rpc.Rank(i)
		assert.EqualValues(t, expected, actual)
	}
	selectCnt := make(map[int]int)
	prev := 0
	for i := 0; i < size*64; i++ {
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
		assert.EqualValues(t, expected, actual)
	}
}

func countBits(idx int, data []uint64) int {
	acc := 0
	j := idx / BLOCKSIZE
	for k := 0; k < j; k++ {
		acc += bits.OnesCount64(data[k])
	}
	shift := BLOCKSIZE - idx%BLOCKSIZE - 1
	acc += bits.OnesCount64(data[j] >> shift)
	return int(acc)
}
