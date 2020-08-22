package recsplit

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
)

func TestFromMap(t *testing.T) {
	nbKeys := 50000000
	mp := make(map[string][]byte, nbKeys)
	for i := 0; i < nbKeys; i++ {
		data := fmt.Sprintf("test-%v", i)
		mp[data] = []byte(data)
	}
	mph := FromMap(mp, 20)
	collisions := make([]bool, len(mp))
	for key, value := range mp {
		i := mph.GetKey(key)
		if i >= len(collisions) {
			t.Errorf("Index out of range, did not learn an mph %v", i)
		} else {
			if i < len(collisions) && collisions[i] {
				t.Errorf("Found collision for value %v that hashed at %v", key, i)
			} else {
				collisions[i] = true
			}
		}
		if !bytes.Equal(mph.Get(key), value) {
			t.Errorf("Values assignement is corrupt, expecting %v, got %v", string(value), string(mph.Get(key)))
		}
	}
}

func TestMPHFromRecsplitLeafs(t *testing.T) {
	var wg sync.WaitGroup
	workCount := int64(1)
	nbKeys := 5000
	splits := make(chan recsplitBucket, nbKeys)
	results := make(chan recsplitLeaf, nbKeys)
	keys := make([]string, 0)
	mp := make(map[string][]byte, nbKeys)
	for i := 0; i < nbKeys; i++ {
		data := fmt.Sprintf("test-%v", i)
		mp[data] = []byte(data)
		keys = append(keys, data)
	}
	b := recsplitBucket{
		keys:    keys,
		parents: []uint32{},
		isRight: []bool{},
		bucket:  0,
	}
	splits <- b
	wg.Add(1)
	go monitorWorkAndCloseChannels(&wg, &workCount, splits, results)
	go recsplitWorker(&wg, &workCount, splits, results)
	res := make(map[int][]recsplitLeaf, 0)
	for r := range results {
		if _, ok := res[r.bucket]; !ok {
			res[r.bucket] = make([]recsplitLeaf, 0, 1)
		}
		res[r.bucket] = append(res[r.bucket], r)
	}
	mph := mphFromRecsplitLeafs(res, 1, mp)
	collisions := make([]bool, len(keys))
	for key, value := range mp {
		i := mph.GetKey(key)
		if i >= len(collisions) {
			t.Errorf("Index out of range, did not learn an mph %v", i)
		}
		if i < len(collisions) && collisions[i] {
			t.Errorf("Found collision for value %v that hashed at %v", key, i)
		} else {
			collisions[i] = true
		}
		if !bytes.Equal(mph.Get(key), value) {
			t.Errorf("Values assignement is corrupt, expecting %v, got %v", string(value), string(mph.Get(key)))
		}
	}
}

func TestRecsplitWorker(t *testing.T) {
	var wg sync.WaitGroup
	workCount := int64(1)
	splits := make(chan recsplitBucket, 3)
	results := make(chan recsplitLeaf, 3)
	b := recsplitBucket{
		keys:    []string{"toto", "tata", "titi", "test", "tardif", "toff", "tiff", "tall", "health", "append", "count"},
		parents: []uint32{5, 7},
		sizes:   []uint32{44, 22},
		isRight: []bool{true, false},
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
		if len(r.parents) != len(r.sizes) {
			t.Errorf("parent and size slices ar of different length for %v", r)
		}
		collisions := make([]bool, len(r.keys))
		if checkForCollisions(uint64(r.parents[len(r.parents)-1]), r.keys, collisions) {
			t.Errorf("checkCollisions Failed for %v because %v", r, collisions)
		}
	}
	if count != 3 {
		t.Errorf("Excepcted two results.")
	}
}

func TestSplit(t *testing.T) {
	keys := []string{"toto", "tata", "titi", "test", "tardif", "tiff", "tall", "health", "append"}
	s := len(keys) / 2
	b := recsplitBucket{
		keys:    keys,
		parents: []uint32{5, 7},
		sizes:   []uint32{44, 22},
		isRight: []bool{true, false},
		bucket:  3,
	}
	l, r := b.split()
	if l.parents[2] != r.parents[2] {
		t.Errorf("Last element in parents should be the same.")
	}
	if l.sizes[2] != r.sizes[2] {
		t.Errorf("Last element in parents should be the same.")
	}
	if l.isRight[2] {
		t.Errorf("isRight value is wrong for the left bucket.")
	}
	if l.bucket != b.bucket {
		t.Errorf("left bucket was not preserved.")
	}
	if !r.isRight[2] {
		t.Errorf("isRight value is wrong for the right bucket")
	}
	if r.bucket != b.bucket {
		t.Errorf("right bucket was not preserved.")
	}
	for _, key := range l.keys {
		if hash(key, uint64(l.parents[2]))%len(b.keys) >= s {
			t.Errorf("%v should be in the right partition", key)
		}
	}
	for _, key := range r.keys {
		if hash(key, uint64(r.parents[2]))%len(b.keys) < s {
			t.Errorf("%v should be in the left partition", key)
		}
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
		if hash(key, uint64(r))%len(b.keys) >= s {
			t.Errorf("%v should be in the right partition", key)
		}
	}
	for _, key := range rKeys {
		if hash(key, uint64(r))%len(b.keys) < s {
			t.Errorf("%v should be in the left partition", key)
		}
	}
	if len(rKeys)+len(lKeys) != len(b.keys) {
		t.Errorf("parttion did not produce the same elements as output.")
	}

	lSet := toStringSet(lKeys)
	rSet := toStringSet(rKeys)
	if haveCommonElements(lSet, rSet) {
		t.Errorf("Left and Right partitions should be disjoint: Got Right %v, and Left %v", rKeys, lKeys)
	}
}

func toStringSet(st []string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, val := range st {
		set[val] = struct{}{}
	}
	return set
}

func haveCommonElements(setA map[string]struct{}, setB map[string]struct{}) bool {
	for key := range setA {
		if _, ok := setB[key]; ok {
			return true
		}
	}
	for key := range setB {
		if _, ok := setA[key]; ok {
			return true
		}
	}
	return false
}

func TestBruteForceMPH(t *testing.T) {
	b := recsplitBucket{
		keys:    []string{"toto", "tata", "titi", "test", "tardif"},
		parents: nil,
	}
	r := b.bruteForceMPH()
	collisions := make([]bool, len(b.keys))
	if checkForCollisions(uint64(r), b.keys, collisions) {
		t.Errorf("Found a collision!")
	}
}
func TestCheckCollisions(t *testing.T) {
	keys := []string{"toto", "tata", "titi", "test", "tardif"}
	collisions := make([]bool, len(keys))
	if !checkForCollisions(uint64(0), keys, collisions) {
		t.Errorf("checkCollisions Failed")
	}
}
