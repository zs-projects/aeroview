package benchmarks

import (
	"math/rand"
	"testing"

	"zs-project.org/aeroview/analysis/randutils"
	"zs-project.org/aeroview/rank"
)

func BenchmarkRankV01M(b *testing.B) {
	slc := randutils.RandSlice64(1000000)
	rk := rank.NewRank0(slc)
	indexes := make([]int, b.N)
	for i := range indexes {
		indexes[i] = rand.Intn(1000000)
	}
	b.ResetTimer()
	for _, idx := range indexes {
		_ = rk.Get(idx)
	}
}

func BenchmarkRankPop1M(b *testing.B) {
	slc := randutils.RandSlice64(1000000)
	rk := rank.MakePopCount(slc)
	indexes := make([]int, b.N)
	for i := range indexes {
		indexes[i] = rand.Intn(1000000)
	}
	b.ResetTimer()
	for _, idx := range indexes {
		_ = rk.Rank(idx)
	}
}
func BenchmarkRankV0100K(b *testing.B) {
	slc := randutils.RandSlice64(100000)
	rk := rank.NewRank0(slc)
	indexes := make([]int, b.N)
	for i := range indexes {
		indexes[i] = rand.Intn(100000)
	}
	b.ResetTimer()
	for _, idx := range indexes {
		_ = rk.Get(idx)
	}
}

func BenchmarkRankPop100K(b *testing.B) {
	slc := randutils.RandSlice64(100000)
	rk := rank.MakePopCount(slc)
	indexes := make([]int, b.N)
	for i := range indexes {
		indexes[i] = rand.Intn(100000)
	}
	b.ResetTimer()
	for _, idx := range indexes {
		_ = rk.Rank(idx)
	}
}

func BenchmarkRankV010K(b *testing.B) {
	slc := randutils.RandSlice64(10000)
	rk := rank.NewRank0(slc)
	indexes := make([]int, b.N)
	for i := range indexes {
		indexes[i] = rand.Intn(10000)
	}
	b.ResetTimer()
	for _, idx := range indexes {
		_ = rk.Get(idx)
	}
}

func BenchmarkRankPop10K(b *testing.B) {
	slc := randutils.RandSlice64(10000)
	rk := rank.MakePopCount(slc)
	indexes := make([]int, b.N)
	for i := range indexes {
		indexes[i] = rand.Intn(10000)
	}
	b.ResetTimer()
	for _, idx := range indexes {
		_ = rk.Rank(idx)
	}
}
