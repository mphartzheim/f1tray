package gui

import (
	"f1tray/internal/preferences"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ShowPreferencesWindow handles the UI for preferences management
func ShowPreferencesWindow(myApp fyne.App) {
	// Create a window for preferences
	prefsWindow := myApp.NewWindow("Preferences")

	// Load current preferences to display
	prefs, err := preferences.LoadPrefs()
	if err != nil {
		fmt.Println("Error loading preferences:", err)
	}

	// Create UI elements for Preferences (like race reminder hours, weekly reminders)
	raceReminderHoursEntry := widget.NewEntry()
	raceReminderHoursEntry.SetText(fmt.Sprintf("%d", prefs.RaceReminderHours))

	weeklyReminderDaySelect := widget.NewSelect([]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, func(value string) {
		// No action needed for now
	})

	// Populate fields with current preferences
	weeklyReminderDaySelect.SetSelected(prefs.WeeklyReminderDay)

	// Time selector for Weekly Reminder Time
	weeklyReminderTimePicker := widget.NewEntry()
	weeklyReminderTimePicker.SetText(prefs.WeeklyReminderTime.Format("15:04"))

	// Button to save preferences
	saveButton := widget.NewButton("Save", func() {
		// Save the preferences to file
		prefs.RaceReminderHours = prefs.RaceReminderHours
		prefs.WeeklyReminderDay = weeklyReminderDaySelect.Selected

		// Parse the time string to save the weekly reminder time
		reminderTime, err := time.Parse("15:04", weeklyReminderTimePicker.Text)
		if err != nil {
			fmt.Println("Error parsing Weekly Reminder Time:", err)
		} else {
			prefs.WeeklyReminderTime = reminderTime
		}

		if err := preferences.SavePrefs(prefs); err != nil {
			fmt.Println("Error saving preferences:", err)
		}
		prefsWindow.Close()
	})

	// Layout the widgets
	content := container.NewVBox(
		widget.NewLabel("Race Reminder Hours:"),
		raceReminderHoursEntry,
		widget.NewLabel("Weekly Reminder Day:"),
		weeklyReminderDaySelect,
		widget.NewLabel("Weekly Reminder Time:"),
		weeklyReminderTimePicker,
		saveButton,
	)

	prefsWindow.SetContent(content)
	prefsWindow.Show()
}
