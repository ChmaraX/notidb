package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	usage   = `notidb`
	example = `  notidb sd
  notidb add
  notidb add --title "Book Idea" --content "A book about the history of the internet"
  notidb a -t "Book Idea" -c "A book about the history of the internet"
  notidb "Book Idea" "A book about the history of the internet"`
)

var rootCmd = &cobra.Command{
	Use:           usage,
	Example:       example,
	Short:         "notidb is a CLI tool for quick interaction with Notion databases",
	Version:       "0.0.1",
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Use != "init" {
			initNotionClient()
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(setDefaultDbCmd)
	rootCmd.AddCommand(addEntryCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
