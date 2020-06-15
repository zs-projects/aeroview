package encoding

import "testing"

func TestEncodeDecode64(t *testing.T) {
	code := EliasEncoding{}
	for i := 1; i < 200; i++ {
		v, size := code.Encode64(uint64(i))
		l, err := code.Decode64(v, size)
		if err != nil {
			t.Errorf("Got unexpected error while decoding %v", err)
		}
		if l != uint64(i) {
			t.Errorf("Decoding failed was expecting %v got %v", i, 0)
		}

	}
}
