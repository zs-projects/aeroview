package benchmarks

import (
	"fmt"
	"math/rand"
)

func makeDataset(size int, prefix string, subSetSize int) (map[string][]byte, []string) {
	out := make(map[string][]byte, size)
	keysSample := make([]string, 0, subSetSize)
	for i := 0; i < size; i++ {
		key := fmt.Sprintf("%s-key-%v", prefix, i)
		value := fmt.Sprintf("%s-value-%v", prefix, i)
		out[key] = []byte(value)
	}
	for i := 0; i < subSetSize; i++ {
		k := rand.Int() % subSetSize
		key := fmt.Sprintf("%s-key-%v", prefix, k)
		keysSample = append(keysSample, key)
	}
	return out, keysSample
}
