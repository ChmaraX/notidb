package cmd

import (
	"fmt"

	"github.com/ChmaraX/notidb/internal"
	"github.com/ChmaraX/notidb/internal/settings"
	"github.com/ChmaraX/notidb/internal/tui"
	"github.com/jomei/notionapi"
	"github.com/spf13/cobra"
)

type cmdArgs struct {
	title   string
	content string
	dbId    string
}

const DefaultTitlePropKey = "title"

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
		entry.Props[DefaultTitlePropKey] = internal.CreateTitleProperty(a.title)
	}

	if a.content != "" {
		entry.Blocks = append(entry.Blocks, internal.CreateContentBlock(a.content))
	}

	return entry
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
		// TODO:
		// 2. form.go - form for adding entry,
		// 3. load.go - generic loading model, with funcs as params and also load/result validating func
		if title == "" && content == "" {
			tui.InitForm(dbId)
			return
		}

		entry := createEntryFromArgs(args)
		tui.InitSave(dbId, entry)
	},
}

var args cmdArgs

func init() {
	addEntryCmd.Flags().StringVarP(&args.title, "title", "t", "", "Title of the new entry")
	addEntryCmd.Flags().StringVarP(&args.content, "content", "c", "", "Content of the new entry")
	addEntryCmd.Flags().StringVarP(&args.dbId, "database", "d", "", "ID of the database to add entry to")
}
