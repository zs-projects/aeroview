package memtable

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"syscall"
)

var (
	ErrHeaderRead = fmt.Errorf("can't read header")
)

func MemtableFromReader(r io.Reader) Memtable {
	return NewReader(r).
		ReadHeaderMetadata().
		ReadKeysAndOffsets().
		ReadValues()
}

func MMapMemtableFromFile(f *os.File) MemoryMappedMemtable {

	return NewReader(f).
		ReadHeaderMetadata().
		ReadKeysAndOffsets().
		MemoryMapValues(f)
}

func NewReader(r io.Reader) readMemtableHeader {
	return readMemtableHeader{r}
}

type readMemtableHeader struct{ io.Reader }

func (m readMemtableHeader) ReadHeaderMetadata() readMemtableMetadata {
	buf := make([]byte, 8)
	n, err := m.Read(buf)
	if n != 8 {
		panic("can't read header size ")
	}
	if err != nil {
		panic("can't read header")
	}
	headerLen := binary.LittleEndian.Uint64(buf)
	n, _ = m.Read(buf)
	if n != 8 {
		panic("can't read nbKeys")
	}
	nbKeys := binary.LittleEndian.Uint64(buf)
	return readMemtableMetadata{
		Reader:    m,
		headerLen: headerLen,
		nbKeys:    nbKeys,
	}
}

type readMemtableMetadata struct {
	io.Reader
	headerLen uint64
	nbKeys    uint64
}

func (m readMemtableMetadata) ReadKeysAndOffsets() valuesReader {
	keyByteLength := 8 * m.nbKeys
	buf := make([]byte, keyByteLength)
	n, err := m.Read(buf)
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
	n, err = m.Read(buf)
	if n != int(offsetsByteLength) || err != nil {
		panic("can't read offsets")
	}
	offsets := make([]uint32, 0, m.nbKeys+1)
	for i := uint64(0); i <= m.nbKeys; i++ {
		offset := binary.LittleEndian.Uint32(buf[i*4 : (i+1)*4])
		offsets = append(offsets, offset)
	}
	return valuesReader{
		reader:    m,
		headerLen: m.headerLen,
		MemtableHeader: MemtableHeader{
			keys:    keys,
			offsets: offsets,
		},
	}
}

type valuesReader struct {
	MemtableHeader
	reader    io.Reader
	headerLen uint64
}

func (m valuesReader) ReadValues() Memtable {
	values, err := ioutil.ReadAll(m.reader)
	if err != nil {
		panic("can't read values")
	}
	return Memtable{
		MemtableHeader: m.MemtableHeader,
		data:           values,
	}
}

func (m valuesReader) MemoryMapValues(f *os.File) MemoryMappedMemtable {
	fi, err := f.Stat()
	if err != nil {
		panic("can't get stats")
	}
	offset := int64(m.headerLen)
	size := int(fi.Size())
	pageSize := syscall.Getpagesize()
	multiple := (int(m.headerLen) / pageSize)
	offset = int64(multiple * pageSize)
	diff := int64(m.headerLen) - offset

	value, err := syscall.Mmap(int(f.Fd()),
		offset, size-int(offset),
		syscall.PROT_WRITE, syscall.MAP_SHARED)

	if err != nil {
		panic(err.Error())
	}

	return MemoryMappedMemtable{
		MemtableHeader: m.MemtableHeader,
		data:           value,
		file:           f,
		offset:         int(diff + 4), // we add the 4 bytes that we used to indicate the header length at the start of the file
	}
}
