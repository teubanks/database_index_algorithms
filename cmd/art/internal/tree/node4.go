package tree

// node4 keeps sorted byte labels beside a compact array of children.
type node4 struct {
	keys     [4]byte
	children [4]node
	count    int
}

func (n *node4) get(edge byte) node {
	for i := 0; i < n.count; i++ {
		if n.keys[i] == edge {
			return n.children[i]
		}
	}
	return nil
}

func (n *node4) put(edge byte, child node) childSet {
	i := 0
	for i < n.count && n.keys[i] < edge {
		i++
	}
	if i < n.count && n.keys[i] == edge {
		n.children[i] = child
		return n
	}
	if n.count == len(n.keys) {
		grown := &node16{}
		n.each(func(edge byte, child node) { grown.put(edge, child) })
		return grown.put(edge, child)
	}
	copy(n.keys[i+1:n.count+1], n.keys[i:n.count])
	copy(n.children[i+1:n.count+1], n.children[i:n.count])
	n.keys[i], n.children[i] = edge, child
	n.count++
	return n
}
func (n *node4) len() int     { return n.count }
func (n *node4) kind() string { return "Node4" }
func (n *node4) each(visit func(byte, node)) {
	for i := 0; i < n.count; i++ {
		visit(n.keys[i], n.children[i])
	}
}
