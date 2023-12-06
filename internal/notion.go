package internal

import (
	"github.com/dstotijn/go-notion"
)

var NotionClient *notion.Client

type NotionDb struct {
	Title       string
	Description string
	Id          string
}

// connect to notion
func CreateNotionClient() {
	c, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	NotionClient = notion.NewClient(c.ApiKey)
}
