// Package tree implements a byte-oriented adaptive radix tree of unique strings.
package tree

import (
	"fmt"
	"io"
)

// Tree uses compressed paths and adaptive child layouts. Its zero value is ready.
// Strings are treated as bytes, including empty strings and embedded zero bytes.
type Tree struct {
	root node
	size int
}

func (t *Tree) Len() int { return t.size }

// RootKind exposes the current root layout for the adaptive-growth walkthrough.
func (t *Tree) RootKind() string {
	switch root := t.root.(type) {
	case nil:
		return "empty"
	case *leaf:
		return "leaf"
	case *branch:
		return root.children.kind()
	default:
		panic("unexpected ART node type")
	}
}

// Insert returns false for a duplicate, demonstrating the basis of uniqueness checks.
func (t *Tree) Insert(key string) bool {
	if t.Search(key, nil) {
		return false
	}
	if t.root == nil {
		t.root = &leaf{key: key}
	} else {
		t.root = t.root.insert(key, 0)
	}
	t.size++
	return true
}

func (t *Tree) Search(key string, trace io.Writer) bool {
	if trace == nil {
		trace = io.Discard
	}
	if t.root == nil {
		fmt.Fprintln(trace, "empty tree; not found")
		return false
	}
	return t.root.lookup(key, 0, trace)
}

func (t *Tree) Print(out io.Writer) {
	if t.root == nil {
		fmt.Fprintln(out, "(empty)")
		return
	}
	t.root.print(out, "")
}

type node interface {
	insert(key string, depth int) node
	lookup(key string, depth int, trace io.Writer) bool
	print(out io.Writer, indent string)
}
