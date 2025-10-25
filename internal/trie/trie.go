package tri

type TrieNode struct {
	Dict	      map[rune]*TrieNode
	PrefixHandles map[int64]struct{}
	WordHandles   map[int64]struct{}
	EndOfWord     bool
}
func NewTrieNode() *TrieNode {
	n := &TrieNode{
		Dict: make(map[rune]*TrieNode),
		PrefixHandles: make(map[int64]struct{}),
		WordHandles: make(map[int64]struct{}),
	}
	return n
}

type Trie struct {
	root *TrieNode
}

func (t *Trie) init() {
	if t.root == nil {
		r := NewTrieNode()
		t.root = r
	}
}


// TRIE PREFIX //
/////////////////

func (t *Trie) InsertPrefix(s string, handle int64) {
	t.init()
	curr := t.root
	for _, r := range s {
		node, ok := curr.Dict[r]
		if !ok {
			node = NewTrieNode()
			curr.Dict[r] = node
		}
		node.PrefixHandles[handle] = struct{}{}
		curr = node
	}
}

func (t *Trie) SearchPrefix(s string) (bool, []int64) {
	curr := t.root
	for _, r := range s {
		node, ok := curr.Dict[r]
		if !ok {
			return false, nil
		}
		curr = node
	}
	var handles []int64
	for k := range curr.PrefixHandles {
		handles = append(handles, k)
	}
	return true, handles
}


// TRIE WORDS //
////////////////

func (t *Trie) InsertWord(w string, handle int64) {
	t.init()
	curr := t.root
	for _, r := range w {
		node, ok := curr.Dict[r]
		if !ok {
			node = NewTrieNode()
			curr.Dict[r] = node
		}
		curr = node
	}
	curr.WordHandles[handle] = struct{}{}
	curr.EndOfWord = true
}

func (t *Trie) SearchWord(w string) (bool, []int64) {
	curr := t.root
	for _, r := range w {
		node, ok := curr.Dict[r]
		if !ok {
			return false, nil
		}
		curr = node
	}
	var handles []int64
	for k := range curr.WordHandles {
		handles = append(handles, k)
	}
	return true, handles
}



