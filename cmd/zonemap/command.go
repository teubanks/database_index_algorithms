// Package zonemap exposes its learning demonstration as a Cobra command.
package zonemap

import (
	"example.com/learn/btrees/cmd/zonemap/internal/demo"
	"github.com/spf13/cobra"
)

// New creates an independent demonstration command.
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "zonemap",
		Short: "Compare min/max segment pruning on ordered and mixed data",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return demo.New(cmd.OutOrStdout()).Run()
		},
	}
}
