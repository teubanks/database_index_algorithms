// Package demo shows prefix compression, uniqueness, and adaptive child growth.
package demo

import (
	"example.com/learn/btrees/cmd/art/internal/tree"
	"fmt"
	"io"
)

type Demo struct{ out io.Writer }

func New(out io.Writer) *Demo {
	if out == nil {
		out = io.Discard
	}
	return &Demo{out: out}
}

func (d *Demo) Run() error {
	var index tree.Tree
	fmt.Fprintln(d.out, "ART: compressed byte paths and adaptive child layouts.")
	for _, key := range []string{"car", "cart", "carbon", "cat", "dog"} {
		fmt.Fprintf(d.out, "\nInsert %q: %v\n", key, index.Insert(key))
		index.Print(d.out)
	}
	fmt.Fprintf(d.out, "\nDuplicate car inserted: %v\n", index.Insert("car"))
	for _, key := range []string{"cart", "car", "cab"} {
		fmt.Fprintf(d.out, "\nSearch for %q\n", key)
		index.Search(key, d.out)
	}
	fmt.Fprintln(d.out, "\nGrow one branch by inserting distinct suffix bytes:")
	var wide tree.Tree
	for i := 0; i < 49; i++ {
		wide.Insert("prefix/" + string([]byte{byte(i)}))
		if i == 3 || i == 4 || i == 16 || i == 48 {
			fmt.Fprintf(d.out, "After %d distinct suffix bytes: %s\n", i+1, wide.RootKind())
		}
	}
	return nil
}
