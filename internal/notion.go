package internal

import (
	"context"

	"github.com/dstotijn/go-notion"
)

var NotionClient *notion.Client

type NotionDb struct {
	Title       string
	Description string
	Id          string
}

func GetAllNotionDbs() []NotionDb {
	res, err := NotionClient.Search(
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

	dbs := make([]NotionDb, len(res.Results))
	for i, db := range res.Results {
		dbs[i] = parseNotionDb(db.(notion.Database))
	}

	return dbs
}

func parseNotionDb(db notion.Database) NotionDb {
	var parsedDb NotionDb
	parsedDb.Id = db.ID
	parsedDb.Title = db.Title[0].PlainText
	if len(db.Description) > 0 {
		parsedDb.Description = db.Description[0].PlainText
	} else {
		parsedDb.Description = ""
	}
	return parsedDb
}

// connect to notion
func CreateNotionClient() {
	c, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	NotionClient = notion.NewClient(c.ApiKey)
}
