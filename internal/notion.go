package internal

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jomei/notionapi"
)

var NotionClient *notionapi.Client

var NotionPagePropTypes = []string{
	"rich_text",
	"number",
	"select",
	"multi_select",
	"date",
	"people",
	"file",
	"checkbox",
	"url",
	"email",
	"phone_number",
	"formula",
	"relation",
	"rollup",
	"created_time",
	"created_by",
	"last_edited_time",
	"last_edited_by",
}

func GetSupportedPagePropTypes() []string {
	return []string{
		"rich_text",
		"number",
		"select",
		"multi_select",
		"date",
		"checkbox",
		"email",
		"phone_number",
	}
}

// no better way to convert PropertyConfigs to Properties atm
func ConvertPropertyConfigsToProps(schema notionapi.PropertyConfigs) notionapi.Properties {
	schemaJson, err := json.Marshal(schema)
	if err != nil {
		fmt.Printf("failed to marshal schema to JSON: %v", err)
	}

	var props notionapi.Properties
	err = json.Unmarshal(schemaJson, &props)
	if err != nil {
		fmt.Printf("failed to unmarshal JSON to Properties: %v", err)
	}

	return props
}

func CreateContentBlock(content string) notionapi.Block {
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
						Content: content,
					},
				},
			},
		},
	}
}

func CreateTitleProperty(title string) notionapi.TitleProperty {
	return notionapi.TitleProperty{Title: []notionapi.RichText{
		{
			Type: "text",
			Text: &notionapi.Text{
				Content: title,
			},
		},
	},
	}
}

func GetAllNotionDbs() ([]notionapi.Database, error) {
	res, err := NotionClient.Search.Do(context.Background(), &notionapi.SearchRequest{
		Filter: notionapi.SearchFilter{
			Value:    "database",
			Property: "object",
		},
	})
	if err != nil {
		return nil, err
	}

	databases := make([]notionapi.Database, len(res.Results))
	for i, obj := range res.Results {
		if db, ok := obj.(*notionapi.Database); ok {
			databases[i] = *db
		}
	}

	return databases, nil
}

func GetDatabaseSchema(dbId string) (notionapi.PropertyConfigs, error) {
	db, err := NotionClient.Database.Get(context.Background(), notionapi.DatabaseID(dbId))
	if err != nil {
		return nil, err
	}
	return db.Properties, nil
}

type DatabaseEntry struct {
	Props  notionapi.Properties
	Blocks []notionapi.Block
}

func AddDatabaseEntry(dbId string, entry DatabaseEntry) (notionapi.Page, error) {
	page, err := NotionClient.Page.Create(context.Background(), &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       "database_id",
			DatabaseID: notionapi.DatabaseID(dbId),
		},
		Properties: entry.Props,
		Children:   entry.Blocks,
	})
	if err != nil {
		return notionapi.Page{}, err
	}
	return *page, nil
}

func CreateNotionClient() {
	c, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	NotionClient = notionapi.NewClient(notionapi.Token(c.ApiKey))
}
