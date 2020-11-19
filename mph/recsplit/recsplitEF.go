package recsplit

import (
	"github.com/zs-projects/aeroview/encoding"
)

// Recsplit represents a minimal perfect hash function found using the recsplit algorithm.
type RecsplitEF struct {
	values [][]byte
	// One binary tree per bucket.
	keys    []CompactFBTree
	cumSums encoding.EliasFanoVector
}

func CompressWithEliasFano(r Recsplit) RecsplitEF {
	u := make([]uint64, 0, len(r.cumSums))
	for _, v := range r.cumSums {
		u = append(u, uint64(v))
	}
	return RecsplitEF{
		values:  r.values,
		keys:    r.keys,
		cumSums: encoding.MakeEliasFanoVector(u),
	}
}

func (r RecsplitEF) SizeInBytes() int {
	u := 0
	for _, fbt := range r.keys {
		u += fbt.SizeInBytes()
	}
	return r.cumSums.Len()/64*8 + u
}

// GetKey Returns the looked for value
func (r RecsplitEF) GetKey(s string) int {
	// 1. Determine bucket.
	bucket := hash(s, rNot) % len(r.keys)
	tree := r.keys[bucket]
	node := tree.Root()
	h := r.cumSums.Get(bucket)
	out := -1
	for node != 0 || (node == 0 && out == -1) {
		R, nbKeys := tree.node(node)
		split := hash(s, uint64(R)) % nbKeys
		halfNbKeys := nbKeys / 2
		out = int(h) + split
		if split < halfNbKeys {
			node = tree.LeftChild(node)
		} else {
			node = tree.RightChild(node)
			h += uint64(halfNbKeys)
		}
	}
	return out
}

// Get Returns the looked for value
func (r RecsplitEF) Get(s string) []byte {
	return r.values[r.GetKey(s)]
}
