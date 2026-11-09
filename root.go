package main

import "github.com/spf13/cobra"

// Register additional learning examples here as subcommands.
func newRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:           "btrees",
		Short:         "Explore tree data structures with Go",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(newBTreeCommand(), newBPlusCommand())
	return root
}
