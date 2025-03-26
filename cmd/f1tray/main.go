package main

import (
	_ "embed"
	"time"

	"f1tray/internal/config"
	"f1tray/internal/models"
	"f1tray/internal/processes"
	"f1tray/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
)

//go:embed assets/tray_icon.png
var trayIconBytes []byte

func main() {
	myApp := app.NewWithID("f1tray")
	myWindow := myApp.NewWindow("F1 Viewer")

	prefs := config.LoadConfig()

	scheduleTab := ui.CreateScheduleTableTab(models.ScheduleURL, processes.ParseSchedule)
	resultsTab := ui.CreateResultsTableTab(models.RaceResultsURL, processes.ParseRaceResults)
	qualifyingTab := ui.CreateResultsTableTab(models.QualifyingURL, processes.ParseQualifyingResults)
	sprintTab := ui.CreateResultsTableTab(models.SprintURL, processes.ParseSprintResults)
	preferencesTab := ui.CreatePreferencesTab(prefs, func(updated config.Preferences) {
		_ = config.SaveConfig(updated)
		prefs = updated // Update in-memory copy for close behavior
	})

	tabs := container.NewAppTabs(
		container.NewTabItem("Schedule", scheduleTab),
		container.NewTabItem("Race Results", resultsTab),
		container.NewTabItem("Qualifying", qualifyingTab),
		container.NewTabItem("Sprint", sprintTab),
		container.NewTabItem("Preferences", preferencesTab),
	)

	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(900, 600))

	// Use embedded tray icon
	iconResource := fyne.NewStaticResource("tray_icon.png", trayIconBytes)

	if desk, ok := myApp.(desktop.App); ok {
		// Delay to allow tray to become ready (Windows quirk)
		// 100ms is usually enough, adjust if needed
		time.Sleep(500 * time.Millisecond)

		showItem := fyne.NewMenuItem("Show", func() {
			myWindow.Show()
			myWindow.RequestFocus()
		})

		quitItem := fyne.NewMenuItem("Quit", func() {
			myApp.Quit()
		})
		desk.SetSystemTrayIcon(iconResource)
		desk.SetSystemTrayMenu(fyne.NewMenu("F1 Tray", showItem, quitItem))
	}

	// Hide window on startup
	myWindow.Hide()

	// Set behavior for clicking the window X based on config
	myWindow.SetCloseIntercept(func() {
		if prefs.CloseBehavior == "exit" {
			myApp.Quit()
		} else {
			myWindow.Hide()
		}
	})

	myApp.Run()
}
