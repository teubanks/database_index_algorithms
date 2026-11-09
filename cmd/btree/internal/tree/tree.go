// Package tree implements the classic in-memory integer B-tree.
package tree

import (
	"fmt"
	"io"
	"strings"
)

// Tree is a set of unique integer keys. Its zero value is ready to use.
type Tree struct{ root node }

// Search follows one child per level. The trace shows each decision.
func (t *Tree) Search(key int, trace io.Writer) bool {
	if trace == nil {
		trace = io.Discard
	}
	for n := t.root; n != nil; {
		keys := n.nodeKeys()
		i := 0
		for i < len(keys) && key > keys[i] {
			i++
		}
		fmt.Fprintf(trace, "visit %v: ", keys)
		if i < len(keys) && key == keys[i] {
			fmt.Fprintln(trace, "found")
			return true
		}
		b, ok := n.(*branch)
		if !ok {
			fmt.Fprintln(trace, "leaf reached; not found")
			return false
		}
		fmt.Fprintf(trace, "follow child %d (numbered from 0)\n", i)
		n = b.children[i]
	}
	return false
}

// Insert treats the tree as a set: duplicate keys are ignored.
// Split full nodes before descending, so insertion always has room.
func (t *Tree) Insert(key int) {
	if t.Search(key, nil) {
		return
	}
	if t.root == nil {
		t.root = &leaf{keys: []int{key}}
		return
	}
	if len(t.root.nodeKeys()) == maxKeys {
		root := &branch{children: []node{t.root}}
		root.splitChild(0)
		t.root = root
	}
	t.root.insertNonFull(key)
}

// Print shows the root first, with its children indented underneath.
func (t *Tree) Print(w io.Writer) {
	var visit func(node, int)
	visit = func(n node, depth int) {
		if n == nil {
			fmt.Fprintln(w, "(empty)")
			return
		}
		fmt.Fprintf(w, "%s%v\n", strings.Repeat("  ", depth), n.nodeKeys())
		if b, ok := n.(*branch); ok {
			for _, child := range b.children {
				visit(child, depth+1)
			}
		}
	}
	visit(t.root, 0)
}
