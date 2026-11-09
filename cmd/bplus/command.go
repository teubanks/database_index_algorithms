// Package bplus exposes the bplus learning subcommand.
package bplus

import (
	"example.com/learn/btrees/cmd/bplus/internal/demo"

	"github.com/spf13/cobra"
)

// New creates the bplus demonstration command.
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "bplus",
		Short: "Run the B+ tree insertion, search, and range scan demonstration",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			demo.New(cmd.OutOrStdout()).Run()
		},
	}
}
