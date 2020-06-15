package encoding

import (
	"bytes"
	"testing"
)

func TestEliasFanoEncode64(t *testing.T) {
	//01011010000010100001100
	code := EliasFanoEncoding{}
	values := []uint64{5, 8, 8, 15, 32}
	val, length := code.Encode64(values)
	if !bytes.Equal(val, []byte{90, 10, 24}) || length != 23 {
		t.Errorf("Elias Fano encoding failed.%v \t %v \n", val, length)
	}
}
