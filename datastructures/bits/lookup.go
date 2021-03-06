package bits

import (
	"math/bits"
)

type LookupTable8 = [256]uint8

// Make8BitLookup creates an 8 bit wide lookup table
func Make8BitLookup() LookupTable8 {
	var lookup [256]uint8
	for i := 0; i <= 255; i++ {
		lookup[i] = uint8(bits.OnesCount8(uint8(i) & 255))
	}
	return lookup
}
