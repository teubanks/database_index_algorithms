package main

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func newBPlusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "bplus",
		Short: "Run the B+ tree insertion, search, and range scan demonstration",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			runBPlusDemo(cmd.OutOrStdout())
		},
	}
}

func runBPlusDemo(out io.Writer) {
	var tree BPlusTree
	fmt.Fprintln(out, "B+ tree: stored keys live in leaves; branch keys route searches.")
	fmt.Fprintln(out, "At most 3 keys per node. Leaf splits copy a separator into the parent.")
	for _, key := range []int{10, 20, 5, 6, 12, 30, 7, 17} {
		fmt.Fprintf(out, "\nInsert %d\n", key)
		tree.Insert(key)
		tree.Print(out)
	}
	for _, key := range []int{17, 15, 10} {
		fmt.Fprintf(out, "\nSearch for %d\n", key)
		tree.Search(key, out)
	}
	fmt.Fprintln(out, "\nRange [6, 20] (inclusive)")
	keys := tree.Range(6, 20, out)
	fmt.Fprintf(out, "Result: %v\n", keys)
}
