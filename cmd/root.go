// Package cmd constructs the root CLI and registers its subcommands.
package cmd

import (
	"example.com/learn/btrees/cmd/bplus"
	"example.com/learn/btrees/cmd/btree"

	"github.com/spf13/cobra"
)

// New creates an independent root command for one execution.
func New() *cobra.Command {
	root := &cobra.Command{
		Use:           "btrees",
		Short:         "Explore tree data structures with Go",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(btree.New(), bplus.New())
	return root
}
