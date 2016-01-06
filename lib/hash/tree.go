package hash

// A Tree is a binary tree with integer values.
type Tree struct {
	Left  *Tree
	Value int64
	Right *Tree
}

func Add(t *Tree, v int64) *Tree {
	if t == nil {
		return &Tree{nil, v, nil}
	}
	if v < t.Value {
		t.Left = Add(t.Left, v)
		return t
	}
	t.Right = Add(t.Right, v)
	return t
}

func FindNext(t *Tree, v int64, p *Tree) {
	if t == nil {
		return
	}
	FindNext(t.Left, v, p)
	if t.Value >= v {
		p = t
		return
	}
	FindNext(t.Right, v, p)
}
