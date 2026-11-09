package tree

// A leaf owns sorted stored keys and has no children.
type leaf struct{ keys []int }

func (l *leaf) nodeKeys() []int { return l.keys }

func (l *leaf) insertNonFull(key int) {
	i := 0
	for i < len(l.keys) && key > l.keys[i] {
		i++
	}
	l.keys = append(l.keys, 0)
	copy(l.keys[i+1:], l.keys[i:])
	l.keys[i] = key
}

// split keeps the left half here and moves the middle key into the parent.
func (l *leaf) split() (separator int, right node) {
	separator = l.keys[degree-1]
	right = &leaf{keys: append([]int(nil), l.keys[degree:]...)}
	l.keys = l.keys[:degree-1]
	return separator, right
}
