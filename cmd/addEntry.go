package cmd

import (
	"github.com/spf13/cobra"
)

var addEntryCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "Adds a new entry to a database",
	Run: func(cmd *cobra.Command, args []string) {
		// Implement add entry functionality
	},
}
