package preferences

import (
	"github.com/mphartzheim/f1tray/internal/appstate"
	"github.com/mphartzheim/f1tray/internal/config"
	"github.com/mphartzheim/f1tray/internal/themes"

	"fyne.io/fyne/v2"
)

// ApplyThemeSelection sets the active theme in the app
func ApplyThemeSelection(app fyne.App, selected string) {
	if theme, ok := themes.AvailableThemes()[selected]; ok {
		app.Settings().SetTheme(theme)
	}
}

// ApplyWindowSettings applies window preferences (such as minimize to tray)
func ApplyWindowSettings(state *appstate.AppState) {
	prefs := config.Get()
	if prefs.Window.MinimizeToTray {
		// Future tray logic could go here
	}
}

// SetClockFormat could affect future time display format (24h vs 12h)
func SetClockFormat(use24 bool) {
	// Placeholder for future time formatting behavior
}
