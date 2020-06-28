package recsplit

import (
	"sync"
	"sync/atomic"

	"zs-project.org/aeroview/mph/farmhash"
	"zs-project.org/aeroview/mph/utils"
)

const (
	rNot       = 0
	maxRetries = 2 << 8
)

type MPH struct {
	// One binary tree per bucket.
	splits [][]node
}
type node struct {
	r      uint32
	isLeaf bool
}

// FromMap computes a minimal perfect hashing function over the map that is provided as an argument.
// nbWorkes controls the number of parallel go-routines to be launched to do the work.
func FromMap(data map[string][]byte, nbWorkers int) MPH {
	nBuckets := len(data) / 100
	partitioner := func(data string) int {
		return hash(data, rNot)
	}
	// 1. Assign data to buckets.
	buckets := utils.AssignToBuckets(partitioner, data, nBuckets)
	splits, results := makeWorkerPool(5, len(buckets))
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
				parents: []uint32{uint32(b.OriginalIndex)},
				bucket:  b.OriginalIndex,
			}
			splits <- rb
		}
	}
	// 5. Collect the finalRecSplitbuckets.
	fBuckets := make([]recsplitLeaf, 0)
	// NOTE: for now this channel is closed nowhere. This will loop forever.
	for res := range results {
		fBuckets = append(fBuckets, res)
	}
	// 6. Reconstruct the mph.
	return mphFromRecsplitLeafs(fBuckets, nBuckets)
}

func makeWorkerPool(nbWokers int, initialWorkCount int) (chan<- recsplitBucket, <-chan recsplitLeaf) {
	var wg sync.WaitGroup
	work := int64(initialWorkCount)
	splits := make(chan recsplitBucket)
	results := make(chan recsplitLeaf)
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
			atomic.AddInt64(workCount, -1)
			results <- recsplitLeaf{b}
		} else {
			left, right := b.split()
			splits <- right
			atomic.AddInt64(workCount, +1)
			splits <- left
			atomic.AddInt64(workCount, +1)
			atomic.AddInt64(workCount, -1)
		}
	}
	wg.Done()
}

func mphFromRecsplitLeafs(res []recsplitLeaf, nbBuckets int) MPH {
	mph := MPH{
		splits: make([][]node, nbBuckets),
	}
	for _, r := range res {
		// the first element is supposed to be a bucket.
		bucket := r.bucket
		if mph.splits[bucket] != nil {
			mph.splits[bucket] = make([]node, 0, 2^(len(r.parents)-1))
		}
		// If the lenght of the array is not big enough to handle the whole tree, we allocate a new one.
		if len(mph.splits[bucket]) <= 2^len(r.parents) {
			dest := make([]node, 2^(len(r.parents)-1))
			copy(dest, mph.splits[bucket])
			mph.splits[bucket] = dest
		}
		for k, n := range r.parents[1:] {
			pos := 2 ^ k - 1
			if !r.isLeft[k+1] {
				pos++
			}

			mph.splits[bucket][pos] = node{
				r: n,
			}
			if k == len(r.parents[1:])-1 {
				mph.splits[bucket][pos].isLeaf = true
			}
		}
	}
	return mph
}

type recsplitBucket struct {
	keys    []string
	parents []uint32
	isLeft  []bool
	bucket  int
}
type recsplitLeaf struct {
	recsplitBucket
}
type recsplitLeafs []recsplitLeaf

func (b recsplitBucket) split() (recsplitBucket, recsplitBucket) {
	r, lKeys, rKeys := b.partitionKeys()
	left := recsplitBucket{
		keys:    lKeys,
		parents: make([]uint32, len(b.parents)),
		isLeft:  make([]bool, len(b.isLeft)),
		bucket:  b.bucket,
	}
	right := recsplitBucket{
		keys:    rKeys,
		parents: make([]uint32, len(b.parents)),
		isLeft:  make([]bool, len(b.isLeft)),
		bucket:  b.bucket,
	}
	copy(right.parents, b.parents)
	copy(right.isLeft, b.isLeft)
	copy(left.parents, b.parents)
	copy(left.isLeft, b.isLeft)
	right.parents = append(right.parents, uint32(r))
	left.parents = append(left.parents, uint32(r))
	right.isLeft = append(right.isLeft, false)
	left.isLeft = append(left.isLeft, true)
	return left, right
}

func (b recsplitBucket) partitionKeys() (int, []string, []string) {
	r := uint32(1)
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
		}
		if lPos == s {
			return int(r), lKeys[:lPos], rKeys[:rPos]
		}
		r++
		if r >= maxRetries {
			panic("Can't partition keys.")
		}
	}
}

func (b recsplitBucket) bruteForceMPH() int {
	r := uint32(1)
	collisions := make([]bool, len(b.keys))
	for {
		for i := range collisions {
			collisions[i] = false
		}
		if !checkForCollisions(r, b.keys, collisions) {
			return int(r)
		}
		r++
	}
}
func checkForCollisions(r uint32, keys []string, collisions []bool) bool {
	for _, key := range keys {
		h := hash(key, uint32(r)) % len(keys)
		if !collisions[h] {
			collisions[h] = true
		} else {
			// found a collision, break!
			return true
		}
	}
	return false
}

func hash(data string, r uint32) int {
	return int(farmhash.Hash32(data) ^ r)
}
