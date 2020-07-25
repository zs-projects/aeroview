package rank

import "testing"

func TestCount(t *testing.T) {
	for i := uint8(1); i < 255; i++ {
		c1 := count(int(i))
		c2 := count2(int(i))
		c3 := count3(i)
		if c1 != c2 {
			t.Errorf("Count2 failed for %v with answer %v, was expecting %v", i, c1, c2)
		}
		if c1 != c3 {
			t.Errorf("Count2 failed for %v with answer %v, was expecting %v", i, c1, c3)
		}
	}
}
