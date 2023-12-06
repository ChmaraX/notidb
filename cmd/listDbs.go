package cmd

import (
	"fmt"
	"log"

	"github.com/ChmaraX/notidb/internal"
	"github.com/ChmaraX/notidb/internal/settings"
	"github.com/ChmaraX/notidb/internal/tui"
	"github.com/spf13/cobra"
)

func defaultDbExists(dbs []internal.NotionDb, defaultDbId string) bool {
	for _, db := range dbs {
		if db.Id == defaultDbId {
			return true
		}
	}
	return false
}

var listDbsCmd = &cobra.Command{
	Use:     "list-dbs",
	Aliases: []string{"ld", "ls"},
	Short:   "Lists all available databases",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list-databases called")

		dbs := internal.GetAllNotionDbs()
		defaultDbId, err := settings.GetDefaultDatabase()

		if err != nil {
			log.Fatalf("Error getting default database: %v", err)
		}

		if len(dbs) == 0 {
			log.Fatalf("No databases found in your workspace or the access is not granted.")
		}

		// check if defaultDbId exists in dbs
		if !defaultDbExists(dbs, defaultDbId) && defaultDbId != settings.NoDefaultDatabaseId {
			log.Fatalf("Database which is set as default (%s) was not found in your workspace or the access is not granted.", defaultDbId)
		}

		tui.GetDbsListTUI(dbs, defaultDbId)
	},
}
