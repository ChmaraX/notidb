package cmd

import (
	"fmt"

	"github.com/ChmaraX/notidb/internal"
	"github.com/ChmaraX/notidb/internal/settings"
	"github.com/jomei/notionapi"
	"github.com/spf13/cobra"
)

type databaseEntry struct {
	title string
	body  string
	dbId  string
}

func (e *databaseEntry) validateDefaultDb() error {
	if e.dbId == "" {
		dbId, err := settings.GetDefaultDatabase()
		if err != nil {
			return err
		}
		e.dbId = dbId
	}
	return nil
}

func createBodyBlock(body string) notionapi.Block {
	return notionapi.ParagraphBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: "block",
			Type:   "paragraph",
		},
		Paragraph: notionapi.Paragraph{
			RichText: []notionapi.RichText{
				{
					Type: "text",
					Text: &notionapi.Text{
						Content: body,
					},
				},
			},
		},
	}
}

func createTitleProperty(title string) notionapi.Properties {
	return notionapi.Properties{
		// "title" is alias for any property of type "title" (unique per database)
		"title": notionapi.TitleProperty{Title: []notionapi.RichText{
			{
				Type: "text",
				Text: &notionapi.Text{
					Content: title,
				},
			},
		}},
	}
}

func addEntry(e databaseEntry) (notionapi.Page, error) {
	props := notionapi.Properties{}
	blocks := []notionapi.Block{}

	if e.title != "" {
		props = createTitleProperty(e.title)
	}

	if e.body != "" {
		blocks = append(blocks, createBodyBlock(e.body))
	}

	res, err := internal.AddDatabaseEntry(e.dbId, props, blocks)
	if err != nil {
		return notionapi.Page{}, err
	}

	return res, nil
}

var addEntryCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "Adds a new entry to a database",
	Run: func(cmd *cobra.Command, args []string) {
		if err := entry.validateDefaultDb(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		schema, err := internal.GetDatabaseSchema(entry.dbId)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if entry.title == "" && entry.body == "" {
			fmt.Printf("Schema: %+v\n", schema)
			// TODO: open form
			return
		}

		page, err := addEntry(entry)
		if err != nil {
			fmt.Printf("Error adding entry: %v\n", err)
			return
		}

		fmt.Printf("Added entry: %+v\n", page)
	},
}

var entry databaseEntry

func init() {
	addEntryCmd.Flags().StringVarP(&entry.title, "title", "t", "", "Title of the new entry")
	addEntryCmd.Flags().StringVarP(&entry.body, "body", "b", "", "Body of the new entry")
	addEntryCmd.Flags().StringVarP(&entry.dbId, "database", "d", "", "ID of the database to add entry to")
}
