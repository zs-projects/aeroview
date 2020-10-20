package benchmarks

import (
	"testing"

	"github.com/zs-projects/aeroview/analysis/randutils"
	"github.com/zs-projects/aeroview/rank"
)

var dummyInt int

func BenchRank0(b *testing.B, sliceSize, indexSize int) {
	slc, indexes := randutils.PrepareSliceAndIndexes(sliceSize, indexSize)
	rk := rank.NewRank0(slc)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, idx := range indexes {
			dummy = rk.Get(idx)
		}
	}
	b.ReportMetric(float64(indexSize), "Rank/op")
}

func BenchRank1(b *testing.B, sliceSize, indexSize int) {
	slc, indexes := randutils.PrepareSliceAndIndexes(sliceSize, indexSize)
	rk := rank.NewRank1(slc)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, idx := range indexes {
			dummy = rk.Get(idx)
		}
	}
	b.ReportMetric(float64(indexSize), "Rank/op")
}

func BenchRankPopCount(b *testing.B, sliceSize, indexSize int) {
	slc, indexes := randutils.PrepareSliceAndIndexes(sliceSize, indexSize)
	rk := rank.MakePopCount(slc)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, idx := range indexes {
			dummyInt = rk.Rank(idx)
		}
	}
	b.ReportMetric(float64(indexSize), "Rank/op")
}
func BenchmarkRankV01M(b *testing.B) {
	sliceSize := 1000000
	indexSize := 1000
	BenchRank0(b, sliceSize, indexSize)
}
func BenchmarkRankV0100K(b *testing.B) {
	sliceSize := 100000
	indexSize := 1000
	BenchRank0(b, sliceSize, indexSize)
}

func BenchmarkRankV010K(b *testing.B) {
	sliceSize := 10000
	indexSize := 1000
	BenchRank0(b, sliceSize, indexSize)
}

func BenchmarkRankV11M(b *testing.B) {
	sliceSize := 1000000
	indexSize := 1000
	BenchRank1(b, sliceSize, indexSize)
}
func BenchmarkRankV1100K(b *testing.B) {
	sliceSize := 100000
	indexSize := 1000
	BenchRank1(b, sliceSize, indexSize)
}

func BenchmarkRankV110K(b *testing.B) {
	sliceSize := 10000
	indexSize := 1000
	BenchRank1(b, sliceSize, indexSize)
}

func BenchmarkRankPopcount1M(b *testing.B) {
	sliceSize := 1000000
	indexSize := 1000
	BenchRankPopCount(b, sliceSize, indexSize)
}
func BenchmarkRankPopCount100K(b *testing.B) {
	sliceSize := 100000
	indexSize := 1000
	BenchRankPopCount(b, sliceSize, indexSize)
}

func BenchmarkRankPopCount10K(b *testing.B) {
	sliceSize := 10000
	indexSize := 1000
	BenchRankPopCount(b, sliceSize, indexSize)
}
