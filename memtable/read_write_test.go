package memtable

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemtableReadWrite(t *testing.T) {
	var buf bytes.Buffer
	size := 1000
	mp := getSampleMap(size)
	memtable := FromMap(mp)
	err := Write(&memtable, &buf)
	if err != nil {
		t.Fatalf(err.Error())
	}
	m := MemtableFromReader(&buf)
	assert.EqualValues(t, m.keys, memtable.keys)
	assert.EqualValues(t, m.offsets, memtable.offsets)
	assert.EqualValues(t, m.data, memtable.data)
}

func TestMmapMemtableRead(t *testing.T) {
	f, err := ioutil.TempFile("", "mmap_memtable")
	if err != nil {
		t.Fail()
	}
	size := 5000
	mp := getSampleMap(size)
	memtable := FromMap(mp)
	err = Write(&memtable, f)
	if err != nil {
		t.Fatalf(err.Error())
	}
	f.Seek(0, os.SEEK_SET)
	m := MMapMemtableFromFile(f)
	assert.EqualValues(t, m.keys, memtable.keys)
	assert.EqualValues(t, m.offsets, memtable.offsets)
	assert.Equal(t, len(m.data)-m.offset, len(memtable.data))
	assert.EqualValues(t, m.data[m.offset:], memtable.data)
	for key, value := range mp {
		if v, ok := m.Get(key); ok {
			assert.EqualValues(t, v, value)
		}
	}
}
