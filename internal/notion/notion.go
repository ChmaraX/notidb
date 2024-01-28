package notion

import (
	"context"
	"fmt"
	"log"
	"net/http"

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

func CreateNotionClient(apiKey string) {
	err := validateNotionAPIKey(apiKey)
	if err != nil {
		log.Fatalf("error validating API key: %v \n", err)
	}
	NotionClient = notionapi.NewClient(notionapi.Token(apiKey))
}

func validateNotionAPIKey(apiKey string) error {
	url := "https://api.notion.com/v1/users/me"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Notion-Version", "2021-08-16")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API key is invalid or doesn't have necessary permissions")
	}

	return nil
}
