package main

import (
	"fmt"
	"io"
	"os"

	"f1tray/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
)

func main() {
	myApp := app.NewWithID("f1tray")
	myWindow := myApp.NewWindow("F1 Viewer")

	resultsTab := ui.CreateResultsTableTab("https://api.jolpi.ca/ergast/f1/current/last/results.json", ui.ParseRaceResults)
	qualifyingTab := ui.CreateResultsTableTab("https://api.jolpi.ca/ergast/f1/current/last/qualifying.json", ui.ParseQualifyingResults)
	sprintTab := ui.CreateResultsTableTab("https://api.jolpi.ca/ergast/f1/current/last/sprint.json", ui.ParseSprintResults)
	scheduleTab := ui.CreateScheduleTableTab("https://api.jolpi.ca/ergast/f1/current.json", ui.ParseSchedule)

	tabs := container.NewAppTabs(
		container.NewTabItem("Race Results", resultsTab),
		container.NewTabItem("Qualifying", qualifyingTab),
		container.NewTabItem("Sprint", sprintTab),
		container.NewTabItem("Schedule", scheduleTab),
	)

	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(900, 600))

	// Load tray icon
	iconFileURI := storage.NewFileURI("assets/tray_icon.png")
	iconReader, err := storage.Reader(iconFileURI)
	if err != nil {
		fmt.Println("Failed to load tray icon:", err)
		os.Exit(1)
	}
	iconBytes, err := io.ReadAll(iconReader)
	if err != nil {
		fmt.Println("Failed to read icon file:", err)
		os.Exit(1)
	}
	iconResource := fyne.NewStaticResource("tray_icon.png", iconBytes)

	// Track hidden state for restoration
	isClosed := false

	if desk, ok := myApp.(desktop.App); ok {
		showItem := fyne.NewMenuItem("Show", func() {
			if isClosed {
				myWindow = myApp.NewWindow("F1 Viewer")
				myWindow.SetContent(tabs)
				myWindow.Resize(fyne.NewSize(900, 600))
				myWindow.SetOnClosed(func() {
					myWindow.Hide()
					isClosed = true
				})
				myWindow.Show()
				isClosed = false
			} else {
				myWindow.Show()
			}
		})
		quitItem := fyne.NewMenuItem("Quit", func() {
			myApp.Quit()
		})
		desk.SetSystemTrayIcon(iconResource)
		desk.SetSystemTrayMenu(fyne.NewMenu("F1 Tray", showItem, quitItem))
	}

	// Hide window on startup
	myWindow.Hide()

	// Hide window when closed instead of quitting
	myWindow.SetOnClosed(func() {
		myWindow.Hide()
		isClosed = true
	})

	myApp.Run()
}
