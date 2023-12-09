package cmd

import (
	"github.com/ChmaraX/notidb/internal/tui"
	"github.com/spf13/cobra"
)

var listDbsCmd = &cobra.Command{
	Use:     "list-dbs",
	Aliases: []string{"ld", "ls"},
	Short:   "Lists all available databases",
	Run: func(cmd *cobra.Command, args []string) {
		tui.LoadDbs()
	},
}
