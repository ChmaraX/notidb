package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type UserSettings struct {
	DefaultDatabaseId string `json:"defaultDatabase"`
}

func GetDefaultDatabase() (string, error) {
	settingsFilePath, err := getSettingsFilePath()
	if err != nil {
		return "", err
	}

	settings, err := readSettings(settingsFilePath)
	if err != nil {
		return "", err
	}

	return settings.DefaultDatabaseId, nil
}

func SetDefaultDatabase(dbId string) error {
	settingsFilePath, err := getSettingsFilePath()
	if err != nil {
		return err
	}

	settings, err := readSettings(settingsFilePath)
	if err != nil {
		return err
	}

	settings.DefaultDatabaseId = dbId
	return writeSettings(settings, settingsFilePath)
}

func EnsureSettingsFileExists() error {
	settingsFilePath, err := getSettingsFilePath()
	if err != nil {
		return err
	}
	appDir := filepath.Dir(settingsFilePath)

	// create the directory if it doesn't exist
	if err := os.MkdirAll(appDir, 0700); err != nil {
		return err
	}

	// check if the settings file exists, and create it if not
	if _, err := os.Stat(settingsFilePath); os.IsNotExist(err) {
		defaultSettings := UserSettings{
			DefaultDatabaseId: "-1",
		}
		return writeSettings(&defaultSettings, settingsFilePath)
	}

	return nil
}

func getSettingsFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(homeDir, ".notidb")
	return filepath.Join(appDir, "settings.json"), nil
}

func readSettings(filePath string) (*UserSettings, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var settings UserSettings
	err = json.Unmarshal(file, &settings)
	return &settings, err
}

func writeSettings(settings *UserSettings, filePath string) error {
	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0600)
}
