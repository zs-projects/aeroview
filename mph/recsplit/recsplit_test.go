package recsplit

import (
	"testing"
)

func TestBruteForceMPH(t *testing.T) {
	b := recsplitBucket{
		keys:    []string{"toto", "tata", "titi", "test", "tardif"},
		parents: nil,
	}
	r := b.bruteForceMPH()
	collisions := make([]bool, len(b.keys))
	if checkForCollisions(uint32(r), b.keys, collisions) {
		t.Errorf("Found a collision!")
	}
}
func TestCheckCollisions(t *testing.T) {
	keys := []string{"toto", "tata", "titi", "test", "tardif"}
	collisions := make([]bool, len(keys))
	if !checkForCollisions(uint32(0), keys, collisions) {
		t.Errorf("checkCollisions Failed")
	}
}
