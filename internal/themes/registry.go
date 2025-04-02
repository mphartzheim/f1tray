package themes

import (
	"sort"

	"fyne.io/fyne/v2"
)

// AvailableThemes returns a map of theme names to their implementations.
func AvailableThemes() map[string]fyne.Theme {
	return map[string]fyne.Theme{
		// Core options
		"System": SystemTheme{},
		"Dark":   DarkTheme{},
		"Light":  LightTheme{},

		// Teams
		"Alpine":       AlpineTheme{},
		"Aston Martin": AstonMartinTheme{},
		"Ferrari":      FerrariTheme{},
		"Haas":         HaasTheme{},
		"McLaren":      McLarenTheme{},
		"Mercedes":     MercedesTheme{},
		"Racing Bulls": RacingBullsTheme{},
		"Red Bull":     RedBullTheme{},
		"Sauber":       SauberTheme{},
		"Williams":     WilliamsTheme{},
	}
}

// SortedThemeList returns theme names with System, Dark, Light first, then teams alphabetically.
func SortedThemeList() []string {
	priority := []string{"System", "Dark", "Light"}

	// Gather remaining keys
	var others []string
	for name := range AvailableThemes() {
		if name != "System" && name != "Dark" && name != "Light" {
			others = append(others, name)
		}
	}

	sort.Strings(others)

	return append(priority, others...)
}
