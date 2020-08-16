package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"zs-project.org/aeroview/analysis/arrays/randutils"
	"zs-project.org/aeroview/encoding"
)

// Observations :
// * Elias fano seems to consistently outperfom VarInt encoding on the samples where all the numbers where encoded
// on the same number of bits.
// * On the cases where you have a few 8 bytes numbers and many 4 and 2 bytes numbers, varInt seems to outperform elias-fano.

func main() {
	sizes := []int{10, 100, 1000, 1000, 10000, 100000, 1000000}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "Experiment ( 64 bits )\tOriginal\tVarInt\tElias-Fano\tImprovement\t\n")
	for _, size := range sizes {
		slc := randutils.RandSlice64(size)
		varIntBuf := randutils.VarInt64(slc)
		eliasFanoBuf := encoding.MakeEliasFanoVector(slc)
		nbBytesOrig := len(slc) * 64 / 8
		nbBytesVarInt := len(varIntBuf)
		nbBytesEliasFano := eliasFanoBuf.Len() / 8
		fmt.Fprintf(w, "nb elements: %v\t%v\t%v\t%v\t%.3f\t\n", size, nbBytesOrig, nbBytesVarInt, nbBytesEliasFano, float64(nbBytesVarInt)/float64(nbBytesEliasFano))
	}
	w.Flush()
	w = tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "Experiment ( 32 bits )\tOriginal\tVarInt\tElias-Fano\tImprovement\t\n")
	for _, size := range sizes {
		slc := randutils.RandSlice32(size)
		varIntBuf := randutils.VarInt32(slc)
		eliasFanoBuf := encoding.MakeEliasFanoVector(randutils.Slice32To64(slc))
		nbBytesOrig := len(slc) * 32 / 8
		nbBytesVarInt := len(varIntBuf)
		nbBytesEliasFano := eliasFanoBuf.Len() / 8
		fmt.Fprintf(w, "nb elements: %v\t%v\t%v\t%v\t%.3f\t\n", size, nbBytesOrig, nbBytesVarInt, nbBytesEliasFano, float64(nbBytesVarInt)/float64(nbBytesEliasFano))
	}
	w.Flush()

	w = tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "Experiment ( 16 bits )\tOriginal\tVarInt\tElias-Fano\tImprovement\t\n")
	for _, size := range sizes {
		slc := randutils.RandSlice16(size)
		varIntBuf := randutils.VarInt16(slc)
		eliasFanoBuf := encoding.MakeEliasFanoVector(randutils.Slice16To64(slc))
		nbBytesOrig := len(slc) * 16 / 8
		nbBytesVarInt := len(varIntBuf)
		nbBytesEliasFano := eliasFanoBuf.Len() / 8
		fmt.Fprintf(w, "nb elements: %v\t%v\t%v\t%v\t%.3f\t\n", size, nbBytesOrig, nbBytesVarInt, nbBytesEliasFano, float64(nbBytesVarInt)/float64(nbBytesEliasFano))
	}
	w.Flush()

	w = tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "Experiment ( Mixed )\tVarInt\tElias-Fano\tImprovement\t\n")
	for _, size := range sizes {
		slc := randutils.RandSliceMixed(size)
		varIntBuf := randutils.VarInt64(slc)
		eliasFanoBuf := encoding.MakeEliasFanoVector(slc)
		nbBytesVarInt := len(varIntBuf)
		nbBytesEliasFano := eliasFanoBuf.Len() / 8
		fmt.Fprintf(w, "nb elements: %v\t%v\t%v\t%.3f\t\n", size, nbBytesVarInt, nbBytesEliasFano, float64(nbBytesVarInt)/float64(nbBytesEliasFano))
	}
	w.Flush()

}
