package cmd

import (
	"fmt"

	"github.com/ChmaraX/notidb/internal/notion"
	"github.com/ChmaraX/notidb/internal/settings"
	"github.com/ChmaraX/notidb/internal/tui"
	"github.com/jomei/notionapi"
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

func dbExists(dbs []notionapi.Database, dbId string) bool {
	for _, db := range dbs {
		if string(db.ID) == dbId {
			return true
		}
	}
	return false
}

var setDefaultDbCmd = &cobra.Command{
	Use:     "set-db",
	Aliases: []string{"sd"},
	Short:   "Set default database",
	Run: func(cmd *cobra.Command, args []string) {
		m := tui.NewLoadingModel("Calling Notion API - loading databases", loadDatabases, loadDefaultDatabase)
		dbs := m.GetResponse("dbs").Data.([]notionapi.Database)
		defaultDbId := m.GetResponse("defaultDb").Data.(string)

		if !dbExists(dbs, defaultDbId) && defaultDbId != settings.NoDefaultDatabaseId {
			fmt.Printf("database which is set as default (%s) was not found in your workspace or the access is not granted", defaultDbId)
			return
		}

		tui.InitDbListModel(dbs, defaultDbId)
	},
}
