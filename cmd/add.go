package cmd

import (
	"fmt"

	"github.com/ChmaraX/notidb/internal/notion"
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

var args cmdArgs

const DefaultTitlePropKey = "title"
const GreenCheckMark = "\033[32m✓\033[0m"

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

func createEntryFromArgs(a cmdArgs) notion.DatabaseEntry {
	entry := notion.DatabaseEntry{
		Props:  make(notionapi.Properties),
		Blocks: make([]notionapi.Block, 0),
	}

	if a.title != "" {
		entry.Props[DefaultTitlePropKey] = notion.CreateTitleProperty(a.title)
	}

	if a.content != "" {
		entry.Blocks = append(entry.Blocks, notion.CreateContentBlock(a.content))
	}

	return entry
}

func saveEntry(dbId string, entry notion.DatabaseEntry) tui.Response {
	page, err := notion.AddDatabaseEntry(dbId, entry)
	id := "save"

	if err != nil {
		return tui.Response{Id: id, Data: nil, Err: fmt.Errorf("error saving entry: %v", err)}
	}

	return tui.Response{Id: id, Data: page.URL, Err: nil}
}

func wrappedSaveEntry(dbId string, entry notion.DatabaseEntry) func() tui.Response {
	return func() tui.Response {
		return saveEntry(dbId, entry)
	}
}

func createEntry() notion.DatabaseEntry {
	if args.title == "" && args.content == "" {
		schema, err := notion.GetDatabaseSchema(args.dbId)
		if err != nil {
			fmt.Printf("Error getting DB schema: %v\n", err)
		}
		return tui.InitForm(schema)
	}
	return createEntryFromArgs(args)
}

/**
 * Shortcut for args provided like this `notidb add "foo" "bar"`
 * where "foo" is title and "bar" is content
 */
func loadShortcutArgs(arguments []string) {
	// check if args are not already set
	if args.title != "" || args.content != "" {
		return
	}
	if len(arguments) > 0 {
		args.title = arguments[0]
	}
	if len(arguments) > 1 {
		args.content = arguments[1]
	}
}

var addEntryCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "Adds a new entry to the database",
	Run: func(cmd *cobra.Command, arguments []string) {
		if err := args.validateDefaultDb(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		loadShortcutArgs(arguments)

		entry := createEntry()

		if entry.Props == nil && entry.Blocks == nil {
			fmt.Println("No content to save")
			return
		}

		m := tui.NewLoadingModel("Saving to Notion", wrappedSaveEntry(args.dbId, entry))
		res := m.GetResponse("save")

		if res.Err != nil {
			fmt.Printf("\n%s\n", res.Err)
			return
		}

		url := res.Data.(string)
		fmt.Printf("\n %s Saved: %s\n\n", GreenCheckMark, url)
	},
}

func init() {
	addEntryCmd.Flags().StringVarP(&args.title, "title", "t", "", "Title of the new entry")
	addEntryCmd.Flags().StringVarP(&args.content, "content", "c", "", "Content of the new entry")
}
