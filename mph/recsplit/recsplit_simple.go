package recsplit

// Simple represents a minimal perfect hash function found using the recsplit algorithm.
type Simple struct {
	values [][]byte
	// One binary tree per bucket.
	keys    []*Node
	cumSums []int
}

func FromRecsplit(r Recsplit) Simple {
	keys := make([]*Node, 0, len(r.keys))
	for _, v := range r.keys {
		keys = append(keys, FromCompactFBTree(v))
	}
	return Simple{
		keys:    keys,
		cumSums: r.cumSums,
		values:  r.values,
	}
}

// GetKey Returns the looked for value
func (r Simple) GetKey(s string) int {
	// 1. Determine bucket.
	bucket := hash(s, rNot) % len(r.keys)
	node := r.keys[bucket]
	h := r.cumSums[bucket]
	out := -1
	for node != nil {
		nbKeys := node.nbKeys
		R := node.R
		split := hash(s, uint64(R)) % nbKeys
		halfNbKeys := nbKeys / 2
		out = h + split
		if split < halfNbKeys {
			node = node.Left
		} else {
			node = node.Right
			h += halfNbKeys
		}
	}
	return out
}

// Get Returns the looked for value
func (r Simple) Get(s string) []byte {
	return r.values[r.GetKey(s)]
}
