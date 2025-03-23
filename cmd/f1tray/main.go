package main

import (
	"flag"
	"fmt"
	"os"

	"f1tray/internal/notify"
	"f1tray/internal/schedule"
	"f1tray/internal/tray"

	"github.com/getlantern/systray"
)

var testReminder bool

func main() {
	flag.BoolVar(&testReminder, "test-reminder", false, "Trigger the race reminder after 10 seconds for testing")
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

	// Menu items
	mTestNotify := systray.AddMenuItem("Test Notification", "Send a test notification")
	mTestAPI := systray.AddMenuItem("Test API Call", "Get next race weekend info")
	systray.AddSeparator()
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
			}
		}
	}()

	go schedule.ScheduleNextRaceReminder(testReminder)
}

func onExit() {
	fmt.Println("Exiting F1 Tray.")
	os.Exit(0)
}
