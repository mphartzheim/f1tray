package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Preferences struct {
	CloseBehavior string `json:"close_behavior"` // "exit" or "minimize"
}

var DefaultPreferences = Preferences{
	CloseBehavior: "minimize",
}

func LoadConfig() Preferences {
	configPath := getConfigPath()

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create config directory if needed
		_ = os.MkdirAll(filepath.Dir(configPath), os.ModePerm)
		// Save default config
		_ = SaveConfig(DefaultPreferences)
		return DefaultPreferences
	}

	// Read existing config
	file, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultPreferences
	}

	var prefs Preferences
	err = json.Unmarshal(file, &prefs)
	if err != nil {
		return DefaultPreferences
	}

	return prefs
}

func SaveConfig(prefs Preferences) error {
	configPath := getConfigPath()
	data, err := json.MarshalIndent(prefs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

func getConfigPath() string {
	dirname, _ := os.UserConfigDir()
	return filepath.Join(dirname, "f1tray", "config.json")
}
