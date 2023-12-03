package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	usage   = `notion-cli`
	example = `  notion-cli ld
  notion-cli a --title "Meeting Notes" --content "Notes from the team meeting on 3rd Dec"
  notion-cli sd --database 1234567890abcdef
  notion-cli le --database 1234567890abcdef`
)

var rootCmd = &cobra.Command{
	Use:           usage,
	Example:       example,
	Short:         "notion-cli is a CLI tool for interacting with Notion databases",
	Version:       "0.0.1",
	SilenceErrors: true,
}

func init() {
	listEntriesCmd.Flags().String("database", "", "ID of the database to list entries from")
	setDefaultCmd.Flags().String("database", "", "ID of the database to set as default")

	addEntryCmd.Flags().String("title", "", "Title of the new entry")
	addEntryCmd.Flags().String("content", "", "Content of the new entry")
	addEntryCmd.Flags().String("database", "", "ID of the database to add entry to")

	rootCmd.AddCommand(listDatabasesCmd)
	rootCmd.AddCommand(addEntryCmd)
	rootCmd.AddCommand(setDefaultCmd)
	rootCmd.AddCommand(listEntriesCmd)

}

var listDatabasesCmd = &cobra.Command{
	Use:     "list-databases",
	Aliases: []string{"ld"},
	Short:   "Lists all available databases",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list-databases called")
	},
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
