package chd0

import (
	"errors"
	"sort"

	"github.com/twmb/murmur3"
	"github.com/zs-projects/aeroview/mph/utils"
)

const (
	rNot     = 0
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
		keyIndex = hash(key, uint32(chd.h[hIndex])) % len(chd.keys)
	}

	if chd.keys[keyIndex] != key {
		return nil, false
	}
	return chd.values[keyIndex], true
}

func FromMap(kv map[string][]byte) (*CHD, error) {

	keys := make([]string, len(kv))
	values := make([][]byte, len(kv))

	// there is a tradeoff here to make in the number of bucket we want.
	// more buckets means faster to build, but more memory to keep the structure in place.
	nBuckets := len(kv) / 2
	hashes := make([]int32, nBuckets)
	buckets := make(utils.Buckets, nBuckets)

	if nBuckets == 0 {
		nBuckets = 1
	}

	// 1. assign the different keys to buckets
	for key := range kv {
		bucketIndex := hash(key, rNot) % nBuckets
		if buckets[bucketIndex] == nil {
			buckets[bucketIndex] = utils.NewBucket(bucketIndex)
		}
		buckets[bucketIndex].Keys = append(buckets[bucketIndex].Keys, key)
	}

	// sort bucket by length, wanna start w/ larger bucket first.
	sort.Sort(sort.Reverse(buckets))

	// 2. choose the correct hashes for buckets with conflict.
	// - when finding non-conflicting hash, write it down.
	// - assign keys and values to the right place
	for _, bucket := range buckets {
		if len(bucket.Keys) == 1 {
			break
		}

		r := uint32(1)

		assignedIndexes := make(map[int]bool)
		bucketIndex := 0
		for bucketIndex < len(bucket.Keys) {
			key := bucket.Keys[bucketIndex]
			keyIndex := hash(key, r) % len(keys)

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
		for _, key := range bucket.Keys {
			keyIndex := hash(key, r) % len(keys)
			keys[keyIndex] = key
			values[keyIndex] = kv[key]
			hashes[bucket.OriginalIndex] = int32(r)
		}
	}

	// 3. get free slots so that we can do some manual assignment where bucket has len 1.
	freeSlots := make([]int, 0)
	for i, key := range keys {
		if len(key) == 0 {
			freeSlots = append(freeSlots, i)
		}
	}
	// manually assign keys and values to where bucket has len 1.
	// manually assigned keys have a negative hash.
	for _, bucket := range buckets {
		if bucket == nil || len(bucket.Keys) != 1 {
			continue
		}
		slotIndex := freeSlots[0]
		freeSlots = freeSlots[1:]
		hashes[bucket.OriginalIndex] = int32(-slotIndex - 1)
		keys[slotIndex] = bucket.Keys[0]
		values[slotIndex] = kv[bucket.Keys[0]]
	}
	return &CHD{
		keys:   keys,
		values: values,
		h:      hashes,
	}, nil
}

var mask uint64 = (1<<64 - 1<<63) - 1

func hash(data string, r uint32) int {
	hash := murmur3.SeedStringSum64(uint64(r), data)
	// put the highest bit to 0, to make sure that we have a positive number when converting.ut the highest bit to 0, to make sure that we have a positive number when converting.
	return int(hash & mask)
}
