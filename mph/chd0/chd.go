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
	if hIndex < 0 {
		keyIndex = -hIndex + 1
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

func from(kv map[string][]byte) (*CHD, error) {

	nBuckets := len(kv)
	keys := make([]string, len(kv))
	values := make([][]byte, len(kv))
	hashes := make([]int32, len(kv))

	if nBuckets == 0 {
		nBuckets = 1
	}

	// 1. assign to buckets
	buckets := make([][]string, nBuckets)
	for key := range kv {
		bucketIndex := int(hash(key, 0)) % nBuckets
		buckets[bucketIndex] = append(buckets[bucketIndex], key)
	}

	// sort bucket by length, wanna start w/ larger bucket first.
	sort.Slice(buckets, func(i, j int) bool {
		return len(buckets[i]) > len(buckets[j])
	})

	// 2. assign keys to the right place
	for ithBucket, bucket := range buckets {
		if len(bucket) == 1 {
			break
		}

		r := uint32(1)

		var leftOverKeys []string
		for len(bucket) != 0 {
			key := bucket[0]
			bucket = bucket[1:]

			keyIndex := int(hash(key, r)) % len(keys)
			if len(keys[keyIndex]) != 0 {
				leftOverKeys = append(leftOverKeys, key)
			} else {
				hashes[ithBucket] = int32(r)
				keys[keyIndex] = key
				values[keyIndex] = kv[key]
			}

			if len(bucket) == 0 {
				for _, leftOverKey := range leftOverKeys {
					bucket = append(bucket, leftOverKey)
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

	for i, bucket := range buckets {
		if len(bucket) != 1 {
			continue
		}
		if len(bucket) == 0 {
			break
		}
		slotIndex := freeSlots[0]
		freeSlots = freeSlots[1:]
		hashes[i] = int32(-slotIndex - 1)
		keys[slotIndex] = bucket[0]
		values[slotIndex] = kv[bucket[0]]
	}
	return &CHD{
		keys:   keys,
		values: values,
		h:      hashes,
	}, nil
}
