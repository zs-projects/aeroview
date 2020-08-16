package benchmarks

import (
	"math/rand"
	"testing"

	"zs-project.org/aeroview/analysis/arrays/randutils"
	"zs-project.org/aeroview/encoding"
)

func BenchmarkEliasFanoGet100K(b *testing.B) {
	slc := randutils.RandSlice64(100000)
	elias := encoding.MakeEliasFanoVector(slc)
	indexes := make([]int, b.N)
	for i := range indexes {
		indexes[i] = rand.Intn(100000)
	}
	b.ResetTimer()
	for _, idx := range indexes {
		elias.Get(idx)
	}
}

func BenchmarkSliceGet100K(b *testing.B) {
	slc := randutils.RandSlice64(100000)
	indexes := make([]int, b.N)
	for i := range indexes {
		indexes[i] = rand.Intn(100000)
	}
	b.ResetTimer()
	for _, idx := range indexes {
		_ = slc[idx]
	}
}

func BenchmarkEliasFanoGet10K(b *testing.B) {
	slc := randutils.RandSlice64(10000)
	elias := encoding.MakeEliasFanoVector(slc)
	indexes := make([]int, b.N)
	for i := range indexes {
		indexes[i] = rand.Intn(10000)
	}
	b.ResetTimer()
	for _, idx := range indexes {
		elias.Get(idx)
	}
}

func BenchmarkSliceGet10K(b *testing.B) {
	slc := randutils.RandSlice64(10000)
	indexes := make([]int, b.N)
	for i := range indexes {
		indexes[i] = rand.Intn(10000)
	}
	b.ResetTimer()
	for _, idx := range indexes {
		_ = slc[idx]
	}
}
