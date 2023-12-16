package cmd

import (
	"fmt"

	"github.com/ChmaraX/notidb/internal"
	"github.com/ChmaraX/notidb/internal/settings"
	"github.com/ChmaraX/notidb/internal/utils"
	"github.com/jomei/notionapi"
	"github.com/spf13/cobra"
)

type cmdArgs struct {
	title   string
	content string
	dbId    string
}

func (a *cmdArgs) validateDefaultDb() error {
	if a.dbId == "" {
		dbId, err := settings.GetDefaultDatabase()
		if err != nil {
			return err
		}
		a.dbId = dbId
	}
	return nil
}

func createEntryFromArgs(a cmdArgs) internal.DatabaseEntry {
	entry := internal.DatabaseEntry{
		Props:  make(notionapi.Properties),
		Blocks: make([]notionapi.Block, 0),
	}

	if a.title != "" {
		entry.Props["title"] = internal.CreateTitleProperty(a.title)
	}

	if a.content != "" {
		entry.Blocks = append(entry.Blocks, internal.CreateContentBlock(a.content))
	}

	return entry
}

func createEmptyEntry(dbId string) internal.DatabaseEntry {
	schema, err := internal.GetDatabaseSchema(dbId)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	schema = filterSupportedProps(schema)
	props := internal.ConvertPropertyConfigsToProps(schema)
	props["title"] = internal.CreateTitleProperty("")
	blocks := []notionapi.Block{internal.CreateContentBlock("")}

	return internal.DatabaseEntry{
		Props:  props,
		Blocks: blocks,
	}
}

func addEntry(dbId string, entry internal.DatabaseEntry) (notionapi.Page, error) {
	res, err := internal.AddDatabaseEntry(dbId, entry)
	if err != nil {
		return notionapi.Page{}, fmt.Errorf("failed to add database entry: %w", err)
	}

	return res, nil
}

func filterSupportedProps(schema notionapi.PropertyConfigs) notionapi.PropertyConfigs {
	supportedPropTypes := internal.GetSupportedPagePropTypes()
	filteredSchema := make(notionapi.PropertyConfigs)
	for key, value := range schema {
		if utils.Contains(supportedPropTypes, string(value.GetType())) {
			filteredSchema[key] = value
		}
	}
	return filteredSchema
}

var addEntryCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "Adds a new entry to a database",
	Run: func(cmd *cobra.Command, arguments []string) {
		if err := args.validateDefaultDb(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		dbId, title, content := args.dbId, args.title, args.content

		if title == "" && content == "" {

			entryForm := createEmptyEntry(dbId)
			fmt.Printf("%+v\n", entryForm)
			// TODO: open form
			// tui.InitForm(schema)
			return
		}

		entry := createEntryFromArgs(args)
		page, err := addEntry(dbId, entry)
		if err != nil {
			fmt.Printf("Error adding entry: %v\n", err)
			return
		}
		fmt.Printf("Entry successfully added: %s\n", page.URL)
	},
}

var args cmdArgs

func init() {
	addEntryCmd.Flags().StringVarP(&args.title, "title", "t", "", "Title of the new entry")
	addEntryCmd.Flags().StringVarP(&args.content, "content", "c", "", "Content of the new entry")
	addEntryCmd.Flags().StringVarP(&args.dbId, "database", "d", "", "ID of the database to add entry to")
}
