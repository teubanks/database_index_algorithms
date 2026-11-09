package main

import (
	"fmt"
	"io"
	"strings"
)

// This is a classic B-tree: keys can live in leaves and branches.
// Degree 2 means each node holds at most 3 keys and branches have at most 4 children.
const degree = 2
const maxKeys = 2*degree - 1

// A child (and the root) can be either a leaf or a branch.
type node interface {
	nodeKeys() []int
}

type leaf struct {
	keys []int
}

func (l *leaf) nodeKeys() []int { return l.keys }

type branch struct {
	keys     []int
	children []node // Always len(keys)+1; leaves have no child field.
}

func (b *branch) nodeKeys() []int { return b.keys }

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
		splitChild(root, 0)
		t.root = root
	}
	insertNonFull(t.root, key)
}

// splitChild turns [a b c] into [a] and [c], promoting b to parent.
// The two halves keep their original type; branch children move with their keys.
func splitChild(parent *branch, i int) {
	left := parent.children[i]
	keys := left.nodeKeys()
	middle := keys[degree-1]
	rightKeys := append([]int(nil), keys[degree:]...)
	var right node
	switch left := left.(type) {
	case *leaf:
		right = &leaf{keys: rightKeys}
		left.keys = left.keys[:degree-1]
	case *branch:
		right = &branch{
			keys:     rightKeys,
			children: append([]node(nil), left.children[degree:]...),
		}
		left.keys = left.keys[:degree-1]
		// Clear moved references from the retained backing array.
		for j := degree; j < len(left.children); j++ {
			left.children[j] = nil
		}
		left.children = left.children[:degree]
	}

	parent.keys = append(parent.keys, 0)
	copy(parent.keys[i+1:], parent.keys[i:])
	parent.keys[i] = middle
	parent.children = append(parent.children, nil)
	copy(parent.children[i+2:], parent.children[i+1:])
	parent.children[i+1] = right
}

func insertNonFull(n node, key int) {
	keys := n.nodeKeys()
	i := 0
	for i < len(keys) && key > keys[i] {
		i++
	}
	switch n := n.(type) {
	case *leaf:
		n.keys = append(n.keys, 0)
		copy(n.keys[i+1:], n.keys[i:])
		n.keys[i] = key
	case *branch:
		if len(n.children[i].nodeKeys()) == maxKeys {
			splitChild(n, i)
			if key > n.keys[i] {
				i++
			}
		}
		insertNonFull(n.children[i], key)
	}
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
