package internal

import (
	"github.com/dstotijn/go-notion"
)

var NotionClient *notion.Client

// connect to notion
func CreateNotionClient() {
	c, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	NotionClient = notion.NewClient(c.ApiKey)
}
