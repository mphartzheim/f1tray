package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Preferences struct {
	CloseBehavior  string `json:"close_behavior"`    // "exit" or "minimize"
	HideOnOpen     bool   `json:"hide_on_open"`      // if true, the window is hidden on launch
	DebugMode      bool   `json:"debug_mode"`        // if true, debug mode is enabled
	EnableSound    bool   `json:"enable_sound"`      // if true, play system sounds
	Use24HourClock bool   `json:"use_24_hour_clock"` // if true, display time in 24-hour format
	Theme          string `json:"theme"`             // selected theme (e.g., "Default", "CustomTheme", etc.)
}

var DefaultPreferences = Preferences{
	CloseBehavior:  "minimize",
	HideOnOpen:     true,
	DebugMode:      false,
	EnableSound:    true,
	Use24HourClock: false,
	Theme:          "System",
}

// LoadConfig reads the configuration file or returns default preferences if none exist.
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

// SaveConfig writes the given preferences to the configuration file.
func SaveConfig(prefs Preferences) error {
	configPath := getConfigPath()
	data, err := json.MarshalIndent(prefs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

// getConfigPath returns the full path to the user's F1Tray configuration file.
func getConfigPath() string {
	dirname, _ := os.UserConfigDir()
	return filepath.Join(dirname, "f1tray", "config.json")
}
