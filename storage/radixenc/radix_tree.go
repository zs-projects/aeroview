package radixenc

type RadixTree struct {
	children map[string]*RadixTree
	isLeaf   bool
}

func MakeRadixTree(data []string) RadixTree {
	tr := MakeTRIE(data)
	rdt := RadixTree{
		children: make(map[string]*RadixTree),
	}
	for k, ssT := range tr.children {
		if ssT != nil {
			subKey, sRdt := compressSubTrie(k, ssT)
			rdt.children[subKey] = sRdt
		} else {
			rdt.children[string(k)] = nil
		}
	}
	return rdt
}

func compressSubTrie(r rune, subTR *TRIE) (string, *RadixTree) {
	runes := []rune{r}
	handle := subTR
	for len(handle.children) == 1 && !handle.isLeaf {
		curRune := handle.childrenKeys()[0]
		curSubTrie := handle.childrenValues()[0]
		runes = append(runes, curRune)
		if curSubTrie == nil {
			break
		}
		handle = curSubTrie
	}
	rdt := RadixTree{
		children: make(map[string]*RadixTree),
		isLeaf:   handle.isLeaf,
	}
	for k, ssT := range handle.children {
		if ssT != nil {
			subKey, sRdt := compressSubTrie(k, ssT)
			rdt.children[subKey] = sRdt
		} else {
			rdt.children[string(k)] = nil
		}
	}
	return string(runes), &rdt
}
