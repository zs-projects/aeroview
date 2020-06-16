package chd0

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCHD_Get(t *testing.T) {
	nValues := 100

	m := make(map[string][]byte)
	for i := 0; i < nValues; i++ {
		m[fmt.Sprintf("key=%d", i)] = []byte(fmt.Sprintf("value=%d", i))
	}

	chd, err := From(m)
	assert.Nil(t, err)

	for i := 0; i < nValues; i++ {
		value, hasVal := chd.Get(fmt.Sprintf("key=%d", i))
		assert.Equal(t, true, hasVal)
		assert.Equal(t, []byte(fmt.Sprintf("value=%d", i)), value, "fail to assert arg %d", i)
	}
}
