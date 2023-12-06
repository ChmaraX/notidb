package cmd

import (
	"context"
	"fmt"

	"github.com/ChmaraX/notidb/internal"
	"github.com/ChmaraX/notidb/internal/tui"
	"github.com/dstotijn/go-notion"
	"github.com/spf13/cobra"
)

var listDbsCmd = &cobra.Command{
	Use:     "list-dbs",
	Aliases: []string{"ld", "ls"},
	Short:   "Lists all available databases",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list-databases called")

		res, err := internal.NotionClient.Search(
			context.TODO(),
			&notion.SearchOpts{
				Filter: &notion.SearchFilter{
					Value:    "database",
					Property: "object",
				},
			},
		)
		if err != nil {
			panic(err)
		}

		dbs := make([]internal.NotionDb, len(res.Results))
		for i, db := range res.Results {
			dbs[i].Id = db.(notion.Database).ID
			dbs[i].Title = db.(notion.Database).Title[0].PlainText // TODO check title length
			if len(db.(notion.Database).Description) > 0 {
				dbs[i].Description = db.(notion.Database).Description[0].PlainText
			} else {
				dbs[i].Description = ""
			}
		}

		// TODO: get default database from storage
		// TODO: check if result contains default (compare by id); if not = error, missing integration on db

		tui.GetDbsListTUI(dbs)
	},
}
