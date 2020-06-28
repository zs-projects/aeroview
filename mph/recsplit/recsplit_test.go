package recsplit

import (
	"sync"
	"testing"
)

func TestMPHFromRecsplitLeafs(t *testing.T) {
	var wg sync.WaitGroup
	workCount := int64(1)
	splits := make(chan recsplitBucket, 3)
	results := make(chan recsplitLeaf, 3)
	b := recsplitBucket{
		keys:    []string{"toto", "tata", "titi", "test", "tardif", "toff", "tiff", "tall", "health", "append", "count"},
		parents: []uint32{},
		isLeft:  []bool{},
		bucket:  0,
	}
	splits <- b
	wg.Add(1)
	go monitorWorkAndCloseChannels(&wg, &workCount, splits, results)
	go recsplitWorker(&wg, &workCount, splits, results)
	res := make([]recsplitLeaf, 0)
	for r := range results {
		res = append(res, r)
	}
	mphFromRecsplitLeafs(res, 1)
}

func TestRecsplitWorker(t *testing.T) {
	var wg sync.WaitGroup
	workCount := int64(1)
	splits := make(chan recsplitBucket, 3)
	results := make(chan recsplitLeaf, 3)
	b := recsplitBucket{
		keys:    []string{"toto", "tata", "titi", "test", "tardif", "toff", "tiff", "tall", "health", "append", "count"},
		parents: []uint32{5, 7},
		isLeft:  []bool{true, false},
		bucket:  3,
	}
	splits <- b
	wg.Add(1)
	go monitorWorkAndCloseChannels(&wg, &workCount, splits, results)
	go recsplitWorker(&wg, &workCount, splits, results)
	count := 0
	for r := range results {
		count++
		if len(r.keys) > 5 {
			t.Errorf("Non leaf node %v was sent with the results.", r)
		}
		collisions := make([]bool, len(r.keys))
		if checkForCollisions(uint32(r.parents[len(r.parents)-1]), r.keys, collisions) {
			t.Errorf("checkCollisions Failed for %v because %v", r, collisions)
		}
	}
	if count != 3 {
		t.Errorf("Excepcted two results.")
	}
}

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
	if l.bucket != b.bucket {
		t.Errorf("left bucket was not preserved.")
	}
	if r.isLeft[2] {
		t.Errorf("isLeft value is wrong for the right bucket")
	}
	if r.bucket != b.bucket {
		t.Errorf("right bucket was not preserved.")
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
