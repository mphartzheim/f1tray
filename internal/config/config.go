package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// Preferences defines user-configurable application settings.
type Preferences struct {
	CloseBehavior  string `json:"close_behavior"`    // "exit" or "minimize"
	HideOnOpen     bool   `json:"hide_on_open"`      // if true, the window is hidden on launch
	DebugMode      bool   `json:"debug_mode"`        // if true, debug mode is enabled
	EnableSound    bool   `json:"enable_sound"`      // if true, play system sounds
	Use24HourClock bool   `json:"use_24_hour_clock"` // if true, display time in 24-hour format
	Theme          string `json:"theme"`             // selected theme (e.g., "Dark", "Light", etc.)
}

// DefaultPreferences provides fallback settings when no config file is present.
var DefaultPreferences = Preferences{
	CloseBehavior:  "minimize",
	HideOnOpen:     true,
	DebugMode:      false,
	EnableSound:    true,
	Use24HourClock: false,
	Theme:          "Dark",
}

var (
	instance *Preferences
	once     sync.Once
)

// loadConfig loads the configuration from disk (or returns default if none exists).
func loadConfig() *Preferences {
	configPath := getConfigPath()

	// Check if file exists.
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		_ = os.MkdirAll(filepath.Dir(configPath), os.ModePerm)
		_ = SaveConfig(DefaultPreferences)
		return &DefaultPreferences
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return &DefaultPreferences
	}

	var prefs Preferences
	if err := json.Unmarshal(data, &prefs); err != nil {
		return &DefaultPreferences
	}
	return &prefs
}

// Get returns the global Preferences instance, loading it once.
func Get() *Preferences {
	once.Do(func() {
		instance = loadConfig()
	})
	return instance
}

// Set updates the global Preferences instance in memory and on disk.
func Set(prefs *Preferences) error {
	instance = prefs
	return SaveConfig(*prefs)
}

// Reload clears the singleton instance so that it gets reloaded from disk.
func Reload() {
	once = sync.Once{} // reset the once so that Get() reloads the config
	Get()
}

// SaveConfig writes the preferences to the configuration file.
func SaveConfig(prefs Preferences) error {
	configPath := getConfigPath()
	data, err := json.MarshalIndent(prefs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

// getConfigPath returns the full path to the configuration file.
func getConfigPath() string {
	dirname, _ := os.UserConfigDir()
	return filepath.Join(dirname, "f1tray", "config.json")
}
