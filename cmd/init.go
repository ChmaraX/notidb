package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/ChmaraX/notidb/internal/keyring"
	"github.com/ChmaraX/notidb/internal/notion"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"i"},
	Short:   "Initialize NotiDB CLI",
	Run: func(cmd *cobra.Command, args []string) {

		keyring, err := keyring.NewKeyringManager()
		if err != nil {
			log.Fatalf("Error initializing keyring: %s\n", err)
		}

		fmt.Print("Please enter your Notion API key: ")
		apiKeyBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatalf("\nError reading API key: %s\n", err)
		}
		apiKey := string(apiKeyBytes)

		fmt.Println() // Print a new line after the user's input

		if err := keyring.SaveAPIKey(apiKey); err != nil {
			log.Fatalf("Error saving API key: %s\n", err)
		}

		// prompt for default database
		notion.CreateNotionClient(apiKey)
		setDefaultDbCmd.Run(cmd, args)

		fmt.Printf("\n %s NotiDB CLI initialized\n\n", GreenCheckMark)
	},
}

func initNotionClient() {
	keyring, err := keyring.NewKeyringManager()
	if err != nil {
		log.Fatalf("Error initializing keyring: %s\n", err)
	}

	apiKey, err := keyring.GetAPIKey()
	if err != nil {
		log.Fatalf("Error getting API key: %s\nNotiDB CLI might not be initialized. Please run `notidb init` first.\n", err)
	}

	notion.CreateNotionClient(apiKey)
}
