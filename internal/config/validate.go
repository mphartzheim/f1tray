package config

import "github.com/mphartzheim/f1tray/internal/themes"

func SanitizePreferences(p Preferences) Preferences {
	// Validate theme selection
	if _, ok := themes.AvailableThemes()[p.Themes.Selected]; !ok {
		p.Themes.Selected = "System"
	}

	// Ensure no negative notification timing
	if p.Notifications.NotifyMinutesBefore < 0 {
		p.Notifications.NotifyMinutesBefore = 10
	}

	return p
}
