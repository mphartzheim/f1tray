package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
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

	// Initialize the app using NewWithID()
	myApp := app.NewWithID("f1tray")

	// Load icon from assets folder
	iconData, err := os.ReadFile("assets/tray_icon.png") // or .ico if you prefer
	if err != nil {
		fmt.Println("Error loading icon:", err)
	} else {
		// Convert icon data to a fyne.Resource and set it as the app icon
		iconResource := fyne.NewStaticResource("tray_icon", iconData)
		myApp.SetIcon(iconResource)
	}

	fmt.Println("Launching F1 Tray with stable Fyne...")

	// Load preferences
	prefs, err := preferences.LoadPrefs()
	if err != nil {
		fmt.Println("Error loading preferences:", err)
		os.Exit(1)
	}

	// Create a quit channel to manage background processes
	quitChannel := make(chan struct{})
	var wg sync.WaitGroup

	// Start scheduled reminders (background tasks)
	wg.Add(1) // Increment WaitGroup counter
	go func() {
		defer wg.Done() // Decrement WaitGroup counter when goroutine completes
		for {
			select {
			case <-quitChannel:
				// Exit the goroutine when quit signal is received
				return
			default:
				// Continue with scheduling tasks
				schedule.ScheduleNextRaceReminder(false, prefs.RaceReminderHours)
				schedule.ScheduleWeeklyReminder(false, prefs.WeeklyReminderDay, prefs.WeeklyReminderTime) // Use WeeklyReminderTime
				time.Sleep(10 * time.Second)                                                              // Sleep to avoid blocking the loop
			}
		}
	}()

	// Basic fallback window acting as our menu
	win := myApp.NewWindow("F1 Tray")

	// Core buttons
	buttons := []*widget.Button{
		widget.NewButton("Preferences", func() {
			go gui.ShowPreferencesWindow(myApp) // Pass myApp to preferences window
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
				go schedule.ScheduleWeeklyReminder(true, prefs.WeeklyReminderDay, prefs.WeeklyReminderTime) // Use WeeklyReminderTime
			}),
		)
	}

	// Quit button logic to stop background tasks and exit the app
	buttons = append(buttons, widget.NewButton("Quit", func() {
		// Signal the background tasks to stop immediately
		close(quitChannel)
		wg.Wait() // Wait for the background task to finish
		fmt.Println("Exiting F1 Tray.")
		myApp.Quit()
	}))

	objects := make([]fyne.CanvasObject, len(buttons))
	for i, b := range buttons {
		objects[i] = b
	}

	// Set window content and size
	win.SetContent(container.NewVBox(objects...))
	win.Resize(fyne.NewSize(300, 200))

	// Ensure the window is displayed
	win.Show()

	// Run the app
	myApp.Run()
}
