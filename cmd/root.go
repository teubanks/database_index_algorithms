// Package cmd constructs the root CLI and registers its subcommands.
package cmd

import (
	"example.com/learn/btrees/cmd/art"
	"example.com/learn/btrees/cmd/bplus"
	"example.com/learn/btrees/cmd/btree"
	"example.com/learn/btrees/cmd/columnscan"
	"example.com/learn/btrees/cmd/zonemap"

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
	root.AddCommand(btree.New(), bplus.New(), art.New(), zonemap.New(), columnscan.New())
	return root
}
