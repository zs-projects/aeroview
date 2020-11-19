package memtable

import (
	"bytes"
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
	m := FromReader(&buf)
	assert.EqualValues(t, m.keys, memtable.keys)
	assert.EqualValues(t, m.offsets, memtable.offsets)
	assert.EqualValues(t, m.data, memtable.data)

}
