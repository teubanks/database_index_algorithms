// Package art exposes its learning demonstration as a Cobra command.
package art

import (
	"example.com/learn/btrees/cmd/art/internal/demo"
	"github.com/spf13/cobra"
)

// New creates an independent demonstration command.
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "art",
		Short: "Explore adaptive radix tree paths, lookups, and node growth",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return demo.New(cmd.OutOrStdout()).Run()
		},
	}
}
