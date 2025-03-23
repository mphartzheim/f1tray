package gui

import (
	"f1tray/internal/preferences"
	"f1tray/internal/schedule"
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var prefsApp fyne.App
var window fyne.Window

func ShowPreferencesWindow() {
	if prefsApp == nil {
		prefsApp = app.NewWithID("f1tray-preferences")
		window = prefsApp.NewWindow("F1 Tray Preferences")
		window.Resize(fyne.NewSize(400, 240))
		buildPreferencesUI()
	} else {
		window.Show()
		window.RequestFocus()
	}
}

func buildPreferencesUI() {
	current, err := preferences.LoadPrefs()
	if err != nil {
		fmt.Println("Failed to load preferences:", err)
		current = preferences.UserPrefs{
			RaceReminderHours:  2,
			WeeklyReminderDay:  "Wednesday",
			WeeklyReminderHour: 12, // ✅ now an int, not a string
		}
	}

	// Widgets
	hoursEntry := widget.NewEntry()
	hoursEntry.SetText(strconv.Itoa(current.RaceReminderHours))

	days := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	daySelect := widget.NewSelect(days, nil)
	daySelect.SetSelected(current.WeeklyReminderDay)

	hourOptions := make([]string, 24)
	for i := 0; i < 24; i++ {
		hourOptions[i] = fmt.Sprintf("%02d:00", i)
	}
	hourSelect := widget.NewSelect(hourOptions, nil)
	hourSelect.SetSelected(fmt.Sprintf("%02d:00", current.WeeklyReminderHour)) // ✅ convert int to string

	// Buttons
	saveBtn := widget.NewButton("Save", func() {
		hours, err := strconv.Atoi(hoursEntry.Text)
		if err != nil || hours < 1 || hours > 48 {
			dialog.ShowError(fmt.Errorf("Reminder hours must be a number between 1 and 48"), window)
			return
		}

		// Parse "14:00" to 14
		hourStr := hourSelect.Selected
		hourInt, err := strconv.Atoi(strings.Split(hourStr, ":")[0])
		if err != nil {
			dialog.ShowError(fmt.Errorf("Invalid hour selected"), window)
			return
		}

		newPrefs := preferences.UserPrefs{
			RaceReminderHours:  hours,
			WeeklyReminderDay:  daySelect.Selected,
			WeeklyReminderHour: hourInt, // ✅ now an int
		}

		err = preferences.SavePrefs(newPrefs)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Failed to save preferences: %v", err), window)
			return
		}

		go schedule.ScheduleNextRaceReminder(false, newPrefs.RaceReminderHours)
		go schedule.ScheduleWeeklyReminder(false, newPrefs.WeeklyReminderDay, newPrefs.WeeklyReminderHour)

		dialog.ShowInformation("Saved", "Preferences saved successfully.", window)
		window.Hide()
	})

	closeBtn := widget.NewButton("Close", func() {
		window.Hide()
	})

	form := container.NewVBox(
		widget.NewLabel("Remind me X hours before each session:"),
		hoursEntry,
		widget.NewLabel("Weekly reminder day:"),
		daySelect,
		widget.NewLabel("Weekly reminder time:"),
		hourSelect,
		container.NewHBox(saveBtn, closeBtn),
	)

	window.SetContent(form)
	window.Show()
}
