package rank

import (
	"math"
	"math/bits"
)

func count(n int) int {
	nBits := 0
	for n != 0 {
		n = n & (n - 1)
		nBits++
	}
	return nBits

}

func count3(n uint8) int {
	return bits.OnesCount8(n)
}

func count2(n int) int {
	nBits := 0
	for n != 0 {
		if n&1 == 1 {
			nBits++
		}
		n >>= 1
	}
	return nBits
}

func nbBits_floor(a int) int {
	return int(math.Floor(math.Log2(float64(a))))
}

func ceil_devide(a, b int) int {
	return int(math.Ceil(float64(a / b)))
}
func floor_devide(a, b int) int {
	return int(math.Floor(float64(a / b)))
}
