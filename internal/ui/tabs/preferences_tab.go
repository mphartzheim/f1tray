package tabs

import (
	"fmt"
	"sort"
	"strconv"

	"f1tray/internal/config"
	"f1tray/internal/processes"
	"f1tray/internal/ui/themes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CreatePreferencesTab builds the Preferences tab UI with sub-tabs for Main and Notifications.
func CreatePreferencesTab(onSave func(config.Preferences), refreshUpcomingTab func()) fyne.CanvasObject {
	prefs := config.Get()

	mainTab := container.NewTabItem("Main", buildMainPreferences(prefs, onSave, refreshUpcomingTab))
	notificationsTab := container.NewTabItem("Notifications", buildNotificationPreferences(prefs, onSave))

	return container.NewAppTabs(mainTab, notificationsTab)
}

// buildMainPreferences returns the main preferences UI.
func buildMainPreferences(prefs *config.Preferences, onSave func(config.Preferences), refreshUpcomingTab func()) fyne.CanvasObject {
	// Theme selector.
	availableThemes := themes.AvailableThemes()
	themeOptions := make([]string, 0, len(availableThemes))
	for name := range availableThemes {
		themeOptions = append(themeOptions, name)
	}
	sort.Strings(themeOptions)

	mapTheme := func(selected string) fyne.Theme {
		if t, ok := availableThemes[selected]; ok {
			return t
		}
		return theme.DefaultTheme()
	}

	selectTheme := widget.NewSelect(themeOptions, func(selected string) {
		prefs.Themes.Theme = selected
		_ = config.Set(prefs)
		onSave(*prefs)
		fyne.CurrentApp().Settings().SetTheme(mapTheme(selected))
	})
	selectTheme.SetSelected(prefs.Themes.Theme)
	themeRow := container.NewHBox(widget.NewLabel("Theme:"), selectTheme)

	// Close on exit checkbox.
	closeCheckbox := widget.NewCheck("Close on exit?", func(checked bool) {
		if checked {
			prefs.Window.CloseBehavior = "exit"
		} else {
			prefs.Window.CloseBehavior = "minimize"
		}
		_ = config.Set(prefs)
		onSave(*prefs)
	})
	closeCheckbox.SetChecked(prefs.Window.CloseBehavior == "exit")

	// Hide on open checkbox.
	hideCheckbox := widget.NewCheck("Hide on open?", func(checked bool) {
		prefs.Window.HideOnOpen = checked
		_ = config.Set(prefs)
		onSave(*prefs)
	})
	hideCheckbox.SetChecked(prefs.Window.HideOnOpen)

	// Sound settings.
	testButton := widget.NewButton("Test", func() {
		processes.PlayNotificationSound()
	})
	soundCheckbox := widget.NewCheck("Enable sounds?", func(checked bool) {
		prefs.Sound.Enable = checked
		if checked {
			testButton.Enable()
		} else {
			testButton.Disable()
		}
		_ = config.Set(prefs)
		onSave(*prefs)
	})
	soundCheckbox.SetChecked(prefs.Sound.Enable)
	if !prefs.Sound.Enable {
		testButton.Disable()
	}
	soundRow := container.NewHBox(soundCheckbox, testButton)

	// 24-hour clock checkbox.
	timeFormatCheckbox := widget.NewCheck("Use 24-hour clock?", func(checked bool) {
		prefs.Clock.Use24Hour = checked
		_ = config.Set(prefs)
		onSave(*prefs)
		refreshUpcomingTab()
	})
	timeFormatCheckbox.SetChecked(prefs.Clock.Use24Hour)

	// Debug mode checkbox.
	debugCheckbox := widget.NewCheck("Debug Mode?", func(checked bool) {
		prefs.Debug.Enabled = checked
		_ = config.Set(prefs)
		onSave(*prefs)
	})
	debugCheckbox.SetChecked(prefs.Debug.Enabled)

	return container.NewVBox(
		themeRow,
		closeCheckbox,
		hideCheckbox,
		soundRow,
		timeFormatCheckbox,
		debugCheckbox,
	)
}

// buildNotificationPreferences returns the notification preferences UI.
func buildNotificationPreferences(prefs *config.Preferences, onSave func(config.Preferences)) fyne.CanvasObject {
	return container.NewVBox(
		buildSessionNotificationSection("Practice", prefs.Notifications.Practice, onSave),
		buildSessionNotificationSection("Qualifying", prefs.Notifications.Qualifying, onSave),
		buildSessionNotificationSection("Race", prefs.Notifications.Race, onSave),
	)
}

// buildSessionNotificationSection builds a section for session-specific notification settings.
func buildSessionNotificationSection(title string, sessionPrefs *config.SessionNotificationSettings, onSave func(config.Preferences)) fyne.CanvasObject {
	// --- Row for "At Session Start" ---
	notifyStartCheck := widget.NewCheck("Notify at session start", func(checked bool) {
		sessionPrefs.NotifyOnStart = checked
		_ = config.Set(config.Get())
		onSave(*config.Get())
	})
	notifyStartCheck.SetChecked(sessionPrefs.NotifyOnStart)

	playSoundStartCheck := widget.NewCheck("Play sound", func(checked bool) {
		sessionPrefs.PlaySoundOnStart = checked
		_ = config.Set(config.Get())
		onSave(*config.Get())
	})
	playSoundStartCheck.SetChecked(sessionPrefs.PlaySoundOnStart)

	testStartButton := widget.NewButton("Test", func() {
		fmt.Printf("Testing %s session start notification\n", title)
		// Insert test logic here.
	})

	startRow := container.NewHBox(
		widget.NewLabelWithStyle("At Session Start:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		notifyStartCheck,
		playSoundStartCheck,
		testStartButton,
	)

	// --- Row for "Before Session" ---
	notifyBeforeCheck := widget.NewCheck("Notify before session", func(checked bool) {
		sessionPrefs.NotifyBefore = checked
		_ = config.Set(config.Get())
		onSave(*config.Get())
	})
	notifyBeforeCheck.SetChecked(sessionPrefs.NotifyBefore)

	beforeValueEntry := widget.NewEntry()
	beforeValueEntry.SetText(strconv.Itoa(sessionPrefs.BeforeValue))
	beforeValueEntry.OnChanged = func(val string) {
		if n, err := strconv.Atoi(val); err == nil {
			sessionPrefs.BeforeValue = n
			_ = config.Set(config.Get())
			onSave(*config.Get())
		}
	}

	beforeUnitSelect := widget.NewSelect([]string{"minutes", "hours"}, func(selected string) {
		sessionPrefs.BeforeUnit = selected
		_ = config.Set(config.Get())
		onSave(*config.Get())
	})
	beforeUnitSelect.SetSelected(sessionPrefs.BeforeUnit)

	playSoundBeforeCheck := widget.NewCheck("Play sound", func(checked bool) {
		sessionPrefs.PlaySoundBefore = checked
		_ = config.Set(config.Get())
		onSave(*config.Get())
	})
	playSoundBeforeCheck.SetChecked(sessionPrefs.PlaySoundBefore)

	testBeforeButton := widget.NewButton("Test", func() {
		fmt.Printf("Testing %s before session notification\n", title)
		// Insert test logic here.
	})

	beforeRow := container.NewHBox(
		widget.NewLabelWithStyle("Before Session:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		notifyBeforeCheck,
		beforeValueEntry,
		beforeUnitSelect,
		playSoundBeforeCheck,
		testBeforeButton,
	)

	// --- Combine with a header ---
	header := widget.NewLabelWithStyle(title+" Notifications", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	return container.NewVBox(
		header,
		startRow,
		beforeRow,
		widget.NewSeparator(),
	)
}
