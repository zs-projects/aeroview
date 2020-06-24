package utils

type Bucket struct {
	OriginalIndex int
	Keys          []string
}

func NewBucket(index int) *Bucket {
	return &Bucket{
		OriginalIndex: index,
		Keys:          nil,
	}
}

func AssignToBuckets(assigner IndexHashFunc, data map[string][]byte, nBuckets int) Buckets {
	buckets := make(Buckets, nBuckets)
	for key := range data {
		bucketIndex := assigner(key) % nBuckets
		if buckets[bucketIndex] == nil {
			buckets[bucketIndex] = NewBucket(bucketIndex)
		}
		buckets[bucketIndex].Keys = append(buckets[bucketIndex].Keys, key)
	}
	return buckets
}

type Buckets []*Bucket

// Implements the Sort interface.
func (b Buckets) Len() int {
	return len(b)
}

func (b Buckets) Less(i, j int) bool {
	if b[i] != nil && b[j] != nil {
		return len(b[i].Keys) < len(b[j].Keys)
	} else if b[i] != nil {
		return false
	}
	return true
}

func (b Buckets) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
