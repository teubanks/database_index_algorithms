package tree

import (
	"fmt"
	"io"
	"strings"
)

// A branch stores a compressed path and adapts its child layout to fanout.
// terminal distinguishes an exact prefix key from longer keys without reserving a byte.
type branch struct {
	prefix   string
	terminal *leaf
	children childSet
}

func (b *branch) attach(key string, depth int, child node) {
	if depth == len(key) {
		b.terminal = child.(*leaf)
		return
	}
	b.children = b.children.put(key[depth], child)
}

func (b *branch) insert(key string, depth int) node {
	shared := 0
	for shared < len(b.prefix) && depth+shared < len(key) && b.prefix[shared] == key[depth+shared] {
		shared++
	}
	if shared < len(b.prefix) {
		// Split a compressed path at its first mismatching byte.
		parent := &branch{prefix: b.prefix[:shared], children: &node4{}}
		oldEdge := b.prefix[shared]
		b.prefix = b.prefix[shared+1:]
		parent.children = parent.children.put(oldEdge, b)
		parent.attach(key, depth+shared, &leaf{key: key})
		return parent
	}
	depth += len(b.prefix)
	if depth == len(key) {
		b.terminal = &leaf{key: key}
		return b
	}
	edge := key[depth]
	child := b.children.get(edge)
	if child == nil {
		child = &leaf{key: key}
	} else {
		child = child.insert(key, depth+1)
	}
	b.children = b.children.put(edge, child)
	return b
}

func (b *branch) lookup(key string, depth int, trace io.Writer) bool {
	fmt.Fprintf(trace, "%s prefix=%q: ", b.children.kind(), b.prefix)
	if !strings.HasPrefix(key[depth:], b.prefix) {
		fmt.Fprintln(trace, "prefix mismatch; not found")
		return false
	}
	depth += len(b.prefix)
	if depth == len(key) {
		fmt.Fprintln(trace, "check end-of-key")
		if b.terminal == nil {
			return false
		}
		return b.terminal.lookup(key, depth, trace)
	}
	edge := key[depth]
	fmt.Fprintf(trace, "follow byte 0x%02x\n", edge)
	child := b.children.get(edge)
	if child == nil {
		fmt.Fprintln(trace, "missing edge; not found")
		return false
	}
	return child.lookup(key, depth+1, trace)
}

func (b *branch) print(out io.Writer, indent string) {
	fmt.Fprintf(out, "%s%s prefix=%q children=%d\n", indent, b.children.kind(), b.prefix, b.children.len())
	if b.terminal != nil {
		fmt.Fprintf(out, "%s  <end> -> %q\n", indent, b.terminal.key)
	}
	b.children.each(func(edge byte, child node) {
		fmt.Fprintf(out, "%s  byte 0x%02x:\n", indent, edge)
		child.print(out, indent+"    ")
	})
}
