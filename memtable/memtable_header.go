package memtable

import (
	"math"

	"github.com/zs-projects/aeroview/mph/farmhash"
)

type MemtableHeader struct {
	keys    []uint64
	offsets []uint32
}

func (m *MemtableHeader) findKey(key string) (int, bool) {
	hashedKey := farmhash.Hash64(key)
	lo := 0
	hi := len(m.keys) - 1

	if len(m.keys) == 0 || hi == -1 {
		return -1, false
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
		return -1, false
	}
	return lo, true
}

func (m *MemtableHeader) linearSearch(lo, high int, hashedKey uint64) (int, bool) {
	for i := lo; i <= high; i++ {
		if hashedKey == m.keys[i] {
			return i, true
		}
	}
	return -1, false
}

func (m *MemtableHeader) slope(lo, hi int) float64 {
	last := float64(m.keys[hi])
	first := float64(m.keys[lo])
	span := float64(len(m.keys))
	return span / (last - first)
}
