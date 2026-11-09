// Package tree implements an in-memory B+ tree of unique integer keys.
package tree

import (
	"fmt"
	"io"
	"strings"
)

// Tree keeps all stored keys in leaves; branch keys only route searches.
// Its zero value is ready to use. It is an in-memory set without deletion.
type Tree struct{ root node }

// Insert splits overflowing nodes on the way back up, growing the root if needed.
// Duplicates are ignored, just as in the classic B-tree example.
func (t *Tree) Insert(key int) {
	if t.root == nil {
		t.root = &leaf{keys: []int{key}}
		return
	}
	separator, right := t.root.insert(key)
	if right != nil {
		t.root = &branch{
			keys:     []int{separator},
			children: []node{t.root, right},
		}
	}
}

func (t *Tree) findLeaf(key int, trace io.Writer) *leaf {
	for n := t.root; n != nil; {
		switch current := n.(type) {
		case *leaf:
			return current
		case *branch:
			i := current.childIndex(key)
			fmt.Fprintf(trace, "branch %v: follow child %d (numbered from 0)\n", current.keys, i)
			n = current.children[i]
		default:
			panic("unexpected B+ tree node type")
		}
	}
	return nil
}

// Search always descends to a leaf, even if a branch contains the same key.
func (t *Tree) Search(key int, trace io.Writer) bool {
	if trace == nil {
		trace = io.Discard
	}
	l := t.findLeaf(key, trace)
	if l == nil {
		fmt.Fprintln(trace, "empty tree; not found")
		return false
	}
	for _, candidate := range l.keys {
		if candidate == key {
			fmt.Fprintf(trace, "leaf %v: found\n", l.keys)
			return true
		}
	}
	fmt.Fprintf(trace, "leaf %v: not found\n", l.keys)
	return false
}

// Range returns keys in the inclusive interval [low, high], in sorted order.
// Descend once to the starting leaf, then follow next links instead of the tree.
func (t *Tree) Range(low, high int, trace io.Writer) []int {
	if trace == nil {
		trace = io.Discard
	}
	if low > high {
		return nil
	}
	var result []int
	for l := t.findLeaf(low, trace); l != nil; l = l.next {
		fmt.Fprintf(trace, "scan leaf %v\n", l.keys)
		for _, key := range l.keys {
			if key > high {
				return result
			}
			if key >= low {
				result = append(result, key)
			}
		}
		if l.next != nil {
			fmt.Fprintln(trace, "follow next leaf")
		}
	}
	return result
}

// Print distinguishes routing keys from stored keys and shows the leaf chain.
func (t *Tree) Print(w io.Writer) {
	if t.root == nil {
		fmt.Fprintln(w, "(empty)")
		return
	}
	var first *leaf
	var visit func(node, int)
	visit = func(n node, depth int) {
		indent := strings.Repeat("  ", depth)
		switch n := n.(type) {
		case *leaf:
			if first == nil {
				first = n
			}
			fmt.Fprintf(w, "%sleaf %v\n", indent, n.keys)
		case *branch:
			fmt.Fprintf(w, "%sbranch %v\n", indent, n.keys)
			for _, child := range n.children {
				visit(child, depth+1)
			}
		}
	}
	visit(t.root, 0)
	fmt.Fprint(w, "Leaf chain: ")
	for l := first; l != nil; l = l.next {
		fmt.Fprintf(w, "%v -> ", l.keys)
	}
	fmt.Fprintln(w, "nil")
}
