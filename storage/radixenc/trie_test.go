package radixenc

import (
	"bytes"
	"fmt"
	"testing"
)

func TestMakeTrie(t *testing.T) {
	data := []string{"string1", "string2", "dance", "opera", "winter", "alstom", "netscape", "hi", "lower", "high"}
	tr := MakeTRIE(data)
	if r := len(tr.children); r != 8 {
		t.Errorf("Expected the root node to have %v children, got %v", 8, r)
	}
	for _, str := range data {
		handle := &tr
		for _, r := range str {
			if _, ok := handle.children[r]; !ok {
				var b bytes.Buffer
				fmt.Fprintf(&b, "nodes:")
				for k := range handle.children {
					fmt.Fprintf(&b, " %s", string(k))
				}
				fmt.Fprintf(&b, "\n")
				t.Errorf("Looking at %v in %v \n", string(r), b.String())
			} else {
				handle = handle.children[r]
			}
		}
	}
}
