package tree

// node16 keeps sorted byte labels beside a compact array of children.
type node16 struct {
	keys     [16]byte
	children [16]node
	count    int
}

func (n *node16) get(edge byte) node {
	for i := 0; i < n.count; i++ {
		if n.keys[i] == edge {
			return n.children[i]
		}
	}
	return nil
}

func (n *node16) put(edge byte, child node) childSet {
	i := 0
	for i < n.count && n.keys[i] < edge {
		i++
	}
	if i < n.count && n.keys[i] == edge {
		n.children[i] = child
		return n
	}
	if n.count == len(n.keys) {
		grown := &node48{}
		n.each(func(edge byte, child node) { grown.put(edge, child) })
		return grown.put(edge, child)
	}
	copy(n.keys[i+1:n.count+1], n.keys[i:n.count])
	copy(n.children[i+1:n.count+1], n.children[i:n.count])
	n.keys[i], n.children[i] = edge, child
	n.count++
	return n
}
func (n *node16) len() int     { return n.count }
func (n *node16) kind() string { return "Node16" }
func (n *node16) each(visit func(byte, node)) {
	for i := 0; i < n.count; i++ {
		visit(n.keys[i], n.children[i])
	}
}
