package benchmarks

import (
	"testing"

	"github.com/zs-projects/aeroview/analysis/randutils"
	"github.com/zs-projects/aeroview/rank"
)

var dummy uint64

func BenchmarkSelectPop100K(b *testing.B) {
	slc, indexes := randutils.PrepareSliceAndIndexes(100000, 1000)
	rk := rank.MakePopCount(slc)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, idx := range indexes {
			dummy = rk.Select(uint64(idx))
		}
	}
	b.ReportMetric(float64(1000), "Select/op")
}

func BenchmarkSelectPop10K(b *testing.B) {
	slc, indexes := randutils.PrepareSliceAndIndexes(10000, 1000)
	rk := rank.MakePopCount(slc)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, idx := range indexes {
			dummy = rk.Select(uint64(idx))
		}
	}
	b.ReportMetric(float64(1000), "Select/op")
}
