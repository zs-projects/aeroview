package benchmarks

import (
	"testing"

	"github.com/zs-projects/aeroview/analysis/randutils"
	"github.com/zs-projects/aeroview/encoding"
)

var dummy uint64

func benchEliasFanoGet(b *testing.B, sliceSize, indexSize int) {
	slc, indexes := randutils.PrepareSliceAndIndexes(sliceSize, indexSize)
	elias := encoding.MakeEliasFanoVector(slc)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, idx := range indexes {
			dummy = elias.Get(idx)
		}
	}
	b.ReportMetric(float64(indexSize), "Get/op")
}

func benchSliceGet(b *testing.B, sliceSize, indexSize int) {
	slc, indexes := randutils.PrepareSliceAndIndexes(sliceSize, indexSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, idx := range indexes {
			dummy = slc[idx]
		}
	}
	b.ReportMetric(float64(indexSize), "Get/op")
}
func BenchmarkEliasFanoGet100K(b *testing.B) {
	sliceSize := 100000
	indexSize := 1000
	benchEliasFanoGet(b, sliceSize, indexSize)
}
func BenchmarkEliasFanoGet10K(b *testing.B) {
	sliceSize := 10000
	indexSize := 1000
	benchEliasFanoGet(b, sliceSize, indexSize)
}

func BenchmarkSliceGet100K(b *testing.B) {
	sliceSize := 100000
	indexSize := 1000
	benchSliceGet(b, sliceSize, indexSize)
}
func BenchmarkSliceGet10K(b *testing.B) {
	sliceSize := 10000
	indexSize := 1000
	benchSliceGet(b, sliceSize, indexSize)
}
