package tree

// A branch owns separator keys and exactly len(keys)+1 children.
// In a classic B-tree its keys are also stored set members.
type branch struct {
	keys     []int
	children []node
}

func (b *branch) nodeKeys() []int { return b.keys }

func (b *branch) insertNonFull(key int) {
	i := 0
	for i < len(b.keys) && key > b.keys[i] {
		i++
	}
	if len(b.children[i].nodeKeys()) == maxKeys {
		b.splitChild(i)
		if key > b.keys[i] {
			i++
		}
	}
	b.children[i].insertNonFull(key)
}

// splitChild turns [a b c] into [a] and [c], promoting b into this branch.
// Each child owns its split behavior and preserves its concrete node type.
func (b *branch) splitChild(i int) {
	middle, right := b.children[i].split()
	b.keys = append(b.keys, 0)
	copy(b.keys[i+1:], b.keys[i:])
	b.keys[i] = middle
	b.children = append(b.children, nil)
	copy(b.children[i+2:], b.children[i+1:])
	b.children[i+1] = right
}

// split moves the middle key upward and the right half into a sibling branch.
func (b *branch) split() (separator int, right node) {
	separator = b.keys[degree-1]
	right = &branch{
		keys:     append([]int(nil), b.keys[degree:]...),
		children: append([]node(nil), b.children[degree:]...),
	}
	b.keys = b.keys[:degree-1]
	// Clear moved references from the retained backing array.
	for j := degree; j < len(b.children); j++ {
		b.children[j] = nil
	}
	b.children = b.children[:degree]
	return separator, right
}
