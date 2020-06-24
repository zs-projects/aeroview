package recsplit

import (
	"sync"

	"zs-project.org/aeroview/mph/farmhash"
	"zs-project.org/aeroview/mph/utils"
)

const (
	rNot       = 0
	maxRetries = 2 << 8
)

func fromMap(data map[string][]byte, nbWorkers int) {
	_, splits, results := makeWorkerPool(5)
	nBuckets := len(data) / 100
	partitioner := func(data string) int {
		return hash(data, rNot)
	}
	// 1. Assign data to buckets.
	buckets := utils.AssignToBuckets(partitioner, data, nBuckets)
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
			}
			splits <- rb
		}
	}
	// 5. Collect the finalRecSplitbuckets.
	fBuckets := make([]finalRecsplitBucket, 0)
	// NOTE: for now this channel is closed nowhere. This will loop forever.
	for res := range results {
		fBuckets = append(fBuckets, res)
	}
	// 6. Reconstruct the mph.
}

func makeWorkerPool(nbWokers int) (*sync.WaitGroup, chan<- recsplitBucket, <-chan finalRecsplitBucket) {
	var wg sync.WaitGroup
	// TO FIGURE OUT : Where to close this channel.
	splits := make(chan recsplitBucket)
	// TO FIGURE OUT : Where to close this channel.
	results := make(chan finalRecsplitBucket)
	for i := 0; i < nbWokers; i++ {
		wg.Add(1)
		go recsplitWorker(&wg, splits, results)
	}
	return &wg, splits, results
}

func recsplitWorker(wg *sync.WaitGroup, splits chan recsplitBucket, results chan<- finalRecsplitBucket) {
	// NOTE: for now this channel is closed nowhere. This will loop forever.
	for b := range splits {
		if len(b.keys) <= 8 {
			r := b.bruteForceMPH()
			b.parents = append(b.parents, uint32(r))
			results <- finalRecsplitBucket{b}
		} else {
			r, lKeys, rKeys := b.partitionKeys()
			left := recsplitBucket{
				keys:    lKeys,
				parents: make([]uint32, len(b.parents)),
				isLeft:  make([]bool, len(b.isLeft)),
			}
			right := recsplitBucket{
				keys:    rKeys,
				parents: make([]uint32, len(b.parents)),
				isLeft:  make([]bool, len(b.isLeft)),
			}
			copy(right.parents, b.parents)
			copy(right.isLeft, b.isLeft)
			copy(left.parents, b.parents)
			copy(left.isLeft, b.isLeft)
			right.parents = append(right.parents, uint32(r))
			left.parents = append(left.parents, uint32(r))
			right.isLeft = append(right.isLeft, false)
			left.isLeft = append(left.isLeft, true)
			splits <- right
			splits <- left
		}
	}
	wg.Done()
}

type recsplitBucket struct {
	keys    []string
	parents []uint32
	isLeft  []bool
}
type finalRecsplitBucket struct {
	recsplitBucket
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
			if hash(key, r) < s {
				lKeys[lPos] = key
				lPos++
			} else {
				rKeys[lPos] = key
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

func newRecSplitBucket() recsplitBucket {
	return recsplitBucket{
		keys:    make([]string, 0),
		parents: make([]uint32, 0)}
}

func (b recsplitBucket) bruteForceMPH() int {
	r := uint32(1)
	for {
		collisions := make([]bool, len(b.keys))
		foundCollision := false
		for _, key := range b.keys {
			h := hash(key, r)
			if collisions[h] != false {
				collisions[h] = true
			} else {
				foundCollision = true
				// found a collision, break!
				break
			}
		}
		if !foundCollision {
			return int(r)
		}
	}
}

func hash(data string, r uint32) int {
	return int(farmhash.Hash32(data) ^ r)
}
