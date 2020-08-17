package radixenc

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/kr/pretty"
	"github.com/zs-projects/aeroview/analysis/randutils"
)

func TestMakeFlatRadixTree(t *testing.T) {
	data := []string{"string1", "string2", "dance", "opera", "ope2", "winter", "alstom", "netscape", "hi", "lower", "high"}
	rdt := MakeRadixTree(data)
	frdt := MakeFlatRadixTree(rdt)
	if nb, pos := frdt.Children(0); nb == 0 || pos != 1 {
		t.Errorf("Root node is not good")
	}
	fmt.Printf("Flat Radix Tree: %# v\n", pretty.Formatter(frdt))
	fmt.Printf("Flat Radix Tree: Overhead %.3f\t size %v\n", frdt.Overhead(), frdt.Size())
	encodedData := frdt.Encode(data)
	decodedData := (frdt.Decode(encodedData))
	if !reflect.DeepEqual(decodedData, data) {
		t.Errorf("Decoding data failed.")
	}
	out := make([]uint64, 0, len(encodedData))
	for _, v := range encodedData {
		for _, nb := range v {
			out = append(out, uint64(nb))
		}
	}
	compressed := randutils.VarInt64(out)
	originalSize := 0
	for _, s := range data {
		originalSize += len(s)
	}
	fmt.Printf("Flat Radix Tree: Original Size %v\t Encoded Size %v\n", originalSize*4, len(compressed))
}
