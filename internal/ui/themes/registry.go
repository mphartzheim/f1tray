package themes

import (
	"fyne.io/fyne/v2"
)

// AvailableThemes returns a map of theme names to their implementations.
func AvailableThemes() map[string]fyne.Theme {
	return map[string]fyne.Theme{
		"System": SystemTheme{}, // Your custom System theme.
		"Dark":   DarkTheme{},   // Your custom Dark theme.
		"Light":  LightTheme{},  // Your custom Light theme.
	}
}
