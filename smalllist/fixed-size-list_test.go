package smalllist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixedSizeListGet(t *testing.T) {
	xs := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// length8
	smallList8 := From(xs, length8)
	for i := 0; i < len(xs); i++ {
		assert.Equal(t, i, smallList8.Get(i))
	}

	// length16
	smallList16 := From(xs, length16)
	for i := 0; i < len(xs); i++ {
		assert.Equal(t, i, smallList16.Get(i))
	}

	// length32
	smallList32 := From(xs, length32)
	for i := 0; i < len(xs); i++ {
		assert.Equal(t, i, smallList32.Get(i))
	}
}

func TestFixedSizedList_Set(t *testing.T) {
	xs := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	ys := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	smallList16 := From(xs, length16)
	for i := 0; i < len(ys); i++ {
		smallList16.Set(i, ys[i])
	}
	for i := 0; i < len(ys); i++ {
		assert.Equal(t, ys[i], smallList16.Get(i))
	}
}
