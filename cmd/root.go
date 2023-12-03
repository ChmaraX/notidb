package cmd

import (
	"fmt"
	"os"

	"github.com/ChmaraX/notidb/internal"
	"github.com/spf13/cobra"
)

const (
	usage   = `notidb`
	example = `  notidb ld
  notidb a --title "Meeting Notes" --content "Notes from the team meeting on 3rd Dec"
  notidb sd --database 1234567890abcdef
  notidb le --database 1234567890abcdef`
)

var rootCmd = &cobra.Command{
	Use:           usage,
	Example:       example,
	Short:         "notidb is a CLI tool for interacting with Notion databases",
	Version:       "0.0.1",
	SilenceErrors: true,
}

func init() {
	listEntriesCmd.Flags().String("database", "", "ID of the database to list entries from")
	setDefaultCmd.Flags().String("database", "", "ID of the database to set as default")

	addEntryCmd.Flags().String("title", "", "Title of the new entry")
	addEntryCmd.Flags().String("content", "", "Content of the new entry")
	addEntryCmd.Flags().String("database", "", "ID of the database to add entry to")

	rootCmd.AddCommand(internal.ListDatabasesCmd())
	rootCmd.AddCommand(addEntryCmd)
	rootCmd.AddCommand(setDefaultCmd)
	rootCmd.AddCommand(listEntriesCmd)

}

var addEntryCmd = &cobra.Command{
	Use:     "add-entry",
	Aliases: []string{"a"},
	Short:   "Adds a new entry to a database",
	Run: func(cmd *cobra.Command, args []string) {
		// Implement add entry functionality
	},
}

var setDefaultCmd = &cobra.Command{
	Use:   "set-default",
	Short: "Sets a default database",
	Run: func(cmd *cobra.Command, args []string) {
		// Implement set default database functionality
	},
}

var listEntriesCmd = &cobra.Command{
	Use:   "list-entries",
	Short: "Lists top 10 entries from a database",
	Run: func(cmd *cobra.Command, args []string) {
		// Implement list entries functionality
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
