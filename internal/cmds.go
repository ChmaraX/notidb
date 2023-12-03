package internal

import (
	"context"
	"fmt"

	"github.com/dstotijn/go-notion"
	"github.com/spf13/cobra"
)

func ListDatabasesCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list-databases",
		Aliases: []string{"ld", "ls"},
		Short:   "Lists all available databases",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("list-databases called")
			dbs, err := NotionClient.Search(
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
			for _, db := range dbs.Results {
				fmt.Println(db.(notion.Database).Title[0].PlainText, db.(notion.Database).ID)
			}
		},
	}
}
