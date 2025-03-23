package main

import (
	"flag"
	"fmt"
	"os"

	"f1tray/internal/gui"
	"f1tray/internal/notify"
	"f1tray/internal/preferences"
	"f1tray/internal/schedule"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var debugMode bool

func main() {
	flag.BoolVar(&debugMode, "debug", false, "Enable debug mode to show test options in the tray menu")
	flag.Parse()

	fmt.Println("Launching F1 Tray with stable Fyne...")

	myApp := app.NewWithID("f1tray")

	prefs, err := preferences.LoadPrefs()
	if err != nil {
		fmt.Println("Error loading preferences:", err)
		os.Exit(1)
	}

	go schedule.ScheduleNextRaceReminder(false, prefs.RaceReminderHours)
	go schedule.ScheduleWeeklyReminder(false, prefs.WeeklyReminderDay, prefs.WeeklyReminderHour)

	// Basic fallback window acting as our menu
	win := myApp.NewWindow("F1 Tray")

	// Core buttons
	buttons := []*widget.Button{
		widget.NewButton("Preferences", func() {
			go gui.ShowPreferencesWindow()
		}),
	}

	if debugMode {
		buttons = append(buttons,
			widget.NewButton("Test Notification", func() {
				go notify.F1Reminder("F1 Tray Test", "This is a test notification!")
			}),
			widget.NewButton("Test API Call", func() {
				go schedule.TestRaceNotification()
			}),
			widget.NewButton("Test Scheduler", func() {
				go schedule.ScheduleNextRaceReminder(true, prefs.RaceReminderHours)
			}),
			widget.NewButton("Test Weekly Reminder", func() {
				go schedule.ScheduleWeeklyReminder(true, prefs.WeeklyReminderDay, prefs.WeeklyReminderHour)
			}),
		)
	}

	buttons = append(buttons, widget.NewButton("Quit", func() {
		myApp.Quit()
	}))

	objects := make([]fyne.CanvasObject, len(buttons))
	for i, b := range buttons {
		objects[i] = b
	}

	win.SetContent(container.NewVBox(objects...))
	win.Resize(fyne.NewSize(300, 200))
	win.ShowAndRun()

}
