package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"time"

	"github.com/mphartzheim/f1tray/internal/config"
	"github.com/mphartzheim/f1tray/internal/models"
	"github.com/mphartzheim/f1tray/internal/processes"
	"github.com/mphartzheim/f1tray/internal/ui"
	"github.com/mphartzheim/f1tray/internal/ui/tabs"
	"github.com/mphartzheim/f1tray/internal/ui/tabs/preferences"
	"github.com/mphartzheim/f1tray/internal/ui/tabs/results"
	"github.com/mphartzheim/f1tray/internal/ui/tabs/standings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	fyneTheme "fyne.io/fyne/v2/theme"
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
	initialTheme := processes.GetThemeFromName(config.Get().Themes.Theme)
	myApp.Settings().SetTheme(initialTheme)

	if config.Get().Debug.Enabled {
		fmt.Printf("Theme: %T\n", fyneTheme.Current())
	}

	myWindow := myApp.NewWindow("F1 Viewer")
	models.MainWindow = myWindow

	// Build a slice of years (as strings) from the current year down to 1950.
	currentYear := time.Now().Year()
	years := []string{}
	for y := currentYear; y >= 1950; y-- {
		years = append(years, strconv.Itoa(y))
	}

	// Create the drop-down widget for year selection.
	yearSelect := widget.NewSelect(years, nil)
	yearSelect.SetSelected(years[0]) // Default to the current year

	// Create a header container that includes the season selector.
	headerContainer := container.NewHBox(widget.NewLabel("Season"), yearSelect)

	// Create initial Schedule tab.
	scheduleTabData := tabs.CreateScheduleTableTab(processes.ParseSchedule, yearSelect.Selected)
	scheduleTab := container.NewTabItem("Schedule", scheduleTabData.Content)

	// Create Upcoming tab.
	upcomingTabData := tabs.CreateUpcomingTab(&state, processes.ParseUpcoming, yearSelect.Selected)
	upcomingTab := container.NewTabItem("Upcoming", upcomingTabData.Content)

	// Create Results tab.
	resultsTabData, resultsInnerTabs := results.CreateResultsTab(yearSelect.Selected, "last")
	resultsOuterTab := container.NewTabItem("Results", resultsTabData.Content)

	// Create Standings tab.
	standingsTabData, standingsInnerTabs := standings.CreateStandingsTab(yearSelect.Selected, yearSelect.Selected)
	standingsOuterTab := container.NewTabItem("Standings", standingsTabData.Content)

	// Create Preferences tab.
	preferencesTab := container.NewTabItem("Preferences", preferences.CreatePreferencesTab(
		func(updated config.Preferences) {
			_ = config.SaveConfig(updated)
		},
		func() {
			upcomingTabData.Refresh()
		},
	))

	// Create the AppTabs container with all tabs.
	tabsContainer := container.NewAppTabs(
		scheduleTab,
		upcomingTab,
		resultsOuterTab,
		standingsOuterTab,
		preferencesTab,
	)

	// Hook up callbacks for updating inner tab content.
	processes.UpdateTabs = func(resultsContent, qualifyingContent, sprintContent fyne.CanvasObject) {
		resultsInnerTabs.Items[0].Content = resultsContent
		resultsInnerTabs.Items[1].Content = qualifyingContent
		resultsInnerTabs.Items[2].Content = sprintContent
		resultsInnerTabs.Refresh()

		if tabsContainer.SelectedIndex() != 2 {
			if len(resultsOuterTab.Text) == 0 || resultsOuterTab.Text[len(resultsOuterTab.Text)-1] != '*' {
				resultsOuterTab.Text += "*"
			}
		}
		tabsContainer.Refresh()
	}
	processes.UpdateStandingsTabs = func(driversContent, constructorsContent fyne.CanvasObject) {
		standingsInnerTabs.Items[0].Content = driversContent
		standingsInnerTabs.Items[1].Content = constructorsContent
		standingsInnerTabs.Refresh()

		if tabsContainer.SelectedIndex() != 3 {
			if len(standingsOuterTab.Text) == 0 || standingsOuterTab.Text[len(standingsOuterTab.Text)-1] != '*' {
				standingsOuterTab.Text += "*"
			}
			tabsContainer.Refresh()
		}
	}

	// Remove a trailing asterisk when a tab is selected.
	tabsContainer.OnSelected = func(selectedTab *container.TabItem) {
		if len(selectedTab.Text) > 0 && selectedTab.Text[len(selectedTab.Text)-1] == '*' {
			selectedTab.Text = selectedTab.Text[:len(selectedTab.Text)-1]
			tabsContainer.Refresh()
		}
	}
	resultsInnerTabs.OnSelected = func(selectedTab *container.TabItem) {
		if len(selectedTab.Text) > 0 && selectedTab.Text[len(selectedTab.Text)-1] == '*' {
			selectedTab.Text = selectedTab.Text[:len(selectedTab.Text)-1]
			resultsInnerTabs.Refresh()
		}
	}
	standingsInnerTabs.OnSelected = func(selectedTab *container.TabItem) {
		if len(selectedTab.Text) > 0 && selectedTab.Text[len(selectedTab.Text)-1] == '*' {
			selectedTab.Text = selectedTab.Text[:len(selectedTab.Text)-1]
			standingsInnerTabs.Refresh()
		}
	}

	// Update content when the selected year changes.
	yearSelect.OnChanged = func(selectedYear string) {
		newScheduleTabData := tabs.CreateScheduleTableTab(processes.ParseSchedule, selectedYear)
		scheduleTab.Content = newScheduleTabData.Content

		selectedStandingsIndex := standingsInnerTabs.SelectedIndex()
		newStandingsTabData, newStandingsInnerTabs := standings.CreateStandingsTab(selectedYear, "last")
		standingsOuterTab.Content = newStandingsTabData.Content
		standingsInnerTabs = newStandingsInnerTabs

		if selectedStandingsIndex >= 0 && selectedStandingsIndex < len(standingsInnerTabs.Items) {
			standingsInnerTabs.SelectIndex(selectedStandingsIndex)
		}

		if tabsContainer.SelectedIndex() != 3 {
			if len(standingsOuterTab.Text) == 0 || standingsOuterTab.Text[len(standingsOuterTab.Text)-1] != '*' {
				standingsOuterTab.Text += "*"
			}
		}

		tabsContainer.Refresh()
	}

	// Create a notification overlay.
	notificationLabel, notificationWrapper := ui.CreateNotification(myWindow)

	// Stack the tabs and the notification overlay.
	stack := container.NewStack(tabsContainer, notificationWrapper)

	// Set the window content.
	myWindow.SetContent(container.NewBorder(headerContainer, nil, nil, nil, stack))
	myWindow.Resize(fyne.NewSize(900, 600))

	// System tray integration.
	iconResource := fyne.NewStaticResource("tray_icon.png", trayIconBytes)
	if desk, ok := myApp.(desktop.App); ok {
		// Pass the full TabItems directly.
		processes.SetTrayIcon(desk, iconResource, tabsContainer, myWindow,
			scheduleTab, upcomingTab, resultsOuterTab, standingsOuterTab, preferencesTab,
		)
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
	go processes.RefreshAllData(notificationLabel, notificationWrapper, upcomingTabData, resultsTabData)

	// Start background auto-refresh.
	go processes.StartAutoRefresh(&state, fmt.Sprintf("%d", time.Now().Year()))

	// Start a ticker to refresh the Upcoming tab if it's a session day.
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			isUpcomingTabVisible := tabsContainer.Selected() == upcomingTab
			isSessionToday := processes.IsSessionDay(state.UpcomingSessions)

			if isUpcomingTabVisible && isSessionToday {
				upcomingTabData.Refresh()
			}
		}
	}()

	myApp.Run()
}
