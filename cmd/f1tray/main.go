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
	prefs := config.LoadConfig()

	myApp := app.NewWithID("f1tray")
	myWindow := myApp.NewWindow("F1 Viewer")

	// Create tab content for each section.
	scheduleContent := ui.CreateScheduleTableTab(models.ScheduleURL, processes.ParseSchedule)
	upcomingContent := ui.CreateUpcomingTab(models.UpcomingURL, processes.ParseUpcoming)
	resultsContent := ui.CreateResultsTableTab(models.RaceResultsURL, processes.ParseRaceResults)
	qualifyingContent := ui.CreateResultsTableTab(models.QualifyingURL, processes.ParseQualifyingResults)
	sprintContent := ui.CreateResultsTableTab(models.SprintURL, processes.ParseSprintResults)
	preferencesContent := ui.CreatePreferencesTab(prefs, func(updated config.Preferences) {
		_ = config.SaveConfig(updated)
		prefs = updated // Update in-memory copy for close behavior
	})

	// Create tab items for each section.
	scheduleTabItem := container.NewTabItem("Schedule", scheduleContent)
	upcomingTabItem := container.NewTabItem("Upcoming", upcomingContent)
	resultsTabItem := container.NewTabItem("Race Results", resultsContent)
	qualifyingTabItem := container.NewTabItem("Qualifying", qualifyingContent)
	sprintTabItem := container.NewTabItem("Sprint", sprintContent)
	preferencesTabItem := container.NewTabItem("Preferences", preferencesContent)

	// Create the AppTabs container.
	tabs := container.NewAppTabs(
		scheduleTabItem,
		upcomingTabItem,
		resultsTabItem,
		qualifyingTabItem,
		sprintTabItem,
		preferencesTabItem,
	)

	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(900, 600))

	// Use embedded tray icon
	iconResource := fyne.NewStaticResource("tray_icon.png", trayIconBytes)

	if desk, ok := myApp.(desktop.App); ok {
		// Delay to allow tray to become ready (Windows quirk)
		time.Sleep(500 * time.Millisecond)

		// Create systray menu items that directly select tabs.
		scheduleItem := fyne.NewMenuItem("Schedule", func() {
			tabs.Select(scheduleTabItem)
			myWindow.Show()
			myWindow.RequestFocus()
		})
		upcomingItem := fyne.NewMenuItem("Upcoming", func() {
			tabs.Select(upcomingTabItem)
			myWindow.Show()
			myWindow.RequestFocus()
		})
		resultsItem := fyne.NewMenuItem("Race Results", func() {
			tabs.Select(resultsTabItem)
			myWindow.Show()
			myWindow.RequestFocus()
		})
		qualifyingItem := fyne.NewMenuItem("Qualifying", func() {
			tabs.Select(qualifyingTabItem)
			myWindow.Show()
			myWindow.RequestFocus()
		})
		sprintItem := fyne.NewMenuItem("Sprint", func() {
			tabs.Select(sprintTabItem)
			myWindow.Show()
			myWindow.RequestFocus()
		})
		preferencesItem := fyne.NewMenuItem("Preferences", func() {
			tabs.Select(preferencesTabItem)
			myWindow.Show()
			myWindow.RequestFocus()
		})
		// A generic "Show" menu item.
		showItem := fyne.NewMenuItem("Show", func() {
			myWindow.Show()
			myWindow.RequestFocus()
		})
		quitItem := fyne.NewMenuItem("Quit", func() {
			myApp.Quit()
		})

		// Create a menu that includes the new items.
		desk.SetSystemTrayIcon(iconResource)
		desk.SetSystemTrayMenu(fyne.NewMenu("F1 Tray",
			scheduleItem,
			upcomingItem,
			resultsItem,
			qualifyingItem,
			sprintItem,
			preferencesItem,
			fyne.NewMenuItemSeparator(),
			showItem,
			quitItem,
		))
	}

	// Hide window on startup.
	myWindow.Hide()

	// Set behavior for clicking the window X based on config.
	myWindow.SetCloseIntercept(func() {
		if prefs.CloseBehavior == "exit" {
			myApp.Quit()
		} else {
			myWindow.Hide()
		}
	})

	myApp.Run()
}
