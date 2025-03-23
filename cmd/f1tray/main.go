package main

import (
	"flag"
	"fmt"
	"os"

	"f1tray/internal/notify"
	"f1tray/internal/preferences"
	"f1tray/internal/schedule"
	"f1tray/internal/tray"

	"github.com/getlantern/systray"
)

var debugMode bool

func main() {
	flag.BoolVar(&debugMode, "debug", false, "Enable debug mode to show test options in the tray menu")
	flag.Parse()
	tray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("F1 Tray")
	systray.SetTooltip("F1 Session Notifier")

	// Load tray icon
	iconData, err := os.ReadFile("assets/tray_icon.png")
	if err == nil {
		systray.SetIcon(iconData)
	} else {
		fmt.Println("Failed to load tray icon:", err)
	}

	prefs, err := preferences.LoadPrefs()
	if err != nil {
		fmt.Println("Error loading preferences:", err)
		return
	}

	var mTestNotify, mTestAPI, mTestScheduler, mTestWeeklyReminder *systray.MenuItem

	if debugMode {
		systray.AddSeparator()
		debugLabel := systray.AddMenuItem("— Debug Options —", "")
		debugLabel.Disable()
		mTestNotify = systray.AddMenuItem("Test Notification", "Send a test notification")
		mTestAPI = systray.AddMenuItem("Test API Call", "Get next race weekend info")
		mTestScheduler = systray.AddMenuItem("Test Scheduler", "Trigger race reminder in 10 seconds")
		mTestWeeklyReminder = systray.AddMenuItem("Test Weekly Reminder", "Trigger weekly reminder in 10 seconds")
		systray.AddSeparator()
	}

	mQuit := systray.AddMenuItem("Quit", "Exit the application")

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
				return

			case <-mTestNotify.ClickedCh:
				go func() {
					err := notify.F1Reminder("F1 Tray Test", "This is a test notification!")
					if err != nil {
						fmt.Println("Notification failed:", err)
					}
				}()

			case <-mTestAPI.ClickedCh:
				go schedule.TestRaceNotification()

			case <-mTestScheduler.ClickedCh:
				go schedule.ScheduleNextRaceReminder(true, prefs.RaceReminderHours)

			case <-mTestWeeklyReminder.ClickedCh:
				go schedule.ScheduleWeeklyReminder(true, prefs.WeeklyReminderDay, prefs.WeeklyReminderHour)
			}
		}
	}()

	go schedule.ScheduleNextRaceReminder(false, prefs.RaceReminderHours)
	go schedule.ScheduleWeeklyReminder(false, prefs.WeeklyReminderDay, prefs.WeeklyReminderHour)
}

func onExit() {
	fmt.Println("Exiting F1 Tray.")
	os.Exit(0)
}
