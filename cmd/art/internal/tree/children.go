package tree

// Child layouts grow at 5, 17, and 49 distinct outgoing bytes.
// Returning a childSet lets a layout replace itself when it outgrows capacity.
type childSet interface {
	get(edge byte) node
	put(edge byte, child node) childSet
	len() int
	kind() string
	each(visit func(byte, node))
}

// node48 maps each possible byte to one of 48 compact slots, plus one.
// Zero means absent; slot+1 allows byte 0 and slot 0 without sentinel collisions.
type node48 struct {
	index    [256]uint8
	children [48]node
	count    int
}

func (n *node48) get(edge byte) node {
	if slot := n.index[edge]; slot != 0 {
		return n.children[int(slot)-1]
	}
	return nil
}

func (n *node48) put(edge byte, child node) childSet {
	if slot := n.index[edge]; slot != 0 {
		n.children[int(slot)-1] = child
		return n
	}
	if n.count == len(n.children) {
		grown := &node256{}
		n.each(func(edge byte, child node) { grown.put(edge, child) })
		return grown.put(edge, child)
	}
	n.children[n.count] = child
	n.index[edge] = uint8(n.count + 1)
	n.count++
	return n
}

func (n *node48) len() int     { return n.count }
func (n *node48) kind() string { return "Node48" }
func (n *node48) each(visit func(byte, node)) {
	for edge, slot := range n.index {
		if slot != 0 {
			visit(byte(edge), n.children[int(slot)-1])
		}
	}
}

// node256 addresses children directly by their next byte.
type node256 struct {
	children [256]node
	count    int
}

func (n *node256) get(edge byte) node { return n.children[edge] }
func (n *node256) put(edge byte, child node) childSet {
	if n.children[edge] == nil {
		n.count++
	}
	n.children[edge] = child
	return n
}
func (n *node256) len() int     { return n.count }
func (n *node256) kind() string { return "Node256" }
func (n *node256) each(visit func(byte, node)) {
	for edge, child := range n.children {
		if child != nil {
			visit(byte(edge), child)
		}
	}
}
