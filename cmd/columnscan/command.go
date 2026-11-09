// Package columnscan exposes its learning demonstration as a Cobra command.
package columnscan

import (
	"example.com/learn/btrees/cmd/columnscan/internal/demo"
	"github.com/spf13/cobra"
)

// New creates an independent demonstration command.
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "columnscan",
		Short: "Explore batch filtering and column projection",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return demo.New(cmd.OutOrStdout()).Run()
		},
	}
}
