package recsplit

import (
	"bytes"
	"fmt"
	"testing"
)

func TestFromMapEF(t *testing.T) {
	nbKeys := 50000
	mp := make(map[string][]byte, nbKeys)
	for i := 0; i < nbKeys; i++ {
		data := fmt.Sprintf("test-%v", i)
		mp[data] = []byte(data)
	}
	mph := CompressWithEliasFano(FromMap(mp, 20))
	collisions := make([]bool, len(mp))
	for key, value := range mp {
		i := mph.GetKey(key)
		if i >= len(collisions) {
			t.Errorf("Index out of range, did not learn an mph %v", i)
		} else {
			if i < len(collisions) && collisions[i] {
				t.Errorf("Found collision for value %v that hashed at %v", key, i)
			} else {
				collisions[i] = true
			}
		}
		if !bytes.Equal(mph.Get(key), value) {
			t.Errorf("Values assignement is corrupt, expecting %v, got %v", string(value), string(mph.Get(key)))
		}
	}
}
