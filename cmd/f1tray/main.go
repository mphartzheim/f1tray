package main

import (
	_ "embed"
	"strconv"
	"time"

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

func main() {
	// Create the Fyne app and window.
	myApp := app.NewWithID("f1tray")
	myWindow := myApp.NewWindow("F1 Viewer")

	// Load user preferences and create the application state.
	prefs := config.LoadConfig()
	state := models.AppState{
		DebugMode:   prefs.DebugMode,
		Preferences: prefs,
		FirstRun:    true,
	}

	// Build a slice of years (as strings) from the current year down to 1950.
	currentYear := time.Now().Year()
	years := []string{}
	for y := currentYear; y >= 1950; y-- {
		years = append(years, strconv.Itoa(y))
	}

	// Create the drop-down widget for year selection.
	yearSelect := widget.NewSelect(years, nil)
	yearSelect.SetSelected(years[0]) // Default to the current year

	// Create a header container that now only includes the schedule selector.
	headerContainer := container.NewHBox(widget.NewLabel("Season"), yearSelect)

	// Create initial schedule table content using the selected year.
	scheduleTabData := tabs.CreateScheduleTableTab(processes.ParseSchedule, yearSelect.Selected)
	scheduleTab := container.NewTabItem("Schedule", scheduleTabData.Content)

	// Create the rest of your tabs using the default year.
	upcomingTabData := tabs.CreateUpcomingTab(processes.ParseUpcoming, yearSelect.Selected)
	resultsTabData := tabs.CreateResultsTableTab(processes.ParseRaceResults, yearSelect.Selected, "last")
	qualifyingTabData := tabs.CreateResultsTableTab(processes.ParseQualifyingResults, yearSelect.Selected, "last")
	sprintTabData := tabs.CreateResultsTableTab(processes.ParseSprintResults, yearSelect.Selected, "last")

	// Create the tabs container.
	tabsContainer := container.NewAppTabs(
		scheduleTab,
		container.NewTabItem("Upcoming", upcomingTabData.Content),
		container.NewTabItem("Race Results", resultsTabData.Content),
		container.NewTabItem("Qualifying", qualifyingTabData.Content),
		container.NewTabItem("Sprint", sprintTabData.Content),
		container.NewTabItem("Preferences", tabs.CreatePreferencesTab(prefs, func(updated config.Preferences) {
			_ = config.SaveConfig(updated)
			prefs = updated
			state.Preferences = updated
			state.DebugMode = updated.DebugMode
		})),
	)

	// Hook up the UpdateTabs callback so that processes.ReloadOtherTabs can update the three tabs.
	processes.UpdateTabs = func(resultsContent, qualifyingContent, sprintContent fyne.CanvasObject) {
		// tabsContainer.Items[2] = Race Results, [3] = Qualifying, [4] = Sprint.
		tabsContainer.Items[2].Content = resultsContent
		tabsContainer.Items[3].Content = qualifyingContent
		tabsContainer.Items[4].Content = sprintContent
		tabsContainer.Refresh()
	}

	// When the selected year changes, update the Schedule tab's content.
	yearSelect.OnChanged = func(selectedYear string) {
		newScheduleTabData := tabs.CreateScheduleTableTab(processes.ParseSchedule, selectedYear)
		// Update the content field of our schedule tab.
		scheduleTab.Content = newScheduleTabData.Content
		// Refresh the tab container to show the updated content.
		tabsContainer.Refresh()
	}

	// Create notification overlay using your dedicated UI function.
	notificationLabel, notificationWrapper := ui.CreateNotification()

	// Define a helper function that refreshes all data.
	refreshData := func(silent bool) {
		go processes.RefreshAllData(&state, notificationLabel, notificationWrapper, silent,
			upcomingTabData, resultsTabData, qualifyingTabData, sprintTabData)
	}

	// Stack the tabs with the notification overlay.
	stack := container.NewStack(tabsContainer, notificationWrapper)

	// Use the header container (with the schedule selector) as the top border.
	myWindow.SetContent(container.NewBorder(headerContainer, nil, nil, nil, stack))
	myWindow.Resize(fyne.NewSize(900, 600))

	// System Tray integration (if supported).
	iconResource := fyne.NewStaticResource("tray_icon.png", trayIconBytes)
	if desk, ok := myApp.(desktop.App); ok {
		processes.SetTrayIcon(desk, iconResource, tabsContainer, myWindow)
	}

	// Show or hide the window based on user preferences.
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
	go processes.StartAutoRefresh(&state, notificationLabel, notificationWrapper,
		upcomingTabData, resultsTabData, qualifyingTabData, sprintTabData)

	myApp.Run()
}
