package chd0

import (
	"fmt"
	"testing"
)

func TestCHD_Get(t *testing.T) {

	m := make(map[string][]byte)
	for i := 0; i < 10; i++ {
		m[fmt.Sprintf("key=%d", i)] = []byte(fmt.Sprintf("value=%d", i))
	}

	chd, _ := from(m)
	fmt.Println(chd.Get("key=0"))
}
