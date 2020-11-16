package memtable

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zs-projects/aeroview/mph/farmhash"
)

func TestFromMap(t *testing.T) {
	size := 10
	mp := getSampleMap(size)
	memtable := FromMap(mp)
	assert.Equal(t, len(memtable.keys), size)
	assert.Equal(t, len(memtable.offsets), size+1)
	assert.EqualValues(t, memtable.offsets[0], 0)
	assert.True(t, sort.SliceIsSorted(memtable.keys, func(i, j int) bool { return memtable.keys[i] < memtable.keys[j] }))
	for idx, key := range memtable.originalKeys {
		assert.EqualValues(t, farmhash.Hash64(key), memtable.keys[idx])
	}
}

func TestMemtableGet(t *testing.T) {
	size := 5000
	mp := getSampleMap(size)
	memtable := FromMap(mp)
	assert.Equal(t, len(memtable.keys), size)
	assert.Equal(t, len(memtable.offsets), size+1)
	assert.EqualValues(t, memtable.offsets[0], 0)
	for key, value := range mp {
		if gotValue, ok := memtable.Get(key); ok {
			assert.EqualValues(t, value, gotValue)
		} else {
			t.Fatalf("failed to get key %v", value)
		}
	}
}
