// Package btree exposes the btree learning subcommand.
package btree

import (
	"example.com/learn/btrees/cmd/btree/internal/demo"

	"github.com/spf13/cobra"
)

// New creates the btree demonstration command.
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "btree",
		Short: "Run the classic B-tree insertion and search demonstration",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			demo.New(cmd.OutOrStdout()).Run()
		},
	}
}
