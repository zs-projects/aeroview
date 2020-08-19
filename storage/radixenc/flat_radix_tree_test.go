package radixenc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/kr/pretty"
	"github.com/pierrec/lz4"
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

func BenchmarkFlatRadixTreeDecode(b *testing.B) {
	//r, err := ioutil.ReadFile("/home/ryad/listOfFiles.list")
	r, err := ioutil.ReadFile("/home/ryad/sample.txt2")
	if err != nil {
		panic("file not found")
	}
	lines := bytes.Split(r, []byte{'\n'})
	data := make([]string, 0, len(r))
	for _, l := range lines {
		data = append(data, string(l))
	}
	rdt := MakeRadixTree(data)
	frdt := MakeFlatRadixTree(rdt)
	encodedData := frdt.Encode(data)
	var buffer bytes.Buffer
	buf := make([]byte, 4)
	for _, line := range encodedData {
		for _, tok := range line {
			binary.LittleEndian.PutUint32(buf, uint32(tok))
			buffer.Write(buf)
		}
	}
	encData := buffer.Bytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err != nil {
			panic(err)
		}
		_ = (frdt.DecodeFast(encData))
	}

}

func BenchmarkLZ4Decode(b *testing.B) {
	//r, err := ioutil.ReadFile("/home/ryad/listOfFiles.list")
	r, err := ioutil.ReadFile("/home/ryad/sample.txt2")
	if err != nil {
		panic("file not found")
	}
	var buf bytes.Buffer

	wBuf := lz4.NewWriter(&buf)
	_, err = wBuf.Write(r)
	if err != nil {
		panic(err)
	}
	data := buf.Bytes()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rBuf := bytes.NewBuffer(data)
		lrBuf := lz4.NewReader(rBuf)
		u, err := ioutil.ReadAll(lrBuf)
		if err != nil {
			panic(fmt.Sprintf("Error, %v, %v, %v", err, len(u), len(r)))
		}
	}

}
