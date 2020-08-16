package encoding

import (
	"testing"
)

func TestEliasFanoVector(t *testing.T) {
	//0101101000001 0100001100
	values := []uint64{5, 8, 8, 15, 32}
	vec := MakeEliasFanoVector(values)
	highBits := vec.highBits.Vector()[0]
	lowBits := vec.lowBits.Vector()[0]
	if highBits != 0b1000001011010 || lowBits != 0b0011000001 {
		t.Errorf("Elias Fano encoding failed")
	}
	for i, v := range values {
		if out := vec.Get(i); out != v {
			t.Errorf("Get(%v) method for Elias Fano encoding failed got %v \t expected %v \n", i, out, v)
		}
	}
}
