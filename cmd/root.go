package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	usage   = `notidb`
	example = `  notidb sd
  notidb a --title "Meeting Notes" --body "Notes from the team meeting on 3rd Dec"
  notidb le --database 1234567890abcdef`
)

var rootCmd = &cobra.Command{
	Use:           usage,
	Example:       example,
	Short:         "notidb is a CLI tool for quick interaction with Notion databases",
	Version:       "0.0.1",
	SilenceErrors: true,
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
