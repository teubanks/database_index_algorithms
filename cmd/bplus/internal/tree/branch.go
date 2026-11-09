package tree

// Each separator is the minimum key in the child to its right.
type branch struct {
	keys     []int
	children []node
}

func (b *branch) nodeKeys() []int { return b.keys }

func (b *branch) insert(key int) (separator int, right node) {
	i := b.childIndex(key)
	separator, right = b.children[i].insert(key)
	if right == nil {
		return 0, nil
	}
	b.keys = append(b.keys, 0)
	copy(b.keys[i+1:], b.keys[i:])
	b.keys[i] = separator
	b.children = append(b.children, nil)
	copy(b.children[i+2:], b.children[i+1:])
	b.children[i+1] = right
	if len(b.keys) <= maxKeys {
		return 0, nil
	}
	mid := len(b.keys) / 2
	separator = b.keys[mid]
	sibling := &branch{
		keys:     append([]int(nil), b.keys[mid+1:]...),
		children: append([]node(nil), b.children[mid+1:]...),
	}
	b.keys = b.keys[:mid]
	for j := mid + 1; j < len(b.children); j++ {
		b.children[j] = nil
	}
	b.children = b.children[:mid+1]
	// MOVE the middle separator up; the underlying key still exists in a leaf.
	return separator, sibling
}

// Equality goes RIGHT at a separator: the matching stored key is in that subtree.
func (b *branch) childIndex(key int) int {
	i := 0
	for i < len(b.keys) && key >= b.keys[i] {
		i++
	}
	return i
}
