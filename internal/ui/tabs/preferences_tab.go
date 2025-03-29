package tabs

import (
	"f1tray/internal/config"
	"f1tray/internal/processes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	// Import your custom themes package.
	"f1tray/internal/ui/themes"
)

// themeMap maps our theme option strings to actual theme instances.
var themeMap = map[string]fyne.Theme{
	"Dark":  themes.DarkTheme{},  // our default dark theme
	"Light": themes.LightTheme{}, // using the default (light) theme
}

// mapTheme returns the theme instance based on the selected string.
// If the selection is not found, it returns the default theme.
func mapTheme(selected string) fyne.Theme {
	if t, ok := themeMap[selected]; ok {
		return t
	}
	return theme.DefaultTheme()
}

// CreatePreferencesTab builds a preferences form for toggling app behavior like close mode, startup visibility, sounds, and debug mode.
func CreatePreferencesTab(currentPrefs config.Preferences, onSave func(config.Preferences), refreshUpcomingTab func()) fyne.CanvasObject {
	// Define theme options (update these to match your available themes)
	themeOptions := []string{"System", "Dark", "Light"}

	// Create the theme drop-down with label "Theme:".
	selectTheme := widget.NewSelect(themeOptions, func(selected string) {
		currentPrefs.Theme = selected
		onSave(currentPrefs)
		// Update the app theme immediately.
		fyne.CurrentApp().Settings().SetTheme(mapTheme(selected))
	})
	// Ensure the drop-down always shows the currently set theme.
	selectTheme.SetSelected(currentPrefs.Theme)
	// Add a label to the right indicating light theme is unsupported.
	themeRow := container.NewHBox(widget.NewLabel("Theme:"), selectTheme)

	isExit := currentPrefs.CloseBehavior == "exit"
	closeCheckbox := widget.NewCheck("Close on exit?", func(checked bool) {
		if checked {
			currentPrefs.CloseBehavior = "exit"
		} else {
			currentPrefs.CloseBehavior = "minimize"
		}
		onSave(currentPrefs)
	})
	closeCheckbox.SetChecked(isExit)

	hideCheckbox := widget.NewCheck("Hide on open?", func(checked bool) {
		currentPrefs.HideOnOpen = checked
		onSave(currentPrefs)
	})
	hideCheckbox.SetChecked(currentPrefs.HideOnOpen)

	// Create the testButton first.
	testButton := widget.NewButton("Test", func() {
		processes.PlayNotificationSound(currentPrefs)
	})

	// Create the soundCheckbox that references the testButton.
	soundCheckbox := widget.NewCheck("Enable sounds?", func(checked bool) {
		currentPrefs.EnableSound = checked
		if checked {
			testButton.Enable()
		} else {
			testButton.Disable()
		}
		onSave(currentPrefs)
	})
	soundCheckbox.SetChecked(currentPrefs.EnableSound)

	// Set the initial state of testButton.
	if !currentPrefs.EnableSound {
		testButton.Disable()
	}

	soundRow := container.NewHBox(
		soundCheckbox,
		testButton,
	)

	// Update the time format checkbox callback to trigger the Upcoming Tab's refresh.
	timeFormatCheckbox := widget.NewCheck("Use 24-hour clock?", func(checked bool) {
		currentPrefs.Use24HourClock = checked
		onSave(currentPrefs)
		// Trigger the Upcoming Tab to refresh so the times are redrawn immediately.
		refreshUpcomingTab()
	})
	timeFormatCheckbox.SetChecked(currentPrefs.Use24HourClock)

	debugCheckbox := widget.NewCheck("Debug Mode?", func(checked bool) {
		currentPrefs.DebugMode = checked
		onSave(currentPrefs)
	})
	debugCheckbox.SetChecked(currentPrefs.DebugMode)

	return container.NewVBox(
		themeRow, // Theme selection drop-down and label at the top.
		closeCheckbox,
		hideCheckbox,
		soundRow,
		timeFormatCheckbox,
		debugCheckbox,
	)
}
