package smalllist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVarInt16_Encode(t *testing.T) {
	for i := 0; i < 1000; i++ {
		buffer := make([]uint8, 10)
		Encode(buffer, int32(i))
		val := Decode(buffer, 0)
		assert.Equal(t, int32(i), val)
	}
}
