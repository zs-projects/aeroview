package memtable

import (
	"encoding/binary"
	"fmt"
	"io"
)

func Write(m Memtable, w io.Writer) error {
	nbKeys := uint64(len(m.keys))
	headerLen := uint64(8 + 8 + 12*nbKeys)
	buf := make([]byte, 16)
	binary.LittleEndian.PutUint64(buf[:8], headerLen)
	binary.LittleEndian.PutUint64(buf[8:], nbKeys)
	n, err := w.Write(buf)
	if n != len(buf) || err != nil {
		return fmt.Errorf("Can't write header")
	}
	buf = buf[:8]
	for _, key := range m.keys {
		binary.LittleEndian.PutUint64(buf, key)
		n, err = w.Write(buf)
		if n != len(buf) || err != nil {
			return fmt.Errorf("Can't write key %v", key)
		}
	}

	buf = buf[:4]
	for _, offset := range m.offsets {
		binary.LittleEndian.PutUint32(buf, offset)
		n, err = w.Write(buf)
		if n != len(buf) || err != nil {
			return fmt.Errorf("Can't write offset %v", offset)
		}
	}
	n, err = w.Write(m.data)
	if n != len(m.data) || err != nil {
		return fmt.Errorf("Can't write data")
	}
	return nil
}
