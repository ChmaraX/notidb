package cmd

import (
	"fmt"

	"github.com/ChmaraX/notidb/internal"
	"github.com/ChmaraX/notidb/internal/settings"
	"github.com/jomei/notionapi"
	"github.com/spf13/cobra"
)

type addEntryArgs struct {
	title string
	body  string
	dbId  string
}

func (c *addEntryArgs) validateDefaultDb() error {
	if c.dbId == "" {
		dbId, err := settings.GetDefaultDatabase()
		if err != nil {
			return err
		}
		c.dbId = dbId
	}
	return nil
}

func addEntry(c addEntryArgs, schema notionapi.PropertyConfigs) (notionapi.Page, error) {
	if c.title != "" {
		props := notionapi.Properties{}
		// "title" is alias for any property that is of type "title" (unique per database)
		props["title"] = notionapi.TitleProperty{Title: []notionapi.RichText{
			{
				Type: "text",
				Text: &notionapi.Text{
					Content: c.title,
				},
			},
		}}
		res, err := internal.AddDatabaseEntry(c.dbId, props)
		if err != nil {
			return notionapi.Page{}, err
		}
		fmt.Printf("Added entry: %+v\n", res)
		return res, nil
	}

	// TODO: add body as block
	return notionapi.Page{}, nil
}

var addEntryCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "Adds a new entry to a database",
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.validateDefaultDb(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		schema, err := internal.GetDatabaseSchema(config.dbId)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if config.title == "" && config.body == "" {

			fmt.Printf("Schema: %+v\n", schema)
			// TODO: open form
			return
		}

		// add entry to db
		addEntry(config, schema)

		// -d flag (optional) - database id; if not provided, the default database is used
		// -t flag (optional) - title of the entry
		// -b flag (optional) - body of the entry
		// no flag, open form to fill in the entry to the default database
	},
}

var config addEntryArgs

func init() {
	addEntryCmd.Flags().StringVarP(&config.title, "title", "t", "", "Title of the new entry")
	addEntryCmd.Flags().StringVarP(&config.body, "body", "b", "", "Body of the new entry")
	addEntryCmd.Flags().StringVarP(&config.dbId, "database", "d", "", "ID of the database to add entry to")
}
