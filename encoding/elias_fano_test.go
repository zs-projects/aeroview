package encoding

import "testing"

func TestEncodeDecode64(t *testing.T) {
	for i := 1; i < 200; i++ {
		v, size := Encode64([]uint64{uint64(i)})
		l, err := Decode64(v, size)
		if err != nil {
			t.Errorf("Got unexpected error while decoding %v", err)
		}
		if len(l) != 1 || l[0] != uint64(i) {
			t.Errorf("Decoding failed was expecting a size of 1 got %v and a value of %v got %v", len(l), i, l[0])
		}

	}
}
