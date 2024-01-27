package cmd

import (
	"fmt"
	"os"

	"github.com/ChmaraX/notidb/internal/keyring"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"i"},
	Short:   "Initialize NotiDB CLI",
	Run: func(cmd *cobra.Command, args []string) {

		authManager, err := keyring.NewKeyringManager()
		if err != nil {
			fmt.Printf("Error initializing AuthManager: %s\n", err)
			return
		}

		fmt.Print("Please enter your Notion API key: ")
		apiKeyBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Printf("\nError reading API key: %s\n", err)
			return
		}
		apiKey := string(apiKeyBytes)

		fmt.Println() // Print a new line after the user's input

		if err := authManager.SaveAPIKey(apiKey); err != nil {
			fmt.Printf("Error saving API key: %s\n", err)
			return
		}

		// prompt for default database
		setDefaultDbCmd.Run(cmd, args)

		fmt.Printf("\n %s NotiDB CLI initialized\n\n", GreenCheckMark)
	},
}
