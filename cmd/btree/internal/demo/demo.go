// Package demo runs the btree learning walkthrough.
package demo

import (
	"fmt"
	"io"

	"example.com/learn/btrees/cmd/btree/internal/tree"
)

// Demo owns the tree and output for a single walkthrough.
type Demo struct {
	tree tree.Tree
	out  io.Writer
}

// New creates a fresh walkthrough. A nil writer discards output.
func New(out io.Writer) *Demo {
	if out == nil {
		out = io.Discard
	}
	return &Demo{out: out}
}

// Run starts a fresh tree and shows each insertion followed by searches.
func (d *Demo) Run() {
	d.tree = tree.Tree{}

	fmt.Fprintln(d.out, "Classic B-tree: at most 3 keys per node. Child order is left to right.")
	for _, key := range []int{10, 20, 5, 6, 12, 30, 7, 17} {
		fmt.Fprintf(d.out, "\nInsert %d\n", key)
		d.tree.Insert(key)
		d.tree.Print(d.out)
	}
	for _, key := range []int{17, 15, 10} {
		fmt.Fprintf(d.out, "\nSearch for %d\n", key)
		d.tree.Search(key, d.out)
	}
}
