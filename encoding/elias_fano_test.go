package encoding

import (
	"bytes"
	"testing"
)

func TestEliasFanoVector(t *testing.T) {
	//01011010000010100001100
	values := []uint64{5, 8, 8, 15, 32}
	vec := MakeEliasFanoVector(values)
	val, length := vec.Data()
	if !bytes.Equal(val, []byte{90, 10, 24}) || length != 23 {
		t.Errorf("Elias Fano encoding failed.%v \t %v \n", val, length)
	}
	for i, v := range values {
		if out := vec.Get(i); out != v {
			t.Errorf("Get(%v) method for Elias Fano encoding failed got %v \t expected %v \n", i, out, v)
		}
	}
}
