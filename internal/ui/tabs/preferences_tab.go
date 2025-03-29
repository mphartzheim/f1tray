package tabs

import (
	"sort"

	"f1tray/internal/config"
	"f1tray/internal/processes"
	"f1tray/internal/ui/themes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CreatePreferencesTab builds a preferences form for toggling app behavior like close mode,
// startup visibility, sounds, and debug mode. It reads and writes the global configuration via config.Get()/Set().
func CreatePreferencesTab(onSave func(config.Preferences), refreshUpcomingTab func()) fyne.CanvasObject {
	// Get the current global preferences.
	prefs := config.Get()

	// Dynamically build theme options from the AvailableThemes registry.
	availableThemes := themes.AvailableThemes()
	themeOptions := make([]string, 0, len(availableThemes))
	for name := range availableThemes {
		themeOptions = append(themeOptions, name)
	}
	sort.Strings(themeOptions) // sort options alphabetically

	// Map the selected string to a theme instance.
	mapTheme := func(selected string) fyne.Theme {
		if t, ok := availableThemes[selected]; ok {
			return t
		}
		return theme.DefaultTheme()
	}

	// Create the theme drop-down with label "Theme:".
	selectTheme := widget.NewSelect(themeOptions, func(selected string) {
		prefs.Theme = selected
		_ = config.Set(prefs) // update global preferences
		onSave(*prefs)        // trigger onSave callback with updated prefs
		fyne.CurrentApp().Settings().SetTheme(mapTheme(selected))
	})
	selectTheme.SetSelected(prefs.Theme)
	themeRow := container.NewHBox(widget.NewLabel("Theme:"), selectTheme)

	// Close on exit checkbox.
	closeCheckbox := widget.NewCheck("Close on exit?", func(checked bool) {
		if checked {
			prefs.CloseBehavior = "exit"
		} else {
			prefs.CloseBehavior = "minimize"
		}
		_ = config.Set(prefs)
		onSave(*prefs)
	})
	closeCheckbox.SetChecked(prefs.CloseBehavior == "exit")

	// Hide on open checkbox.
	hideCheckbox := widget.NewCheck("Hide on open?", func(checked bool) {
		prefs.HideOnOpen = checked
		_ = config.Set(prefs)
		onSave(*prefs)
	})
	hideCheckbox.SetChecked(prefs.HideOnOpen)

	// Sound settings.
	testButton := widget.NewButton("Test", func() {
		processes.PlayNotificationSound()
	})
	soundCheckbox := widget.NewCheck("Enable sounds?", func(checked bool) {
		prefs.EnableSound = checked
		if checked {
			testButton.Enable()
		} else {
			testButton.Disable()
		}
		_ = config.Set(prefs)
		onSave(*prefs)
	})
	soundCheckbox.SetChecked(prefs.EnableSound)
	if !prefs.EnableSound {
		testButton.Disable()
	}
	soundRow := container.NewHBox(soundCheckbox, testButton)

	// 24-hour clock checkbox.
	timeFormatCheckbox := widget.NewCheck("Use 24-hour clock?", func(checked bool) {
		prefs.Use24HourClock = checked
		_ = config.Set(prefs)
		onSave(*prefs)
		// Trigger the Upcoming Tab to refresh so the times are redrawn immediately.
		refreshUpcomingTab()
	})
	timeFormatCheckbox.SetChecked(prefs.Use24HourClock)

	// Debug mode checkbox.
	debugCheckbox := widget.NewCheck("Debug Mode?", func(checked bool) {
		prefs.DebugMode = checked
		_ = config.Set(prefs)
		onSave(*prefs)
	})
	debugCheckbox.SetChecked(prefs.DebugMode)

	return container.NewVBox(
		themeRow,
		closeCheckbox,
		hideCheckbox,
		soundRow,
		timeFormatCheckbox,
		debugCheckbox,
	)
}
