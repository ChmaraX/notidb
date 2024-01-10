package notion

import (
	"context"

	"github.com/ChmaraX/notidb/internal"
	"github.com/jomei/notionapi"
)

var NotionClient *notionapi.Client

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
	c, err := internal.LoadConfig()
	if err != nil {
		panic(err)
	}
	NotionClient = notionapi.NewClient(notionapi.Token(c.ApiKey))
}
