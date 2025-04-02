package config

type Preferences struct {
	Themes              ThemesPreferences       `json:"themes"`
	Window              WindowPreferences       `json:"window"`
	Clock               ClockPreferences        `json:"clock"`
	Notifications       NotificationPreferences `json:"notifications"`
	FavoriteDrivers     []string                `json:"favorite_drivers"`
	FavoriteConstructor string                  `json:"favorite_constructor"`
}

type ThemesPreferences struct {
	Selected string `json:"selected"`
}

type WindowPreferences struct {
	MinimizeToTray bool `json:"minimize_to_tray"`
}

type ClockPreferences struct {
	Use24Hour bool `json:"use_24_hour"`
}

type NotificationPreferences struct {
	EnableSessionStart  bool `json:"enable_session_start"`
	NotifyMinutesBefore int  `json:"notify_minutes_before"`
}
