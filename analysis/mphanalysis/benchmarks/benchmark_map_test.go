package benchmarks

import (
	"testing"

	"github.com/zs-projects/aeroview/mph/chd0"
	"github.com/zs-projects/aeroview/mph/recsplit"
)

var dummy []byte

func BenchmarkMapGet100K(b *testing.B) {
	ds, keys := makeDataset(100000, "string", 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			dummy = ds[key]
		}
	}
	b.ReportMetric(float64(1000), "Get/op")
}

func BenchmarkRecsplitGet100K(b *testing.B) {
	ds, keys := makeDataset(100000, "string", 1000)

	mph := recsplit.FromMap(ds, 20)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			dummy = mph.Get(key)
		}
	}
	b.ReportMetric(float64(1000), "Get/op")
}

func BenchmarkRecsplitSimpleGet100K(b *testing.B) {
	ds, keys := makeDataset(100000, "string", 1000)

	mph := recsplit.Recsplit(recsplit.FromMap(ds, 20))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			dummy = mph.Get(key)
		}
	}
	b.ReportMetric(float64(1000), "Get/op")
}
func BenchmarkCHDGet100K(b *testing.B) {
	ds, keys := makeDataset(100000, "string", 1000)

	mph, err := chd0.FromMap(ds)
	if err != nil {
		panic("CHD failed")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			dummy, _ = mph.Get(key)
		}
	}
	b.ReportMetric(float64(1000), "Get/op")
}
