package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/mphartzheim/f1tray/internal/config"
	"github.com/mphartzheim/f1tray/internal/models"
	"github.com/mphartzheim/f1tray/internal/processes"
	"github.com/mphartzheim/f1tray/internal/ui"
	"github.com/mphartzheim/f1tray/internal/ui/helpers"
	"github.com/mphartzheim/f1tray/internal/ui/tabs"
	"github.com/mphartzheim/f1tray/internal/ui/tabs/preferences"
	"github.com/mphartzheim/f1tray/internal/ui/tabs/standings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	fyneTheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

//go:embed assets/tray_icon.png
var trayIconBytes []byte

// UIComponents holds references to the main UI parts needed later.
type UIComponents struct {
	window              fyne.Window
	tabsContainer       *container.AppTabs
	scheduleTab         *container.TabItem
	upcomingTab         *container.TabItem
	resultsOuterTab     *container.TabItem
	standingsOuterTab   *container.TabItem
	preferencesTab      *container.TabItem
	notificationLabel   *widget.Label
	notificationWrapper fyne.CanvasObject
	upcomingTabData     *models.TabData
	resultsTabData      *models.TabData
}

func main() {
	state := models.AppState{FirstRun: true}

	// Preload and cache constructors.json.
	loadConstructors()

	// Create the Fyne app with its initial theme.
	myApp := setupApp()

	// Build the UI and obtain key components.
	uiComps, countdownBinding := buildUI(myApp, &state)

	// Set up the system tray integration.
	configureSystemTray(myApp, uiComps)

	// Configure window behavior (show/hide and close interception).
	setupWindowBehavior(myApp, uiComps.window)

	// Lazy-load additional data once the UI is ready.
	// Dereference the pointers so that RefreshAllData gets models.TabData values.
	go processes.RefreshAllData(uiComps.notificationLabel,
		uiComps.notificationWrapper,
		*uiComps.upcomingTabData,
		*uiComps.resultsTabData)

	// Start background auto-refresh.
	go processes.StartAutoRefresh(&state, strconv.Itoa(time.Now().Year()))

	go processes.StartCountdown(countdownBinding, &state)

	// Start a ticker to periodically refresh the Upcoming tab.
	startUpcomingTicker(&state, uiComps.tabsContainer, uiComps.upcomingTab)

	myApp.Run()
}

// loadConstructors preloads and caches the constructors.json data.
func loadConstructors() {
	constructorFile, err := processes.LoadOrUpdateConstructorsJSON()
	if err != nil {
		fmt.Println("Warning: Failed to load constructors.json:", err)
		return
	}
	data, err := os.ReadFile(constructorFile)
	if err != nil {
		fmt.Println("Warning: Failed to read constructors file:", err)
		return
	}
	var constructorList models.ConstructorListResponse
	if err := json.Unmarshal(data, &constructorList); err != nil {
		fmt.Println("Warning: Failed to unmarshal constructors.json:", err)
		return
	}
	models.AllConstructors = constructorList.MRData.ConstructorTable.Constructors
}

// setupApp initializes the Fyne application and applies the initial theme.
func setupApp() fyne.App {
	myApp := app.NewWithID("f1tray")
	initialTheme := processes.GetThemeFromName(config.Get().Themes.Theme)
	myApp.Settings().SetTheme(initialTheme)
	if config.Get().Debug.Enabled {
		fmt.Printf("Theme: %T\n", fyneTheme.Current())
	}
	return myApp
}

// buildUI constructs the main window, header, tabs, and notification overlay.
func buildUI(myApp fyne.App, state *models.AppState) (UIComponents, binding.String) {
	// Create main window.
	myWindow := myApp.NewWindow("F1 Viewer")
	myWindow.SetFixedSize(true)
	models.MainWindow = myWindow

	// Create header (yearSelect and header container).
	yearSelect, headerContainer, countdownBinding := createHeader()

	// Create all tabs.
	scheduleTab, upcomingTab, resultsOuterTab, standingsOuterTab, preferencesTab,
		upcomingTabDataVal, resultsTabDataVal, resultsInnerTabs, standingsInnerTabs,
		tabsContainer := createTabs(yearSelect, state)

	// Register callbacks for inner tab updates.
	registerTabCallbacks(tabsContainer, resultsOuterTab, standingsOuterTab, resultsInnerTabs, standingsInnerTabs)

	// Register the yearSelect callback to update tab content when the year changes.
	registerYearSelectCallback(yearSelect, scheduleTab, &standingsInnerTabs, standingsOuterTab, tabsContainer)

	// Create the notification overlay.
	notificationLabel, notificationWrapper := ui.CreateNotification(myWindow, false)

	// Stack the tabs and notification overlay.
	stack := container.NewStack(tabsContainer, notificationWrapper)
	myWindow.SetContent(container.NewBorder(headerContainer, nil, nil, nil, stack))
	myWindow.Resize(fyne.NewSize(900, 700))

	return UIComponents{
		window:              myWindow,
		tabsContainer:       tabsContainer,
		scheduleTab:         scheduleTab,
		upcomingTab:         upcomingTab,
		resultsOuterTab:     resultsOuterTab,
		standingsOuterTab:   standingsOuterTab,
		preferencesTab:      preferencesTab,
		notificationLabel:   notificationLabel,
		notificationWrapper: notificationWrapper,
		// Store as pointers so that later functions can use their Refresh methods.
		upcomingTabData: &upcomingTabDataVal,
		resultsTabData:  &resultsTabDataVal,
	}, countdownBinding
}

