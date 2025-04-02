package userconfig

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// UserConfig represents saved user preferences
type UserConfig struct {
	SelectedTheme string `json:"selected_theme"`
}

var defaultConfig = UserConfig{
	SelectedTheme: "System",
}

func configFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appConfigDir := filepath.Join(configDir, "f1tray")
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(appConfigDir, "userconfig.json"), nil
}

// Load loads the user configuration from userconfig.json (or returns defaults)
func Load() (*UserConfig, error) {
	path, err := configFilePath()
	if err != nil {
		return &defaultConfig, nil
	}
	file, err := os.Open(path)
	if err != nil {
		return &defaultConfig, nil
	}
	defer file.Close()

	var cfg UserConfig
	err = json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return &defaultConfig, nil
	}
	return &cfg, nil
}

// Save writes the user configuration to userconfig.json
func Save(cfg *UserConfig) error {
	path, err := configFilePath()
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(cfg)
}
