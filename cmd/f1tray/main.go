package main

import (
	_ "embed"

	"f1tray/internal/config"
	"f1tray/internal/models"
	"f1tray/internal/processes"
	"f1tray/internal/ui/tabs"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

//go:embed assets/tray_icon.png
var trayIconBytes []byte

func main() {
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

	// Notification overlay
	notificationLabel := widget.NewLabel("")
	notificationLabel.Alignment = fyne.TextAlignCenter
	notificationWrapper := container.NewWithoutLayout()

	closeButton := widget.NewButton("✕", func() {
		notificationWrapper.Hide()
	})
	closeButton.Importance = widget.LowImportance

	popup := container.NewPadded(container.NewHBox(
		notificationLabel, layout.NewSpacer(), closeButton,
	))

	popupBG := canvas.NewRectangle(theme.Color(theme.ColorNamePrimary))
	popupBG.SetMinSize(fyne.NewSize(320, 50))
	notificationContainer := container.NewStack(popupBG, popup)
	notificationWrapper = container.NewCenter(notificationContainer)
	notificationWrapper.Hide()

	// Preferences tab with save + refresh callback
	preferencesContent := tabs.CreatePreferencesTab(prefs, func(updated config.Preferences) {
		_ = config.SaveConfig(updated)
		prefs = updated
		state.Preferences = updated
		state.DebugMode = updated.DebugMode

		// Silent refresh for preference update
		go processes.RefreshAllData(state, notificationLabel, notificationWrapper, true,
			scheduleTabData, upcomingTabData, resultsTabData, qualifyingTabData, sprintTabData)
	})

	// Tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("Schedule", scheduleTabData.Content),
		container.NewTabItem("Upcoming", upcomingTabData.Content),
		container.NewTabItem("Race Results", resultsTabData.Content),
		container.NewTabItem("Qualifying", qualifyingTabData.Content),
		container.NewTabItem("Sprint", sprintTabData.Content),
		container.NewTabItem("Preferences", preferencesContent),
	)

	// Layout stack
	stack := container.NewStack(tabs, notificationWrapper)

	// Manual Refresh
	refreshButton := widget.NewButton("Refresh All Data", func() {
		go processes.RefreshAllData(state, notificationLabel, notificationWrapper, false,
			scheduleTabData, upcomingTabData, resultsTabData, qualifyingTabData, sprintTabData)
	})

	myWindow.SetContent(container.NewBorder(refreshButton, nil, nil, nil, stack))
	myWindow.Resize(fyne.NewSize(900, 600))

	// System Tray
	iconResource := fyne.NewStaticResource("tray_icon.png", trayIconBytes)
	if desk, ok := myApp.(desktop.App); ok {
		processes.SetTrayIcon(desk, iconResource, tabs, myWindow)
	}

	// Hide or Show
	if prefs.HideOnOpen {
		myWindow.Hide()
	} else {
		myWindow.Show()
	}

	// Handle X close
	myWindow.SetCloseIntercept(func() {
		if prefs.CloseBehavior == "exit" {
			myApp.Quit()
		} else {
			myWindow.Hide()
		}
	})

	// Lazy-load data after UI is ready
	go processes.RefreshAllData(state, notificationLabel, notificationWrapper, true,
		scheduleTabData, upcomingTabData, resultsTabData, qualifyingTabData, sprintTabData)

	// Background auto-refresh
	go processes.StartAutoRefresh(state, notificationLabel, notificationWrapper,
		scheduleTabData, upcomingTabData, resultsTabData, qualifyingTabData, sprintTabData)

	myApp.Run()
}
