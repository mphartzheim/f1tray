package main

import (
	"flag"
	"fmt"
	"os"
	"time"

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

	// Create a quit channel to manage background processes
	quitChannel := make(chan struct{})

	// Start scheduled reminders (background tasks)
	go func() {
		for {
			select {
			case <-quitChannel:
				// Stop the scheduled tasks when quit signal is received
				fmt.Println("Stopping scheduled tasks...")
				return
			default:
				// Continue with scheduling tasks
				schedule.ScheduleNextRaceReminder(false, prefs.RaceReminderHours)
				schedule.ScheduleWeeklyReminder(false, prefs.WeeklyReminderDay, prefs.WeeklyReminderHour)
				time.Sleep(10 * time.Second) // Sleep to avoid blocking the loop
			}
		}
	}()

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

	// Quit button logic to stop background tasks and exit the app
	buttons = append(buttons, widget.NewButton("Quit", func() {
		// Signal the background tasks to stop
		close(quitChannel)
		fmt.Println("Exiting F1 Tray.")
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
