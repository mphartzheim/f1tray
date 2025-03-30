package preferences

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/mphartzheim/f1tray/internal/config"
	"github.com/mphartzheim/f1tray/internal/models"
	"github.com/mphartzheim/f1tray/internal/notifications"
	"github.com/mphartzheim/f1tray/internal/ui/themes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// setWidgetsEnabled sets the enabled state for all widgets passed in.
func setWidgetsEnabled(enabled bool, widgets ...interface{}) {
	for _, w := range widgets {
		// Check if the widget implements both Enable and Disable.
		if d, ok := w.(interface {
			Enable()
			Disable()
		}); ok {
			if enabled {
				d.Enable()
			} else {
				d.Disable()
			}
		}
	}
}

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
	playSoundStartCheck := widget.NewCheck("Play sound", func(checked bool) {
		sessionPrefs.PlaySoundOnStart = checked
		_ = config.Set(config.Get())
		onSave(*config.Get())
	})
	playSoundStartCheck.SetChecked(sessionPrefs.PlaySoundOnStart)

	testStartButton := widget.NewButton("Test", func() {
		fmt.Printf("Testing %s session start notification\n", title)
		// Simulate a dummy session event for "At Session Start"
		dummySession := models.SessionInfo{
			Type:      title,
			StartTime: time.Now().Add(-30 * time.Second), // Assume session started 30 seconds ago
			Label:     fmt.Sprintf("%s - Test Start", title),
		}
		notifications.CheckAndSendNotifications(dummySession)
	})

	notifyStartCheck := widget.NewCheck("Notify at session start", func(checked bool) {
		sessionPrefs.NotifyOnStart = checked
		setWidgetsEnabled(checked, playSoundStartCheck, testStartButton)
		_ = config.Set(config.Get())
		onSave(*config.Get())
	})
	notifyStartCheck.SetChecked(sessionPrefs.NotifyOnStart)

	// Set initial enable/disable state.
	setWidgetsEnabled(sessionPrefs.NotifyOnStart, playSoundStartCheck, testStartButton)

	startRow := container.NewHBox(
		widget.NewLabelWithStyle("At Session Start:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		notifyStartCheck,
		playSoundStartCheck,
		testStartButton,
	)

	// --- Row for "Before Session" ---
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
		// Simulate a dummy session event for "Before Session"
		dummySession := models.SessionInfo{
			Type:      title,
			StartTime: time.Now().Add(10 * time.Minute), // Assume session starts in 10 minutes
			Label:     fmt.Sprintf("%s - Test Before", title),
		}
		notifications.CheckAndSendNotifications(dummySession)
	})

	notifyBeforeCheck := widget.NewCheck("Notify before session", func(checked bool) {
		sessionPrefs.NotifyBefore = checked
		setWidgetsEnabled(checked, beforeValueEntry, beforeUnitSelect, playSoundBeforeCheck, testBeforeButton)
		_ = config.Set(config.Get())
		onSave(*config.Get())
	})
	notifyBeforeCheck.SetChecked(sessionPrefs.NotifyBefore)

	// Set initial enable/disable state.
	setWidgetsEnabled(sessionPrefs.NotifyBefore, beforeValueEntry, beforeUnitSelect, playSoundBeforeCheck, testBeforeButton)

	beforeRow := container.NewHBox(
		widget.NewLabelWithStyle("Before Session:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		notifyBeforeCheck,
		beforeValueEntry,
		beforeUnitSelect,
		playSoundBeforeCheck,
		testBeforeButton,
	)

	header := widget.NewLabelWithStyle(title+" Notifications", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	return container.NewVBox(
		header,
		startRow,
		beforeRow,
		widget.NewSeparator(),
	)
}
