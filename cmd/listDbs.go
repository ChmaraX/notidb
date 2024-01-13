package cmd

import (
	"fmt"

	"github.com/ChmaraX/notidb/internal/notion"
	"github.com/ChmaraX/notidb/internal/settings"
	"github.com/ChmaraX/notidb/internal/tui"
	"github.com/spf13/cobra"
)

func loadDatabases() tui.Response {
	databases, err := notion.GetAllNotionDbs()
	id := "dbs"

	if err != nil {
		return tui.Response{Id: id, Data: nil, Err: fmt.Errorf("error loading databases: %v", err)}
	}
	if len(databases) == 0 {
		return tui.Response{Id: id, Data: nil, Err: fmt.Errorf("no databases found in your workspace or the access is not granted")}
	}
	return tui.Response{Id: id, Data: databases, Err: nil}
}

func loadDefaultDatabase() tui.Response {
	defaultDbId, err := settings.GetDefaultDatabase()
	id := "defaultDb"

	if err != nil {
		return tui.Response{Id: id, Data: nil, Err: fmt.Errorf("error loading default database: %v", err)}
	}
	return tui.Response{Id: id, Data: defaultDbId, Err: nil}
}

var setDefaultDbCmd = &cobra.Command{
	Use:     "set-db",
	Aliases: []string{"sd"},
	Short:   "Set default database",
	Run: func(cmd *cobra.Command, args []string) {
		// tui.list()

		m := tui.NewLoadingModel(loadDatabases, loadDefaultDatabase)
		// dbs := m.GetResponse("dbs")
		defaultDb := m.GetResponse("defaultDb")

		fmt.Printf("model: %v\n", defaultDb)

	},
}
