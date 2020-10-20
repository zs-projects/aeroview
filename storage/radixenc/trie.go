package radixenc

type TRIE struct {
	children map[rune]*TRIE
	isLeaf   bool
}

func MakeTRIE(data []string) TRIE {
	r := TRIE{
		children: make(map[rune]*TRIE),
	}
	for _, s := range data {
		handle := &r
		for i, b := range s {
			if child, ok := handle.children[b]; ok {
				handle = child
			} else {
				handle.children[b] = &TRIE{
					children: make(map[rune]*TRIE),
				}
				handle = handle.children[b]
			}
			if i == len(s)-1 {
				handle.isLeaf = true
			}
		}
	}
	return r
}

func (t TRIE) childrenKeys() []rune {
	r := make([]rune, 0, len(t.children))
	for k := range t.children {
		r = append(r, k)
	}
	return r
}

func (t TRIE) childrenValues() []*TRIE {
	r := make([]*TRIE, 0, len(t.children))
	for _, k := range t.children {
		r = append(r, k)
	}
	return r
}
