package memtable

import (
	"golang.org/x/sys/unix"
	"os"
	"syscall"
)

type MemoryMappedMemtable struct {
	MemtableHeader
	data   []byte
	file   *os.File
	offset int
}

func (m *MemoryMappedMemtable) Get(key string) ([]byte, bool) {
	if data, ok := m.getUnsafe(key); ok {
		out := make([]byte, len(data))
		copy(out, data)
		return out, true
	}
	return nil, false
}

func (m *MemoryMappedMemtable) getUnsafe(key string) ([]byte, bool) {
	if lo, ok := m.MemtableHeader.findKey(key); ok {
		return m.data[m.offset:][m.MemtableHeader.offsets[lo]:m.MemtableHeader.offsets[lo+1]], true
	}
	return nil, false
}

func (m *MemoryMappedMemtable) Close() error {
	if err := syscall.Munmap(m.data); err != nil {
		return err
	}
	return m.file.Close()
}
