package datastructures

import (
	"testing"

	"github.com/kr/pretty"
)

func TestMakeRadixTree(t *testing.T) {
	data := []string{"string1", "string2", "dance", "opera", "ope2", "winter", "alstom", "netscape", "hi", "lower", "high"}
	rdt := MakeRadixTree(data)
	if _, ok := rdt.children["string"]; !ok {
		t.Errorf("Expected the root node to contain 'string'")
	}

	if _, ok := rdt.children["hi"]; !ok {
		t.Errorf("Expected the root node to contain 'h'")
	}
	t.Errorf("%# v", pretty.Formatter(rdt))
}

func TestCompressSubTrie(t *testing.T) {

}
