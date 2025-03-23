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
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var debugMode bool

func main() {
	flag.BoolVar(&debugMode, "debug", false, "Enable debug mode to show test options in the tray menu")
	flag.Parse()

	// Initialize app with tray support
	myApp := app.NewWithTray("f1tray")
	myApp.SetIcon(theme.ComputerIcon()) // Optional: Replace with your own tray icon

	fmt.Println("Starting F1 Tray...")

	// Load preferences
	prefs, err := preferences.LoadPrefs()
	if err != nil {
		fmt.Println("Error loading preferences:", err)
		os.Exit(1)
	}

	// Start scheduled reminders
	go schedule.ScheduleNextRaceReminder(false, prefs.RaceReminderHours)
	go schedule.ScheduleWeeklyReminder(false, prefs.WeeklyReminderDay, prefs.WeeklyReminderHour)

	// Build system tray menu
	menuItems := []*fyne.MenuItem{
		fyne.NewMenuItem("Preferences", func() {
			fmt.Println("Opening Preferences")
			go gui.ShowPreferencesWindow()
		}),
	}

	if debugMode {
		menuItems = append(menuItems,
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Test Notification", func() {
				go notify.F1Reminder("F1 Tray Test", "This is a test notification!")
			}),
			fyne.NewMenuItem("Test API Call", func() {
				go schedule.TestRaceNotification()
			}),
			fyne.NewMenuItem("Test Scheduler", func() {
				go schedule.ScheduleNextRaceReminder(true, prefs.RaceReminderHours)
			}),
			fyne.NewMenuItem("Test Weekly Reminder", func() {
				go schedule.ScheduleWeeklyReminder(true, prefs.WeeklyReminderDay, prefs.WeeklyReminderHour)
			}),
		)
	}

	menuItems = append(menuItems,
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quit", func() {
			fmt.Println("Exiting F1 Tray.")
			myApp.Quit()
		}),
	)

	trayMenu := fyne.NewMenu("F1 Tray", menuItems...)
	myApp.SetSystemTrayMenu(trayMenu)

	// Set the system tray icon (ensure you have the correct image data)
	iconData, _ := os.ReadFile("assets/tray_icon.png")
	myApp.SetSystemTrayIcon(iconData)

	fmt.Println("Tray menu set ✅")

	// Create a hidden window to keep the app running
	win := myApp.NewWindow("F1 Tray")
	win.SetContent(container.NewVBox(
		widget.NewLabel("F1 Tray is running in the system tray."),
	))
	win.Hide()

	myApp.Run()
}
