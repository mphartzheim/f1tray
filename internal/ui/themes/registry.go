package themes

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// AvailableThemes returns a map of theme names to their implementations.
func AvailableThemes() map[string]fyne.Theme {
	return map[string]fyne.Theme{
		"System": theme.DefaultTheme(), // Use Fyne's built-in system default theme.
		"Dark":   DarkTheme{},          // Your custom Dark theme.
		"Light":  LightTheme{},         // Your custom Light theme.
	}
}
