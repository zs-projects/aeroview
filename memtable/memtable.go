package memtable

import (
	"sort"

	"github.com/zs-projects/aeroview/mph/farmhash"
)

type Memtable struct {
	MemtableHeader
	data []byte
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
		MemtableHeader: MemtableHeader{
			keys:    keys,
			offsets: offsets,
		},
		data: dataBytes,
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

func (m *Memtable) GetUnsafe(key string) ([]byte, bool) {
	if lo, ok := m.MemtableHeader.findKey(key); ok {
		return m.data[m.MemtableHeader.offsets[lo]:m.MemtableHeader.offsets[lo+1]], true
	}
	return nil, false
}
