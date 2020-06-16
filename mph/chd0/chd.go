package chd0

import (
	"errors"
	"sort"
)

const (
	rNot = 0
	maxTries = 1000
)

type CHD struct {
	keys   []string
	values [][]byte
	h      []int32
}

func (chd *CHD) Get(key string) ([]byte, bool) {
	// 1. get the bucket
	hIndex := int(hash(key, rNot)) % len(chd.h)

	// 2. apply the correct hash on the key
	var keyIndex int
	if chd.h[hIndex] < 0 {
		keyIndex = int(-chd.h[hIndex]) - 1
	} else {
		keyIndex = int(hash(key, uint32(chd.h[hIndex]))) % len(chd.keys)
	}

	if chd.keys[keyIndex] != key {
		return nil, false
	}
	return chd.values[keyIndex], true
}

func From(kv map[string][]byte) (*CHD, error) {

	keys := make([]string, len(kv))
	values := make([][]byte, len(kv))

	// there is a tradeoff here to make in the number of bucket we want.
	// more buckets means faster to build, but more memory to keep the structure in place.
	nBuckets := len(kv) / 2
	hashes := make([]int32, nBuckets)
	buckets := make(buckets, nBuckets)

	if nBuckets == 0 {
		nBuckets = 1
	}

	// 1. assign the different keys to buckets
	for key := range kv {
		bucketIndex := int(hash(key, rNot)) % nBuckets
		if buckets[bucketIndex] == nil {
			buckets[bucketIndex] = newBucket(bucketIndex)
		}
		buckets[bucketIndex].keys = append(buckets[bucketIndex].keys, key)
	}

	// sort bucket by length, wanna start w/ larger bucket first.
	sort.Sort(sort.Reverse(buckets))

	// 2. choose the correct hashes for buckets with conflict.
	// - when finding non-conflicting hash, write it down.
	// - assign keys and values to the right place
	for _, bucket := range buckets {
		if len(bucket.keys) == 1 {
			break
		}

		r := uint32(1)

		assignedIndexes := make(map[int]bool)
		bucketIndex := 0
		for bucketIndex < len(bucket.keys) {
			key := bucket.keys[bucketIndex]
			keyIndex := int(hash(key, r)) % len(keys)

			// if conflict remains, re-init
			if len(keys[keyIndex]) != 0 || assignedIndexes[keyIndex] {
				bucketIndex = 0
				assignedIndexes = map[int]bool{}
				r++
				if r > maxTries {
					return nil, errors.New("fail to generate a CHD")
				}
				continue
			}
			bucketIndex++
			assignedIndexes[keyIndex] = true
		}
		// stable config was found, write down keys, values and the stable r.
		for _, key := range bucket.keys {
			keyIndex := int(hash(key, r)) % len(keys)
			keys[keyIndex] = key
			values[keyIndex] = kv[key]
			hashes[bucket.originalIndex] = int32(r)
		}
	}

	// 3. get free slots so that we can manually assign them to len1 bucket.
	freeSlots := make([]int, 0)
	for i, key := range keys {
		if len(key) == 0 {
			freeSlots = append(freeSlots, i)
		}
	}
	// manually assign keys and values to bucket w/ 1 key only.
	// manually assigned keys have a negative hash.
	for _, bucket := range buckets {
		if bucket == nil || len(bucket.keys) != 1 {
			continue
		}
		slotIndex := freeSlots[0]
		freeSlots = freeSlots[1:]
		hashes[bucket.originalIndex] = int32(-slotIndex - 1)
		keys[slotIndex] = bucket.keys[0]
		values[slotIndex] = kv[bucket.keys[0]]
	}
	return &CHD{
		keys:   keys,
		values: values,
		h:      hashes,
	}, nil
}

func hash(data string, r uint32) uint32 {
	var hash uint32 = 0x01000193
	for _, c := range data {
		hash ^= uint32(c)
		hash *= 0x01000193
	}
	return hash ^ r
}

type bucket struct {
	originalIndex int
	keys []string
}

func newBucket(index int) *bucket{
	return &bucket{
		originalIndex: index,
		keys:          nil,
	}
}

type buckets []*bucket

// Implements the Sort interface.
func (b buckets) Len() int {
	return len(b)
}

func (b buckets) Less(i, j int) bool {
	if b[i] != nil && b[j] != nil {
		return len(b[i].keys) < len(b[j].keys)
	} else if b[i] != nil {
		return false
	}
	return true
}

func (b buckets) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
