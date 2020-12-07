package recsplit

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
)

func TestFromMap(t *testing.T) {
	nbKeys := 50000
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
	nbKeys := 5000
	splits := make(chan recsplitBucket)
	results := make(chan recsplitSubTree)
	keys := make([]string, 0)
	mp := make(map[string][]byte, nbKeys)
	for i := 0; i < nbKeys; i++ {
		data := fmt.Sprintf("test-%v", i)
		mp[data] = []byte(data)
		keys = append(keys, data)
	}
	b := recsplitBucket{
		keys:   keys,
		bucket: 0,
	}
	wg.Add(1)
	go recsplitWorker(&wg, splits, results)
	splits <- b
	close(splits)
	res := make(map[int]recsplitSubTree, 0)
	r := <-results
	res[r.bucket] = r
	wg.Wait()
	close(results)
	mph := mphFromRecsplitLeafs(res, mp)
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

func TestSplit(t *testing.T) {
	keys := []string{"toto", "tata", "titi", "test", "tardif", "tiff", "tall", "health", "append"}
	s := len(keys) / 2
	b := recsplitBucket{
		keys:   keys,
		bucket: 3,
	}
	R, l, r := b.split()
	if l.bucket != b.bucket {
		t.Errorf("left bucket was not preserved.")
	}
	if r.bucket != b.bucket {
		t.Errorf("right bucket was not preserved.")
	}
	for _, key := range l.keys {
		if hash(key, uint64(R))%len(b.keys) >= s {
			t.Errorf("%v should be in the right partition", key)
		}
	}
	for _, key := range r.keys {
		if hash(key, uint64(R))%len(b.keys) < s {
			t.Errorf("%v should be in the left partition", key)
		}
	}
}

func TestPartitionKeys(t *testing.T) {
	b := recsplitBucket{
		keys: []string{"toto", "tata", "titi", "test", "tardif", "toff", "tiff", "tall", "health", "append", "count"},
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

func haveCommonElements(setA, setB map[string]struct{}) bool {
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
		keys: []string{"toto", "tata", "titi", "test", "tardif"},
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
