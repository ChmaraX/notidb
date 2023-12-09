package cmd

import (
	"github.com/ChmaraX/notidb/internal/tui"
	"github.com/spf13/cobra"
)

var setDefaultDbCmd = &cobra.Command{
	Use:     "set-db",
	Aliases: []string{"sd"},
	Short:   "Set default database",
	Run: func(cmd *cobra.Command, args []string) {
		tui.LoadDbs()
	},
}
