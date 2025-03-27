package main

import (
	_ "embed"

	"f1tray/internal/config"
	"f1tray/internal/models"
	"f1tray/internal/processes"
	"f1tray/internal/ui"
	"f1tray/internal/ui/tabs"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

//go:embed assets/tray_icon.png
var trayIconBytes []byte

// main initializes the F1Tray application, builds the UI, and starts background data refresh routines.
func main() {
	// Load user preferences and build the application state.
	prefs := config.LoadConfig()
	state := models.AppState{
		DebugMode:   prefs.DebugMode,
		Preferences: prefs,
	}

	myApp := app.NewWithID("f1tray")
	myWindow := myApp.NewWindow("F1 Viewer")

	// Create tab content (data will be lazy-loaded)
	scheduleTabData := tabs.CreateScheduleTableTab(models.ScheduleURL, processes.ParseSchedule)
	upcomingTabData := tabs.CreateUpcomingTab(models.UpcomingURL, processes.ParseUpcoming)
	resultsTabData := tabs.CreateResultsTableTab(models.RaceResultsURL, processes.ParseRaceResults)
	qualifyingTabData := tabs.CreateResultsTableTab(models.QualifyingURL, processes.ParseQualifyingResults)
	sprintTabData := tabs.CreateResultsTableTab(models.SprintURL, processes.ParseSprintResults)

	// Create notification overlay using a dedicated UI function.
	notificationLabel, notificationWrapper := ui.CreateNotification()

	// Define a helper function that only takes the silent flag.
	// It captures state, notificationLabel, notificationWrapper, and all tab data.
	refreshData := func(silent bool) {
		go processes.RefreshAllData(state, notificationLabel, notificationWrapper, silent,
			scheduleTabData, upcomingTabData, resultsTabData, qualifyingTabData, sprintTabData)
	}

	// Preferences tab with a save callback that triggers a silent refresh.
	preferencesContent := tabs.CreatePreferencesTab(prefs, func(updated config.Preferences) {
		_ = config.SaveConfig(updated)
		prefs = updated
		state.Preferences = updated
		state.DebugMode = updated.DebugMode

		// Trigger a silent refresh when preferences change.
		refreshData(true)
	})

	// Set up the tabs container.
	tabsContainer := container.NewAppTabs(
		container.NewTabItem("Schedule", scheduleTabData.Content),
		container.NewTabItem("Upcoming", upcomingTabData.Content),
		container.NewTabItem("Race Results", resultsTabData.Content),
		container.NewTabItem("Qualifying", qualifyingTabData.Content),
		container.NewTabItem("Sprint", sprintTabData.Content),
		container.NewTabItem("Preferences", preferencesContent),
	)

	// Stack the tabs with the notification overlay.
	stack := container.NewStack(tabsContainer, notificationWrapper)

	// Create a manual refresh button that uses the refreshData helper.
	refreshButton := widget.NewButton("Refresh All Data", func() {
		refreshData(false)
	})

	myWindow.SetContent(container.NewBorder(refreshButton, nil, nil, nil, stack))
	myWindow.Resize(fyne.NewSize(900, 600))

	// System Tray integration (if supported)
	iconResource := fyne.NewStaticResource("tray_icon.png", trayIconBytes)
	if desk, ok := myApp.(desktop.App); ok {
		processes.SetTrayIcon(desk, iconResource, tabsContainer, myWindow)
	}

	// Determine whether to hide or show the main window on startup.
	if prefs.HideOnOpen {
		myWindow.Hide()
	} else {
		myWindow.Show()
	}

	// Handle window close events.
	myWindow.SetCloseIntercept(func() {
		if prefs.CloseBehavior == "exit" {
			myApp.Quit()
		} else {
			myWindow.Hide()
		}
	})

	// Lazy-load data once the UI is ready.
	refreshData(true)

	// Start background auto-refresh.
	go processes.StartAutoRefresh(state, notificationLabel, notificationWrapper,
		scheduleTabData, upcomingTabData, resultsTabData, qualifyingTabData, sprintTabData)

	myApp.Run()
}
