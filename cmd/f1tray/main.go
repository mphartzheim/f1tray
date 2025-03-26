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
	"fyne.io/fyne/v2/widget"
)

//go:embed assets/tray_icon.png
var trayIconBytes []byte

func main() {
	prefs := config.LoadConfig()

	myApp := app.NewWithID("f1tray")
	myWindow := myApp.NewWindow("F1 Viewer")

	// Create tab content for each section using TabData.
	scheduleTabData := ui.CreateScheduleTableTab(models.ScheduleURL, processes.ParseSchedule)
	upcomingTabData := ui.CreateUpcomingTab(models.UpcomingURL, processes.ParseUpcoming)
	resultsTabData := ui.CreateResultsTableTab(models.RaceResultsURL, processes.ParseRaceResults)
	qualifyingTabData := ui.CreateResultsTableTab(models.QualifyingURL, processes.ParseQualifyingResults)
	sprintTabData := ui.CreateResultsTableTab(models.SprintURL, processes.ParseSprintResults)
	preferencesContent := ui.CreatePreferencesTab(prefs, func(updated config.Preferences) {
		_ = config.SaveConfig(updated)
		prefs = updated // Update in-memory copy for close behavior
	})

	// Create tab items using the Content field from TabData.
	scheduleTabItem := container.NewTabItem("Schedule", scheduleTabData.Content)
	upcomingTabItem := container.NewTabItem("Upcoming", upcomingTabData.Content)
	resultsTabItem := container.NewTabItem("Race Results", resultsTabData.Content)
	qualifyingTabItem := container.NewTabItem("Qualifying", qualifyingTabData.Content)
	sprintTabItem := container.NewTabItem("Sprint", sprintTabData.Content)
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

	refreshButton := widget.NewButton("Refresh All Data", func() {
		// Call refresh functions for all tabs.
		scheduleTabData.Refresh()
		upcomingTabData.Refresh()
		resultsTabData.Refresh()
		qualifyingTabData.Refresh()
		sprintTabData.Refresh()

		processes.ShowInAppNotification(myWindow, "Data Refresh", "All data has been refreshed.")
	})

	// Wrap the refresh button above the tabs.
	content := container.NewBorder(refreshButton, nil, nil, nil, tabs)
	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(900, 600))

	// Use embedded tray icon.
	iconResource := fyne.NewStaticResource("tray_icon.png", trayIconBytes)
	if desk, ok := myApp.(desktop.App); ok {
		// Delay to allow tray to become ready (Windows quirk).
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
		showItem := fyne.NewMenuItem("Show", func() {
			myWindow.Show()
			myWindow.RequestFocus()
		})
		quitItem := fyne.NewMenuItem("Quit", func() {
			myApp.Quit()
		})

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
