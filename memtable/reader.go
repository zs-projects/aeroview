package memtable

import (
	"encoding/binary"
	"io"
	"io/ioutil"
)

func FromReader(r io.Reader) Memtable {
	return NewReader(r).ReadHeader().ReadKeysAndOffsets().ReadValues()
}

type memtableReader struct{ reader io.Reader }

func NewReader(r io.Reader) memtableReader {
	return memtableReader{reader: r}
}

func (m memtableReader) ReadHeader() memtablePartialHeader {
	buf := make([]byte, 8)
	n, _ := m.reader.Read(buf)
	if n != 8 {
		panic("can't read headerLen")
	}
	headerLen := binary.LittleEndian.Uint64(buf)
	n, _ = m.reader.Read(buf)
	if n != 8 {
		panic("can't read nbKeys")
	}
	nbKeys := binary.LittleEndian.Uint64(buf)
	return memtablePartialHeader{
		reader:    m.reader,
		headerLen: headerLen,
		nbKeys:    nbKeys,
	}

}

type memtablePartialHeader struct {
	reader    io.Reader
	headerLen uint64
	nbKeys    uint64
}

func (m memtablePartialHeader) ReadKeysAndOffsets() memtableHeader {
	keyByteLength := 8 * m.nbKeys
	buf := make([]byte, keyByteLength)
	n, err := m.reader.Read(buf)
	if n != int(keyByteLength) || err != nil {
		panic("can't read keys")
	}
	keys := make([]uint64, 0, m.nbKeys)
	for i := uint64(0); i < m.nbKeys; i++ {
		key := binary.LittleEndian.Uint64(buf[i*8 : (i+1)*8])
		keys = append(keys, key)
	}

	offsetsByteLength := 4*m.nbKeys + 4
	buf = make([]byte, offsetsByteLength)
	n, err = m.reader.Read(buf)
	if n != int(offsetsByteLength) || err != nil {
		panic("can't read offsets")
	}
	offsets := make([]uint32, 0, m.nbKeys+1)
	for i := uint64(0); i <= m.nbKeys; i++ {
		offset := binary.LittleEndian.Uint32(buf[i*4 : (i+1)*4])
		offsets = append(offsets, offset)
	}
	return memtableHeader{
		reader:    m.reader,
		headerLen: m.headerLen,
		keys:      keys,
		offsets:   offsets,
	}
}

type memtableHeader struct {
	reader    io.Reader
	headerLen uint64
	keys      []uint64
	offsets   []uint32
}

func (m memtableHeader) ReadValues() Memtable {
	values, err := ioutil.ReadAll(m.reader)
	if err != nil {
		panic("can't read values")
	}
	return Memtable{
		originalKeys: nil,
		data:         values,
		keys:         m.keys,
		offsets:      m.offsets,
	}
}
