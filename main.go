package main

import (
	"log"

	"github.com/ChmaraX/notidb/cmd"
	"github.com/ChmaraX/notidb/internal"
	"github.com/ChmaraX/notidb/internal/settings"
)

func main() {
	if err := settings.EnsureSettingsFileExists(); err != nil {
		log.Fatalf("Failed to ensure settings file exists: %v", err)
	}
	internal.CreateNotionClient()
	cmd.Execute()
}
