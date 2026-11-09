package main

import (
	"fmt"
	"io"
	"strings"
)

// Small capacities make splits visible. This example stores unique integer keys.
const bplusMaxKeys = 3

type bplusNode interface {
	bplusKeys() []int
}

type bplusLeaf struct {
	keys []int
	next *bplusLeaf
}

func (l *bplusLeaf) bplusKeys() []int { return l.keys }

type bplusBranch struct {
	keys     []int       // Each key is the minimum key of the child to its right.
	children []bplusNode // Always len(keys)+1.
}

func (b *bplusBranch) bplusKeys() []int { return b.keys }

// BPlusTree keeps all stored keys in leaves; branch keys only route searches.
// Its zero value is ready to use. It is an in-memory set without deletion.
type BPlusTree struct{ root bplusNode }

// Insert splits overflowing nodes on the way back up, growing the root if needed.
// Duplicates are ignored, just as in the classic B-tree example.
func (t *BPlusTree) Insert(key int) {
	if t.root == nil {
		t.root = &bplusLeaf{keys: []int{key}}
		return
	}
	separator, right := bplusInsert(t.root, key)
	if right != nil {
		t.root = &bplusBranch{
			keys:     []int{separator},
			children: []bplusNode{t.root, right},
		}
	}
}

// The existing node retains its left half. A non-nil right signals a split.
func bplusInsert(n bplusNode, key int) (separator int, right bplusNode) {
	switch n := n.(type) {
	case *bplusLeaf:
		i := 0
		for i < len(n.keys) && n.keys[i] < key {
			i++
		}
		if i < len(n.keys) && n.keys[i] == key {
			return 0, nil
		}
		n.keys = append(n.keys, 0)
		copy(n.keys[i+1:], n.keys[i:])
		n.keys[i] = key
		if len(n.keys) <= bplusMaxKeys {
			return 0, nil
		}
		mid := len(n.keys) / 2
		sibling := &bplusLeaf{
			keys: append([]int(nil), n.keys[mid:]...),
			next: n.next,
		}
		n.keys = n.keys[:mid]
		n.next = sibling
		// COPY the right leaf's first key upward: it remains stored in the leaf.
		return sibling.keys[0], sibling
	case *bplusBranch:
		i := bplusChildIndex(n.keys, key)
		separator, right := bplusInsert(n.children[i], key)
		if right == nil {
			return 0, nil
		}
		n.keys = append(n.keys, 0)
		copy(n.keys[i+1:], n.keys[i:])
		n.keys[i] = separator
		n.children = append(n.children, nil)
		copy(n.children[i+2:], n.children[i+1:])
		n.children[i+1] = right
		if len(n.keys) <= bplusMaxKeys {
			return 0, nil
		}
		mid := len(n.keys) / 2
		separator = n.keys[mid]
		sibling := &bplusBranch{
			keys:     append([]int(nil), n.keys[mid+1:]...),
			children: append([]bplusNode(nil), n.children[mid+1:]...),
		}
		n.keys = n.keys[:mid]
		for j := mid + 1; j < len(n.children); j++ {
			n.children[j] = nil
		}
		n.children = n.children[:mid+1]
		// MOVE the middle separator up; the underlying key still exists in a leaf.
		return separator, sibling
	default:
		panic("unexpected B+ tree node type")
	}
}

// Equality goes RIGHT at a separator: the matching stored key is in that subtree.
func bplusChildIndex(keys []int, key int) int {
	i := 0
	for i < len(keys) && key >= keys[i] {
		i++
	}
	return i
}

func (t *BPlusTree) findLeaf(key int, trace io.Writer) *bplusLeaf {
	for n := t.root; n != nil; {
		switch current := n.(type) {
		case *bplusLeaf:
			return current
		case *bplusBranch:
			i := bplusChildIndex(current.keys, key)
			fmt.Fprintf(trace, "branch %v: follow child %d (numbered from 0)\n", current.keys, i)
			n = current.children[i]
		default:
			panic("unexpected B+ tree node type")
		}
	}
	return nil
}

// Search always descends to a leaf, even if a branch contains the same key.
func (t *BPlusTree) Search(key int, trace io.Writer) bool {
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
func (t *BPlusTree) Range(low, high int, trace io.Writer) []int {
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
func (t *BPlusTree) Print(w io.Writer) {
	if t.root == nil {
		fmt.Fprintln(w, "(empty)")
		return
	}
	var first *bplusLeaf
	var visit func(bplusNode, int)
	visit = func(n bplusNode, depth int) {
		indent := strings.Repeat("  ", depth)
		switch n := n.(type) {
		case *bplusLeaf:
			if first == nil {
				first = n
			}
			fmt.Fprintf(w, "%sleaf %v\n", indent, n.keys)
		case *bplusBranch:
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
