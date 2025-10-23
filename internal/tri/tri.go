package tri

type insertType int
const (
	prefixType insertType = iota
	wordType
)

type TriNode struct {
	Dict      map[rune]*TriNode
	EndOfWord bool
	Handles   map[int64]insertType
}

func NewTriNode() *TriNode {
	n := &TriNode{
		Dict: make(map[rune]*TriNode),
		Handles: make(map[int64]insertType),
	}
	return n
}

func (n *TriNode) addHandle(h int64, t insertType) {
	n.Handles[h] = t
}

type Tri struct {
	root *TriNode
}

func (t *Tri) init() {
	if t.root == nil {
		r := NewTriNode()
		t.root = r
	}
}

func (t *Tri) InsertWord(s string, handle int64) {
	t.init()
	curr := t.root
	for _, r := range s {
		node, ok := curr.Dict[r]
		if !ok {
			node = NewTriNode()
			curr.Dict[r] = node
		}
		curr = node
	}
	curr.addHandle(handle, wordType)
	curr.EndOfWord = true
}

func (t *Tri) InsertPrefix(s string, handle int64) {
	t.init()
	curr := t.root
	for _, r := range s {
		node, ok := curr.Dict[r]
		if !ok {
			node = NewTriNode()
			curr.Dict[r] = node
		}
		node.addHandle(handle, prefixType)
		curr = node
	}
}

func (t *Tri) SearchPrefix(s string) (bool, []int64) {
	curr := t.root
	for _, r := range s {
		node, ok := curr.Dict[r]
		if !ok {
			return false, nil
		}
		curr = node
	}
	var handles []int64
	for k := range curr.Handles {
		handles = append(handles, k)
	}
	return true, handles
}




