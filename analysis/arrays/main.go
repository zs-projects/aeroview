package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"text/tabwriter"

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
		slc := randSlice64(size)
		varIntBuf := varInt64(slc)
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
		slc := randSlice32(size)
		varIntBuf := varInt32(slc)
		eliasFanoBuf := encoding.MakeEliasFanoVector(slice32To64(slc))
		nbBytesOrig := len(slc) * 32 / 8
		nbBytesVarInt := len(varIntBuf)
		nbBytesEliasFano := eliasFanoBuf.Len() / 8
		fmt.Fprintf(w, "nb elements: %v\t%v\t%v\t%v\t%.3f\t\n", size, nbBytesOrig, nbBytesVarInt, nbBytesEliasFano, float64(nbBytesVarInt)/float64(nbBytesEliasFano))
	}
	w.Flush()

	w = tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "Experiment ( 16 bits )\tOriginal\tVarInt\tElias-Fano\tImprovement\t\n")
	for _, size := range sizes {
		slc := randSlice16(size)
		varIntBuf := varInt16(slc)
		eliasFanoBuf := encoding.MakeEliasFanoVector(slice16To64(slc))
		nbBytesOrig := len(slc) * 16 / 8
		nbBytesVarInt := len(varIntBuf)
		nbBytesEliasFano := eliasFanoBuf.Len() / 8
		fmt.Fprintf(w, "nb elements: %v\t%v\t%v\t%v\t%.3f\t\n", size, nbBytesOrig, nbBytesVarInt, nbBytesEliasFano, float64(nbBytesVarInt)/float64(nbBytesEliasFano))
	}
	w.Flush()

	w = tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "Experiment ( Mixed )\tVarInt\tElias-Fano\tImprovement\t\n")
	for _, size := range sizes {
		slc := randSliceMixed(size)
		varIntBuf := varInt64(slc)
		eliasFanoBuf := encoding.MakeEliasFanoVector(slc)
		nbBytesVarInt := len(varIntBuf)
		nbBytesEliasFano := eliasFanoBuf.Len() / 8
		fmt.Fprintf(w, "nb elements: %v\t%v\t%v\t%.3f\t\n", size, nbBytesVarInt, nbBytesEliasFano, float64(nbBytesVarInt)/float64(nbBytesEliasFano))
	}
	w.Flush()

}

func randSlice32(size int) []uint32 {
	out := make([]uint32, size)
	for i := range out {
		out[i] = rand.Uint32()
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func randSlice64(size int) []uint64 {
	out := make([]uint64, size)
	for i := range out {
		out[i] = rand.Uint64()
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func randSliceMixed(size int) []uint64 {
	out := make([]uint64, size)
	for i := range out {
		switch rand.Uint32() % 6 {
		case 0:
			out[i] = rand.Uint64()
		case 1, 3, 4, 5:
			out[i] = uint64(rand.Uint32())
		case 2:
			out[i] = uint64(uint16(rand.Uint32()))
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func randSlice16(size int) []uint16 {
	out := make([]uint16, size)
	for i := range out {
		out[i] = uint16(rand.Uint32())
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func varInt16(slice []uint16) []byte {
	out := make([]byte, len(slice)*4)
	i := 0
	for _, k := range slice {
		i += binary.PutUvarint(out[i:], uint64(k))
	}
	return out[:i]
}

func varInt32(slice []uint32) []byte {
	out := make([]byte, len(slice)*8)
	i := 0
	for _, k := range slice {
		i += binary.PutUvarint(out[i:], uint64(k))
	}
	return out[:i]
}

func varInt64(slice []uint64) []byte {
	out := make([]byte, len(slice)*12)
	i := 0
	for _, k := range slice {
		i += binary.PutUvarint(out[i:], uint64(k))
	}
	return out[:i]
}

func slice32To64(slice []uint32) []uint64 {
	out := make([]uint64, 0, len(slice))
	for _, v := range slice {
		out = append(out, uint64(v))
	}
	return out
}

func slice16To64(slice []uint16) []uint64 {
	out := make([]uint64, 0, len(slice))
	for _, v := range slice {
		out = append(out, uint64(v))
	}
	return out
}
