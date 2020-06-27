package recsplit

import (
	"testing"
)

func TestSplit(t *testing.T) {
	b := recsplitBucket{
		keys:    []string{"toto", "tata", "titi", "test", "tardif", "toff", "tiff", "tall", "health", "append", "count"},
		parents: []uint32{5, 7},
		isLeft:  []bool{true, false},
		bucket:  3,
	}
	l, r := b.split()
	if l.parents[2] != r.parents[2] {
		t.Errorf("Last element in parents should be the same.")
	}
	if !l.isLeft[2] {
		t.Errorf("isLeft value is wrong for the left bucket.")
	}
	if r.isLeft[2] {
		t.Errorf("isLeft value is wrong for the right bucket")
	}
}

func TestPartitionKeys(t *testing.T) {
	b := recsplitBucket{
		keys:    []string{"toto", "tata", "titi", "test", "tardif", "toff", "tiff", "tall", "health", "append", "count"},
		parents: nil,
	}
	s := len(b.keys) / 2
	r, lKeys, rKeys := b.partitionKeys()
	for _, key := range lKeys {
		if hash(key, uint32(r))%len(b.keys) >= s {
			t.Errorf("%v should be in the right partition", key)
		}
	}
	for _, key := range rKeys {
		if hash(key, uint32(r))%len(b.keys) < s {
			t.Errorf("%v should be in the left partition", key)
		}
	}
}
func TestBruteForceMPH(t *testing.T) {
	b := recsplitBucket{
		keys:    []string{"toto", "tata", "titi", "test", "tardif"},
		parents: nil,
	}
	r := b.bruteForceMPH()
	collisions := make([]bool, len(b.keys))
	if checkForCollisions(uint32(r), b.keys, collisions) {
		t.Errorf("Found a collision!")
	}
}
func TestCheckCollisions(t *testing.T) {
	keys := []string{"toto", "tata", "titi", "test", "tardif"}
	collisions := make([]bool, len(keys))
	if !checkForCollisions(uint32(0), keys, collisions) {
		t.Errorf("checkCollisions Failed")
	}
}
