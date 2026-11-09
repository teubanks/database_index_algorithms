package tree

// A leaf owns stored keys and the link used by range scans.
type leaf struct {
	keys []int
	next *leaf
}

func (l *leaf) nodeKeys() []int { return l.keys }

func (l *leaf) insert(key int) (separator int, right node) {
	i := 0
	for i < len(l.keys) && l.keys[i] < key {
		i++
	}
	if i < len(l.keys) && l.keys[i] == key {
		return 0, nil
	}
	l.keys = append(l.keys, 0)
	copy(l.keys[i+1:], l.keys[i:])
	l.keys[i] = key
	if len(l.keys) <= maxKeys {
		return 0, nil
	}
	mid := len(l.keys) / 2
	sibling := &leaf{
		keys: append([]int(nil), l.keys[mid:]...),
		next: l.next,
	}
	l.keys = l.keys[:mid]
	l.next = sibling
	// COPY the right leaf's first key upward: it remains stored in the leaf.
	return sibling.keys[0], sibling
}
