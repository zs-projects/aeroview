package memtable

import (
	"math"
	"sort"

	"github.com/zs-projects/aeroview/mph/farmhash"
)

type Memtable struct {
	originalKeys []string
	keys         []uint64
	offsets      []uint32
	data         []byte
}

// FromMap builds a Memtable from a map.
func FromMap(data map[string][]byte) Memtable {
	originalKeys := make([]string, 0, len(data))
	keys := make([]uint64, 0, len(data))
	nbBytes := 0
	for key, data := range data {
		originalKeys = append(originalKeys, key)
		keys = append(keys, farmhash.Hash64(key))
		nbBytes += len(data)

	}
	sort.Slice(originalKeys, func(i, j int) bool { return farmhash.Hash64(originalKeys[i]) < farmhash.Hash64(originalKeys[j]) })
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	// still suboptimal becase data can be unicode
	dataBytes := make([]byte, 0, nbBytes)
	offsets := make([]uint32, 1, len(data)-1)
	currOffset := uint32(0)
	for _, stringKey := range originalKeys {
		value := data[stringKey]
		currOffset += uint32(len(value))
		offsets = append(offsets, currOffset)
		dataBytes = append(dataBytes, value...)
	}

	return Memtable{
		originalKeys: originalKeys,
		keys:         keys,
		offsets:      offsets,
		data:         dataBytes,
	}
}

func (m *Memtable) MapView() map[string][]byte {
	return nil
}

func (m *Memtable) Get(key string) ([]byte, bool) {
	if data, ok := m.GetUnsafe(key); ok {
		out := make([]byte, len(data))
		copy(out, data)
		return out, true
	}
	return nil, false
}

func (m *Memtable) slope(lo, hi int) float64 {
	last := float64(m.keys[hi])
	first := float64(m.keys[lo])
	span := float64(len(m.keys))
	return span / (last - first)
}

func (m *Memtable) linearSearch(lo, high int, hashedKey uint64) ([]byte, bool) {
	for i := lo; i <= high; i++ {
		if hashedKey == m.keys[i] {
			return m.data[m.offsets[i]:m.offsets[i+1]], true
		}
	}
	return nil, false
}

func (m *Memtable) GetUnsafe(key string) ([]byte, bool) {
	hashedKey := farmhash.Hash64(key)
	lo := 0
	hi := len(m.keys) - 1

	if len(m.keys) == 0 || hi == -1 {
		return nil, false
	}
	slope := m.slope(lo, hi)

	for m.keys[lo] < hashedKey {
		midFP := slope * float64(hashedKey-m.keys[lo])
		mid := lo + int(math.Min(float64(hi-lo-1), midFP))

		if m.keys[mid] < hashedKey {
			lo = mid + 1
		} else {
			hi = mid
		}

		if mid+64 > hi || mid-64 <= lo {
			return m.linearSearch(lo, hi, hashedKey)
		}
	}

	if m.keys[lo] != hashedKey {
		return nil, false
	}
	return m.data[m.offsets[lo]:m.offsets[lo+1]], true
}
