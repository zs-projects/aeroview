package chd0

import (
	"errors"
	"sort"
)

type CHD struct {
	keys   []string
	values [][]byte
	h      []int32
}

func (chd *CHD) Get(key string) ([]byte, bool) {
	// 1. get the bucket
	hIndex := int(hash(key, 0)) % len(chd.h)

	// 2. apply the correct hash.
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

func from(kv map[string][]byte) (*CHD, error) {

	nBuckets := len(kv)
	keys := make([]string, len(kv))
	values := make([][]byte, len(kv))
	hashes := make([]int32, len(kv))

	if nBuckets == 0 {
		nBuckets = 1
	}

	// 1. assign to buckets
	buckets := make([]*bucket, nBuckets)
	for key := range kv {
		bucketIndex := int(hash(key, 0)) % nBuckets
		if buckets[bucketIndex] == nil {
			buckets[bucketIndex] = &bucket{
				originalIndex: bucketIndex,
				keys:          []string{},
			}
		}
		buckets[bucketIndex].keys = append(buckets[bucketIndex].keys, key)
	}

	// sort bucket by length, wanna start w/ larger bucket first.
	sort.Slice(buckets, func(i, j int) bool {
		if buckets[i] != nil && buckets[j] != nil {
			return len(buckets[i].keys) > len(buckets[j].keys)
		} else if buckets[i] != nil {
			return true
		} else {
			return false
		}
	})

	// 2. assign keys to the right place
	for _, bucket := range buckets {
		if len(bucket.keys) == 1 {
			break
		}

		r := uint32(1)

		var leftOverKeys []string
		for len(bucket.keys) != 0 {
			key := bucket.keys[0]
			bucket.keys = bucket.keys[1:]

			keyIndex := int(hash(key, r)) % len(keys)
			if len(keys[keyIndex]) != 0 {
				leftOverKeys = append(leftOverKeys, key)
			} else {
				hashes[bucket.originalIndex] = int32(r)
				keys[keyIndex] = key
				values[keyIndex] = kv[key]
			}

			if len(bucket.keys) == 0 {
				for _, leftOverKey := range leftOverKeys {
					bucket.keys = append(bucket.keys, leftOverKey)
				}
				leftOverKeys = []string{}
				r++
			}
			if r > 1000 {
				return nil, errors.New("fail to generate a CHD")
			}
		}
	}

	// 3. get free slots so that we can manually assign them to len1 bucket.
	freeSlots := make([]int, 0)
	for i, key := range keys {
		if len(key) == 0 {
			freeSlots = append(freeSlots, i)
		}
	}

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
