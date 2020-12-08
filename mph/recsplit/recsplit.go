package recsplit

import (
	"fmt"
	"math"
	"sync"

	"github.com/twmb/murmur3"
	"github.com/zs-projects/aeroview/mph/utils"
)

const (
	rNot       = 0
	maxRetries = 2 << 16
)

// Recsplit represents a minimal perfect hash function found using the recsplit algorithm.
type Recsplit struct {
	values [][]byte
	// One binary tree per bucket.
	keys    []RecsplitStree
	cumSums []int
}

func (r *Recsplit) SizeInBytes() int {
	u := 0
	for _, fbt := range r.keys {
		u += fbt.SizeInBytes()
	}
	return len(r.cumSums)*8 + u
}

// GetKey Returns the looked for value
func (r Recsplit) GetKey(s string) int {
	// 1. Determine bucket.
	bucket := hash(s, rNot) % len(r.keys)
	tree := r.keys[bucket]
	node := tree.Root()
	h := r.cumSums[bucket]
	out := -1
	for node != 0 || (node == 0 && out == -1) {
		R, nbKeys := tree.node(node)
		split := hash(s, uint64(R)) % nbKeys
		halfNbKeys := nbKeys / 2
		out = h + split
		if split < halfNbKeys {
			node = tree.LeftChild(node)
		} else {
			node = tree.RightChild(node)
			h += halfNbKeys
		}
	}
	return out
}

// Get Returns the looked for value
func (r Recsplit) Get(s string) []byte {
	return r.values[r.GetKey(s)]
}

// FromMap computes a minimal perfect hashing function over the map that is provided as an argument.
// nbWorkes controls the number of parallel go-routines to be launched to do the work.
func FromMap(data map[string][]byte, nbWorkers int) Recsplit {
	nBuckets := int(math.Max(float64(len(data))/100, 1))

	partitioner := func(data string) int {
		return hash(data, rNot) % nBuckets
	}
	// 1. Assign data to buckets.
	buckets := utils.AssignToBuckets(partitioner, data, nBuckets)
	splits := make(chan recsplitBucket, nBuckets)
	results := makeWorkerPool(nbWorkers, nBuckets, splits)
	for _, b := range buckets {
		// 4. send the bucket to the workers. ( it will be split or brute forced.)
		rb := recsplitBucket{
			keys:   b.Keys,
			bucket: b.OriginalIndex,
		}
		splits <- rb
	}
	close(splits)
	// 5. Collect the finalRecSplitbuckets.
	fBuckets := make(map[int]recsplitSubTree, nBuckets)
	// NOTE: for now this channel is closed nowhere. This will loop forever.
	for res := range results {
		fBuckets[res.bucket] = res
	}
	// 6. Reconstruct the mph.
	mph := mphFromRecsplitLeafs(fBuckets, data)
	return mph
}

func makeWorkerPool(nbWokers, nbBuckets int, splits <-chan recsplitBucket) <-chan recsplitSubTree {
	var wg sync.WaitGroup
	results := make(chan recsplitSubTree, nbBuckets)
	for i := 0; i < nbWokers; i++ {
		wg.Add(1)
		go recsplitWorker(&wg, splits, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()
	return results
}

func recsplitWorker(wg *sync.WaitGroup, buckets <-chan recsplitBucket, results chan<- recsplitSubTree) {
	for b := range buckets {
		tree := PartitionBucket(b)
		results <- recsplitSubTree{recsplitBucket: b, Node: tree}
	}
	wg.Done()
}

func PartitionBucket(b recsplitBucket) *Node {
	if nbKeys := len(b.keys); nbKeys <= 5 {
		r := b.bruteForceMPH()
		return &Node{
			Left:   nil,
			Right:  nil,
			R:      r,
			nbKeys: nbKeys,
		}
	} else {
		r, left, right := b.split()
		return &Node{
			Left:   PartitionBucket(left),
			Right:  PartitionBucket(right),
			R:      r,
			nbKeys: nbKeys,
		}
	}
}

func mphFromRecsplitLeafs(res map[int]recsplitSubTree, values map[string][]byte) Recsplit {
	nbBuckets := len(res)
	mph := Recsplit{
		keys:    make([]RecsplitStree, nbBuckets),
		values:  make([][]byte, len(values)),
		cumSums: make([]int, 1),
	}
	cumSum := 0
	for bucket, leafs := range res {
		mph.keys[bucket] = MakeFbTreeFromRecSplitSubTree(leafs)
	}
	for _, value := range mph.keys {
		_, nbKeys := value.node(value.Root())
		cumSum += nbKeys
		mph.cumSums = append(mph.cumSums, cumSum)
	}
	for key, value := range values {
		mph.values[mph.GetKey(key)] = value
	}
	return mph
}

type recsplitBucket struct {
	keys   []string
	bucket int
}
type recsplitSubTree struct {
	recsplitBucket
	*Node
}

func (b recsplitBucket) split() (int, recsplitBucket, recsplitBucket) {
	r, lKeys, rKeys := b.partitionKeys()
	left := recsplitBucket{
		keys:   lKeys,
		bucket: b.bucket,
	}
	right := recsplitBucket{
		keys:   rKeys,
		bucket: b.bucket,
	}
	return r, left, right
}

func (b recsplitBucket) partitionKeys() (int, []string, []string) {
	r := uint64(1)
	s := len(b.keys) / 2
	lKeys := make([]string, s+1)
	rKeys := make([]string, s+1)
	for {
		lPos := 0
		rPos := 0
		for _, key := range b.keys {
			if hash(key, r)%len(b.keys) < s {
				lKeys[lPos] = key
				lPos++
			} else {
				rKeys[rPos] = key
				rPos++
			}
			if rPos >= s+1 || lPos >= s+1 {
				break
			}
		}
		if lPos == s {
			return int(r), lKeys[:lPos], rKeys[:rPos]
		}
		r++
		if r >= maxRetries {
			panic(fmt.Sprintf("Can't partition keys. %v", b.keys))
		}
	}
}

func (b recsplitBucket) bruteForceMPH() int {
	r := uint64(1)
	collisions := make([]bool, len(b.keys))
	for {
		for i := range collisions {
			collisions[i] = false
		}
		if !checkForCollisions(r, b.keys, collisions) {
			return int(r)
		}
		r++
		if r >= maxRetries {
			panic(fmt.Sprintf("Can't brute force %v", b.keys))
		}
	}
}

func checkForCollisions(r uint64, keys []string, collisions []bool) bool {
	if len(collisions) < len(keys) {
		panic("in recplist.checkForCollision collisions buffer is too small.")
	}

	for _, key := range keys {
		h := hash(key, r) % len(keys)
		if !collisions[h] {
			collisions[h] = true
		} else {
			// found a collision, break!
			return true
		}
	}
	return false
}

var mask uint64 = (1<<64 - 1<<63) - 1

func hash(data string, r uint64) int {
	hash := murmur3.SeedStringSum64(r, data)
	// put the highest bit to 0, to make sure that we have a positive number when converting.ut the highest bit to 0, to make sure that we have a positive number when converting.
	return int(hash & mask)
}