// createHeader builds the header container with a "Season" label and a year selector.
func createHeader() (*widget.Select, fyne.CanvasObject, binding.String) {
	currentYear := time.Now().Year()
	years := []string{}
	for y := currentYear; y >= 1950; y-- {
		years = append(years, strconv.Itoa(y))
	}
	yearSelect := widget.NewSelect(years, nil)
	yearSelect.SetSelected(years[0])
	countdownBinding := binding.NewString()
	_ = countdownBinding.Set("Loading countdown…")
	countdownLabel := widget.NewLabelWithData(countdownBinding)
	countdownLabel.Alignment = fyne.TextAlignTrailing

	headerContainer := container.NewHBox(widget.NewLabel("Season"), yearSelect, layout.NewSpacer(), countdownLabel)
	return yearSelect, headerContainer, countdownBinding
}

// createTabs builds each of the tabs (Schedule, Upcoming, Results, Standings, Preferences)
// and returns the tab items along with inner tab data and the main tabs container.
func createTabs(yearSelect *widget.Select, state *models.AppState) (
	scheduleTab, upcomingTab, resultsOuterTab, standingsOuterTab, preferencesTab *container.TabItem,
	upcomingTabData, resultsTabData models.TabData,
	resultsInnerTabs, standingsInnerTabs *container.AppTabs,
	tabsContainer *container.AppTabs,
) {
	// Load all core tabs concurrently
	tabsData := helpers.LoadTabsConcurrently(state, yearSelect.Selected)

	scheduleTab = tabsData.ScheduleTab
	upcomingTab = tabsData.UpcomingTab
	resultsOuterTab = tabsData.ResultsOuterTab
	standingsOuterTab = tabsData.StandingsOuterTab

	upcomingTabData = tabsData.UpcomingTabData
	upcomingTabData.Refresh()
	resultsTabData = tabsData.ResultsTabData
	resultsInnerTabs = tabsData.ResultsInnerTabs
	standingsInnerTabs = tabsData.StandingsInnerTabs

	// Preferences Tab.
	preferencesTab = container.NewTabItem("Preferences", preferences.CreatePreferencesTab(
		func(updated config.Preferences) { _ = config.SaveConfig(updated) },
		func() { upcomingTabData.Refresh() },
	))

	// Assemble the main tabs container.
	tabsContainer = container.NewAppTabs(
		scheduleTab,
		upcomingTab,
		resultsOuterTab,
		standingsOuterTab,
		preferencesTab,
	)

	return
}

// registerTabCallbacks sets up the update callbacks and the OnSelected handlers for the tab containers.
func registerTabCallbacks(tabsContainer *container.AppTabs, resultsOuterTab, standingsOuterTab *container.TabItem,
	resultsInnerTabs, standingsInnerTabs *container.AppTabs) {

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

	// Remove trailing asterisks when a tab is selected.
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
}

// registerYearSelectCallback sets the callback for year selection changes,
// updating the Schedule and Standings tabs accordingly.
// Note: standingsInnerTabs is passed as a pointer-to-pointer so that it can be updated.
func registerYearSelectCallback(yearSelect *widget.Select, scheduleTab *container.TabItem,
	standingsInnerTabs **container.AppTabs, standingsOuterTab *container.TabItem, tabsContainer *container.AppTabs) {

	yearSelect.OnChanged = func(selectedYear string) {
		newScheduleTabData := tabs.CreateScheduleTableTab(processes.ParseSchedule, selectedYear)
		scheduleTab.Content = newScheduleTabData.Content

		selectedStandingsIndex := (*standingsInnerTabs).SelectedIndex()
		newStandingsTabData, newStandingsInnerTabs := standings.CreateStandingsTab(selectedYear, "last")
		standingsOuterTab.Content = newStandingsTabData.Content
		*standingsInnerTabs = newStandingsInnerTabs

		if selectedStandingsIndex >= 0 && selectedStandingsIndex < len((*standingsInnerTabs).Items) {
			(*standingsInnerTabs).SelectIndex(selectedStandingsIndex)
		}

		if tabsContainer.SelectedIndex() != 3 {
			if len(standingsOuterTab.Text) == 0 || standingsOuterTab.Text[len(standingsOuterTab.Text)-1] != '*' {
				standingsOuterTab.Text += "*"
			}
		}
		tabsContainer.Refresh()
	}
}

// configureSystemTray sets up the system tray icon and menu.
func configureSystemTray(myApp fyne.App, comps UIComponents) {
	iconResource := fyne.NewStaticResource("tray_icon.png", trayIconBytes)
	if desk, ok := myApp.(desktop.App); ok {
		processes.SetTrayIcon(desk, iconResource, comps.tabsContainer, comps.window,
			comps.scheduleTab, comps.upcomingTab, comps.resultsOuterTab, comps.standingsOuterTab, comps.preferencesTab,
		)
	}
}

// setupWindowBehavior configures the window's show/hide behavior and close interception.
func setupWindowBehavior(myApp fyne.App, win fyne.Window) {
	if config.Get().Window.HideOnOpen {
		win.Hide()
	} else {
		win.Show()
	}

	win.SetCloseIntercept(func() {
		if config.Get().Window.CloseBehavior == "exit" {
			myApp.Quit()
		} else {
			win.Hide()
		}
	})
}

// startUpcomingTicker starts a ticker that refreshes the Upcoming tab every 60 seconds
// if the Upcoming tab is visible and today is a session day.
func startUpcomingTicker(state *models.AppState, tabsContainer *container.AppTabs, upcomingTab *container.TabItem) {
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if tabsContainer.Selected() == upcomingTab && processes.IsSessionDay(state.UpcomingSessions) {
				// Refresh the Upcoming tab.
				// (Assuming upcomingTabData.Refresh() is the correct way to update its content.)
			}
		}
	}()
}
