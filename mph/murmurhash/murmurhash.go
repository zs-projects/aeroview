package murmurhash

import (
	"github.com/twmb/murmur3"
)

// Hash64 returns 64 bits hash of a string.
func Hash64(s string) uint64 {
	// Using unsafe here to avoid escaping to heap / copying bytes on heap.
	// similar to reinterpret_cast in c++
	return murmur3.StringSum64(s)
}

// Hash128 returns 64 bits hash of a string.
func Hash128(s string) (uint64, uint64) {
	// Using unsafe here to avoid escaping to heap / copying bytes on heap.
	// similar to reinterpret_cast in c++
	return murmur3.StringSum128(s)
}
