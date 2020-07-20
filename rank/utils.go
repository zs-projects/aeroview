package rank

func Make8BitLookup() [256]uint8 {
	var lookup [256]uint8
	for i := 0; i < 256; i++ {
		lookup[i] = uint8(count(i & 255))
	}
	return lookup
}

func count(n int) int {
	nBits := 0
	for n != 0 {
		if n&1 == 1 {
			nBits++
		}
		n >>= 1
	}
	return nBits
}
