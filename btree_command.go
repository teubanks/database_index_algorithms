package main

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func newBTreeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "btree",
		Short: "Run the classic B-tree insertion and search demonstration",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			runBTreeDemo(cmd.OutOrStdout())
		},
	}
}

func runBTreeDemo(out io.Writer) {
	var tree Tree
	fmt.Fprintln(out, "Classic B-tree: at most 3 keys per node. Child order is left to right.")
	for _, key := range []int{10, 20, 5, 6, 12, 30, 7, 17} {
		fmt.Fprintf(out, "\nInsert %d\n", key)
		tree.Insert(key)
		tree.Print(out)
	}
	for _, key := range []int{17, 15, 10} {
		fmt.Fprintf(out, "\nSearch for %d\n", key)
		tree.Search(key, out)
	}
}
