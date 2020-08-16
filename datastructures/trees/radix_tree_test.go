package trees

import (
	"testing"

	"github.com/kr/pretty"
)

func TestMakeRadixTree(t *testing.T) {
	data := []string{"string1", "string2", "dance", "opera", "ope2", "winter", "alstom", "netscape", "hi", "lower", "high"}
	root := []string{"string", "dance", "ope", "winter", "alstom", "netscape", "hi", "lower"}
	expectedStructure := RadixTree{
		children: map[string]*RadixTree{
			"string": {
				children: map[string]*RadixTree{
					"1": {isLeaf: true},
					"2": {isLeaf: true},
				},
			},
			"ope": {
				children: map[string]*RadixTree{
					"ra": {isLeaf: true},
					"2":  {isLeaf: true},
				},
			},
			"hi": {
				children: map[string]*RadixTree{
					"gh": {isLeaf: true},
				},
				isLeaf: true,
			},
			"alstom":   {isLeaf: true},
			"lower":    {isLeaf: true},
			"netscape": {isLeaf: true},
			"winter":   {isLeaf: true},
			"dance":    {isLeaf: true},
		},
	}
	rdt := MakeRadixTree(data)
	for _, v := range root {
		if _, ok := rdt.children[v]; !ok {
			t.Errorf("Expected to find %v as prefix", v)
		}
	}
	diff := pretty.Diff(rdt, expectedStructure)
	if len(diff) != 0 || diff != nil {
		t.Errorf("Radix Tree not as expected %# v", diff)
	}
}

func TestCompressSubTrie(t *testing.T) {

}
