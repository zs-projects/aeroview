package recsplit

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"

	"zs-project.org/aeroview/datastructures"
	"zs-project.org/aeroview/mph/murmurhash"
	"zs-project.org/aeroview/mph/utils"
)

const (
	rNot       = 0
	maxRetries = 2 << 32
)

// Recsplit represents a minimal perfect hash function found using the recsplit algorithm.
type Recsplit struct {
	values [][]byte
	// One binary tree per bucket.
	keys    []datastructures.FBTree
	cumSums []int
}

// GetKey Returns the looked for value
func (r Recsplit) GetKey(s string) int {
	// 1. Determine bucket.
	bucket := hash(s, 0) % len(r.keys)
	tree := r.keys[bucket]
	node := tree.Root()
	h := r.cumSums[bucket]
	for {
		if tree.IsLeaf(*node) {
			return h + hash(s, uint64(node.Value.R))%node.Value.NbKeys
		}
		split := hash(s, uint64(node.Value.R)) % node.Value.NbKeys
		if split < node.Value.NbKeys/2 {
			node = tree.LeftChild(*node)
		} else {
			h = h + int(node.Value.NbKeys)/2
			node = tree.RightChild(*node)
		}

	}
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
	splits, results := makeWorkerPool(nbWorkers, nBuckets)
	overflows := make(map[int]struct{})
	for i, b := range buckets {
		// 2. Check for overflows.
		if len(b.Keys) > 2000 {
			// 3. Treat overflows with BDZ ( to be implemented_later.)
			overflows[i] = struct{}{}
		} else {
			// 4. send the bucket to the workers. ( it will be split or brute forced.)
			rb := recsplitBucket{
				keys:    b.Keys,
				parents: []uint32{},
				sizes:   []uint32{},
				bucket:  b.OriginalIndex,
			}
			splits <- rb
		}
	}
	// 5. Collect the finalRecSplitbuckets.
	fBuckets := make(map[int][]recsplitLeaf, 0)
	// NOTE: for now this channel is closed nowhere. This will loop forever.
	for res := range results {
		if _, ok := fBuckets[res.bucket]; !ok {
			fBuckets[res.bucket] = make([]recsplitLeaf, 0, 1)
		}
		fBuckets[res.bucket] = append(fBuckets[res.bucket], res)
	}
	// 6. Reconstruct the mph.
	mph := mphFromRecsplitLeafs(fBuckets, nBuckets, data)
	return mph
}

func makeWorkerPool(nbWokers int, initialWorkCount int) (chan<- recsplitBucket, <-chan recsplitLeaf) {
	var wg sync.WaitGroup
	work := int64(initialWorkCount)
	// Buffering in the channels here is important to avoid deadlocks
	// Given the pattern of concurrency we chose.
	// IT'S UGLY....
	splits := make(chan recsplitBucket, 5*nbWokers+1)
	results := make(chan recsplitLeaf, 5*nbWokers+1)
	for i := 0; i < nbWokers; i++ {
		wg.Add(1)
		go recsplitWorker(&wg, &work, splits, results)
	}
	go monitorWorkAndCloseChannels(&wg, &work, splits, results)
	return splits, results
}

func monitorWorkAndCloseChannels(wg *sync.WaitGroup, work *int64, splits chan recsplitBucket, results chan recsplitLeaf) {
	go func() {
		for {
			if atomic.LoadInt64(work) == 0 {
				close(splits)
				break
			}
		}
		wg.Wait()
		close(results)
	}()
}

func recsplitWorker(wg *sync.WaitGroup, workCount *int64, splits chan recsplitBucket, results chan<- recsplitLeaf) {
	for b := range splits {
		if len(b.keys) <= 5 {
			r := b.bruteForceMPH()
			b.parents = append(b.parents, uint32(r))
			b.sizes = append(b.sizes, uint32(len(b.keys)))
			atomic.AddInt64(workCount, -1)
			results <- recsplitLeaf{b}
		} else {
			left, right := b.split()
			atomic.AddInt64(workCount, +2)
			splits <- right
			splits <- left
			atomic.AddInt64(workCount, -1)
		}
	}
	wg.Done()
}

func mphFromRecsplitLeafs(res map[int][]recsplitLeaf, nbBuckets int, values map[string][]byte) Recsplit {
	mph := Recsplit{
		keys:    make([]datastructures.FBTree, nbBuckets),
		values:  make([][]byte, len(values)),
		cumSums: make([]int, 1, nbBuckets),
	}
	cumSum := 0
	for bucket, leafs := range res {
		tleafs := make([]datastructures.TreeLeaf, 0, len(leafs))
		for _, leaf := range leafs {
			tleafs = append(tleafs, leaf)
		}
		mph.keys[bucket] = datastructures.MakeFBTreeFromLeafs(tleafs)
		cumSum += mph.keys[bucket].Root().Value.NbKeys
		mph.cumSums = append(mph.cumSums, cumSum)
	}
	for key, value := range values {
		mph.values[mph.GetKey(key)] = value
	}
	return mph
}

type recsplitBucket struct {
	keys    []string
	parents []uint32
	sizes   []uint32
	isRight []bool
	bucket  int
}
type recsplitLeaf struct {
	recsplitBucket
}

func (r recsplitLeaf) Values() []datastructures.FBValue {
	vals := make([]datastructures.FBValue, 0, len(r.parents))
	for i, v := range r.parents {
		vals = append(vals, datastructures.FBValue{NbKeys: int(r.sizes[i]), R: int(v)})
	}
	return vals
}
func (r recsplitLeaf) Path() []bool {
	return r.isRight
}

type recsplitLeafs []recsplitLeaf

func (b recsplitBucket) split() (recsplitBucket, recsplitBucket) {
	r, lKeys, rKeys := b.partitionKeys()
	left := recsplitBucket{
		keys:    lKeys,
		parents: make([]uint32, len(b.parents)),
		sizes:   make([]uint32, len(b.sizes)),
		isRight: make([]bool, len(b.isRight)),
		bucket:  b.bucket,
	}
	right := recsplitBucket{
		keys:    rKeys,
		parents: make([]uint32, len(b.parents)),
		sizes:   make([]uint32, len(b.sizes)),
		isRight: make([]bool, len(b.isRight)),
		bucket:  b.bucket,
	}
	copy(right.parents, b.parents)
	copy(right.sizes, b.sizes)
	copy(right.isRight, b.isRight)
	copy(left.parents, b.parents)
	copy(left.sizes, b.sizes)
	copy(left.isRight, b.isRight)
	right.parents = append(right.parents, uint32(r))
	left.parents = append(left.parents, uint32(r))
	right.sizes = append(right.sizes, uint32(len(b.keys)))
	left.sizes = append(left.sizes, uint32(len(b.keys)))
	right.isRight = append(right.isRight, true)
	left.isRight = append(left.isRight, false)
	return left, right
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
		h := hash(key, uint64(r)) % len(keys)
		if !collisions[h] {
			collisions[h] = true
		} else {
			// found a collision, break!
			return true
		}
	}
	return false
}

func hash(data string, r uint64) int {
	v := murmurhash.Hash64(data) ^ r
	if int(v) < 0 {
		return -int(v)
	}
	return int(v)
}
