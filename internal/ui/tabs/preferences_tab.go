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

// CreatePreferencesTab builds a preferences form for toggling app behavior like close mode, startup visibility, sounds, and debug mode.
func CreatePreferencesTab(currentPrefs config.Preferences, onSave func(config.Preferences), refreshUpcomingTab func()) fyne.CanvasObject {
	// Dynamically build theme options from the AvailableThemes registry.
	availableThemes := themes.AvailableThemes()
	themeOptions := make([]string, 0, len(availableThemes))
	for name := range availableThemes {
		themeOptions = append(themeOptions, name)
	}
	sort.Strings(themeOptions) // optional: sort alphabetically

	// Map the selected string to a theme instance.
	mapTheme := func(selected string) fyne.Theme {
		if t, ok := availableThemes[selected]; ok {
			return t
		}
		return theme.DefaultTheme()
	}

	// Create the theme drop-down with label "Theme:".
	selectTheme := widget.NewSelect(themeOptions, func(selected string) {
		currentPrefs.Theme = selected
		onSave(currentPrefs)
		// Update the app theme immediately.
		fyne.CurrentApp().Settings().SetTheme(mapTheme(selected))
	})
	// Ensure the drop-down always shows the currently set theme.
	selectTheme.SetSelected(currentPrefs.Theme)
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
		processes.PlayNotificationSound()
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
