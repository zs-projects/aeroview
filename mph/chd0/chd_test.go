package chd0

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCHD_Get(t *testing.T) {

	m := make(map[string][]byte)
	for i := 0; i < 10; i++ {
		m[fmt.Sprintf("key=%d", i)] = []byte(fmt.Sprintf("value=%d", i))
	}

	chd, _ := from(m)

	for i := 0; i < 10; i++ {
		value, _ := chd.Get(fmt.Sprintf("key=%d", i))
		fmt.Println(i, string(value))
		assert.Equal(t, []byte(fmt.Sprintf("value=%d", i)), value, "fail to assert arg %d", i)
	}

}
