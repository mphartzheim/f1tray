package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"time"

	"f1tray/internal/config"
	"f1tray/internal/models"
	"f1tray/internal/processes"
	"f1tray/internal/ui"
	"f1tray/internal/ui/tabs"
	"f1tray/internal/ui/themes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

//go:embed assets/tray_icon.png
var trayIconBytes []byte

func main() {
	state := models.AppState{
		FirstRun: true,
	}

	// Create the Fyne app.
	myApp := app.NewWithID("f1tray")

	// Choose the initial theme based on preferences.
	var initialTheme fyne.Theme
	switch config.Get().Themes.Theme {
	case "Light":
		initialTheme = themes.LightTheme{}
	case "Dark":
		initialTheme = themes.DarkTheme{}
	default:
		initialTheme = themes.SystemTheme{}
	}
	myApp.Settings().SetTheme(initialTheme)

	if config.Get().Debug.Enabled {
		fmt.Printf("Theme: %T\n", theme.Current())
	}

	myWindow := myApp.NewWindow("F1 Viewer")

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
		container.NewTabItem("Preferences", tabs.CreatePreferencesTab(func(updated config.Preferences) {
			_ = config.SaveConfig(updated)
		}, func() {
			upcomingTabData.Refresh()
		})),
	)

	// Hook up the UpdateTabs callback so that processes.ReloadOtherTabs can update the three tabs.
	processes.UpdateTabs = func(resultsContent, qualifyingContent, sprintContent fyne.CanvasObject) {
		tabsContainer.Items[2].Content = resultsContent
		tabsContainer.Items[3].Content = qualifyingContent
		tabsContainer.Items[4].Content = sprintContent

		// Mark updated tabs with an asterisk if they are not currently selected.
		if tabsContainer.Selected() != tabsContainer.Items[2] {
			if len(tabsContainer.Items[2].Text) == 0 || tabsContainer.Items[2].Text[len(tabsContainer.Items[2].Text)-1] != '*' {
				tabsContainer.Items[2].Text += "*"
			}
		}
		if tabsContainer.Selected() != tabsContainer.Items[3] {
			if len(tabsContainer.Items[3].Text) == 0 || tabsContainer.Items[3].Text[len(tabsContainer.Items[3].Text)-1] != '*' {
				tabsContainer.Items[3].Text += "*"
			}
		}
		if tabsContainer.Selected() != tabsContainer.Items[4] {
			if len(tabsContainer.Items[4].Text) == 0 || tabsContainer.Items[4].Text[len(tabsContainer.Items[4].Text)-1] != '*' {
				tabsContainer.Items[4].Text += "*"
			}
		}
		tabsContainer.Refresh()
	}

	// Remove a trailing asterisk when a tab is selected.
	tabsContainer.OnSelected = func(selectedTab *container.TabItem) {
		if len(selectedTab.Text) > 0 && selectedTab.Text[len(selectedTab.Text)-1] == '*' {
			selectedTab.Text = selectedTab.Text[:len(selectedTab.Text)-1]
			tabsContainer.Refresh()
		}
	}

	// When the selected year changes, update the Schedule tab's content.
	yearSelect.OnChanged = func(selectedYear string) {
		newScheduleTabData := tabs.CreateScheduleTableTab(processes.ParseSchedule, selectedYear)
		scheduleTab.Content = newScheduleTabData.Content
		tabsContainer.Refresh()
	}

	// Create notification overlay using your dedicated UI function.
	notificationLabel, notificationWrapper := ui.CreateNotification()

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
	if config.Get().Window.HideOnOpen {
		myWindow.Hide()
	} else {
		myWindow.Show()
	}

	// Handle window close events.
	myWindow.SetCloseIntercept(func() {
		if config.Get().Window.CloseBehavior == "exit" {
			myApp.Quit()
		} else {
			myWindow.Hide()
		}
	})

	// Lazy-load data once the UI is ready.
	go processes.RefreshAllData(notificationLabel, notificationWrapper,
		upcomingTabData, resultsTabData, qualifyingTabData, sprintTabData)

	// Start background auto-refresh.
	go processes.StartAutoRefresh(&state, fmt.Sprintf("%d", time.Now().Year()))

	myApp.Run()
}
