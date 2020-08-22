package main

import (
	"fmt"
	"os"
	"runtime"
	"text/tabwriter"

	"github.com/zs-projects/aeroview/analysis/mphanalysis/benchmarks"
	"github.com/zs-projects/aeroview/mph/chd0"
	"github.com/zs-projects/aeroview/mph/recsplit"
)

func memUsage(m2, m1 *runtime.MemStats) *runtime.MemStats {
	return &runtime.MemStats{
		Alloc:      m2.Alloc - m1.Alloc,
		TotalAlloc: m2.TotalAlloc - m1.TotalAlloc,
		HeapAlloc:  m2.HeapAlloc - m1.HeapAlloc,
	}
}

var (
	chd     *chd0.CHD
	rcsplit recsplit.Recsplit
)

func main() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "Experiment ( 64 bits )\tCHD\tRecsplit\tOriginal Map\tCHD Saving\tRecsplit Saving\t\n")
	var m0, m1, m2, m3 runtime.MemStats
	sizes := []int{1000, 10000, 100000, 1000000}
	for _, size := range sizes {
		runtime.ReadMemStats(&m0)
		dataset, _ := benchmarks.MakeDataset(size, "string", 100)
		runtime.GC()
		runtime.ReadMemStats(&m1)
		chd, _ = chd0.FromMap(dataset)
		runtime.GC()
		runtime.ReadMemStats(&m2)
		rcsplit = recsplit.FromMap(dataset, 10)
		runtime.GC()
		runtime.ReadMemStats(&m3)
		memUsageMap := memUsage(&m1, &m0)
		memUsageCHD := memUsage(&m2, &m1)
		memUsageRECSPLIT := memUsage(&m3, &m2)
		fmt.Fprintf(w, "Number of elements: %v\t%v\t%v\t%v\t%v\t%v\t\n",
			len(dataset),
			memUsageCHD.HeapAlloc,
			memUsageRECSPLIT.HeapAlloc,
			memUsageMap.HeapAlloc,
			float64(memUsageCHD.HeapAlloc)/float64(memUsageMap.HeapAlloc),
			float64(memUsageRECSPLIT.HeapAlloc)/float64(memUsageMap.HeapAlloc))
	}
	w.Flush()

}
