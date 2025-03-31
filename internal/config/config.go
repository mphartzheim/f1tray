package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// Preferences defines user-configurable application settings.
type Preferences struct {
	Window              WindowPreferences       `json:"window"`
	Themes              ThemesPreferences       `json:"theme"`
	Clock               ClockPreferences        `json:"clock"`
	Debug               DebugPreferences        `json:"debug"`
	Notifications       NotificationPreferences `json:"notifications"`
	FavoriteDrivers     []string                `json:"favorite_drivers"`
	FavoriteConstructor string                  `json:"favorite_constructor"`
}

// WindowPreferences groups window-related settings.
type WindowPreferences struct {
	CloseBehavior string `json:"close_behavior"` // "exit" or "minimize"
	HideOnOpen    bool   `json:"hide_on_open"`   // if true, the window is hidden on launch
}

// ThemesPreferences groups theme-related settings.
type ThemesPreferences struct {
	Theme string `json:"theme"` // e.g., "Dark", "Light", etc.
}

// ClockPreferences groups clock-related settings.
type ClockPreferences struct {
	Use24Hour bool `json:"use_24_hour_clock"` // if true, display time in 24-hour format
}

// DebugPreferences groups debug-related settings.
type DebugPreferences struct {
	Enabled bool `json:"debug_mode"` // if true, debug mode is enabled
}

// NotificationPreferences groups notification-related settings.
type NotificationPreferences struct {
	Practice   *SessionNotificationSettings `json:"practice"`
	Qualifying *SessionNotificationSettings `json:"qualifying"`
	Race       *SessionNotificationSettings `json:"race"`
}

// SessionNotificationSettings defines notification settings for a session.
type SessionNotificationSettings struct {
	NotifyOnStart    bool   `json:"notify_on_start"`     // notify at session start
	PlaySoundOnStart bool   `json:"play_sound_on_start"` // play sound at session start
	NotifyBefore     bool   `json:"notify_before"`       // notify before session
	BeforeValue      int    `json:"before_value"`        // numeric value for before notification
	BeforeUnit       string `json:"before_unit"`         // "minutes" or "hours"
	PlaySoundBefore  bool   `json:"play_sound_before"`   // play sound before session
}

// defaultSessionNotificationSettings returns default settings for a session.
func defaultSessionNotificationSettings() *SessionNotificationSettings {
	return &SessionNotificationSettings{
		NotifyOnStart:    false,
		PlaySoundOnStart: false,
		NotifyBefore:     false,
		BeforeValue:      10,
		BeforeUnit:       "minutes",
		PlaySoundBefore:  false,
	}
}

// DefaultPreferences provides fallback settings when no config file is present.
var DefaultPreferences = Preferences{
	Window: WindowPreferences{
		CloseBehavior: "minimize",
		HideOnOpen:    true,
	},
	Themes: ThemesPreferences{
		Theme: "Dark",
	},
	Clock: ClockPreferences{
		Use24Hour: false,
	},
	Debug: DebugPreferences{
		Enabled: false,
	},
	Notifications: NotificationPreferences{
		Practice:   defaultSessionNotificationSettings(),
		Qualifying: defaultSessionNotificationSettings(),
		Race:       defaultSessionNotificationSettings(),
	},
	FavoriteDrivers:     []string{}, // default empty favorites list
	FavoriteConstructor: "",         // default no favorite constructor
}

var (
	instance *Preferences
	once     sync.Once
)

// loadConfig loads config from disk (or defaults) and validates/migrates legacy settings.
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
		// Attempt to load legacy config.
		var legacy legacyPreferences
		if err := json.Unmarshal(data, &legacy); err == nil {
			prefs = migrateLegacy(legacy)
			_ = SaveConfig(prefs)
			return &prefs
		}
		return &DefaultPreferences
	}

	// Validate the loaded config.
	prefs = validatePreferences(prefs)
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

// validatePreferences ensures that all nested preference values have sane defaults.
func validatePreferences(prefs Preferences) Preferences {
	// Validate Window settings.
	if prefs.Window.CloseBehavior == "" {
		prefs.Window.CloseBehavior = DefaultPreferences.Window.CloseBehavior
	}
	// Validate Themes settings.
	if prefs.Themes.Theme == "" {
		prefs.Themes.Theme = DefaultPreferences.Themes.Theme
	}
	// For Notifications, if any session settings are nil, replace with defaults.
	if prefs.Notifications.Practice == nil {
		prefs.Notifications.Practice = defaultSessionNotificationSettings()
	}
	if prefs.Notifications.Qualifying == nil {
		prefs.Notifications.Qualifying = defaultSessionNotificationSettings()
	}
	if prefs.Notifications.Race == nil {
		prefs.Notifications.Race = defaultSessionNotificationSettings()
	}
	// Ensure FavoriteDrivers is non-nil.
	if prefs.FavoriteDrivers == nil {
		prefs.FavoriteDrivers = []string{}
	}
	// Ensure FavoriteConstructor is set (default to empty string if missing).
	if prefs.FavoriteConstructor == "" {
		prefs.FavoriteConstructor = ""
	}
	return prefs
}

// legacyPreferences represents the old flat configuration structure.
type legacyPreferences struct {
	CloseBehavior  string `json:"close_behavior"`
	HideOnOpen     bool   `json:"hide_on_open"`
	DebugMode      bool   `json:"debug_mode"`
	EnableSound    bool   `json:"enable_sound"` // Removed in new version
	Use24HourClock bool   `json:"use_24_hour_clock"`
	Theme          string `json:"theme"`
}

// migrateLegacy migrates a legacyPreferences struct to the new Preferences struct.
func migrateLegacy(old legacyPreferences) Preferences {
	return Preferences{
		Window: WindowPreferences{
			CloseBehavior: old.CloseBehavior,
			HideOnOpen:    old.HideOnOpen,
		},
		Themes: ThemesPreferences{
			Theme: old.Theme,
		},
		Clock: ClockPreferences{
			Use24Hour: old.Use24HourClock,
		},
		Debug: DebugPreferences{
			Enabled: old.DebugMode,
		},
		Notifications: NotificationPreferences{
			Practice:   defaultSessionNotificationSettings(),
			Qualifying: defaultSessionNotificationSettings(),
			Race:       defaultSessionNotificationSettings(),
		},
		FavoriteDrivers:     []string{}, // no legacy data; start with an empty list
		FavoriteConstructor: "",
	}
}
