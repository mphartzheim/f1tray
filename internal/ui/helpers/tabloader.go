package helpers

import (
	"github.com/mphartzheim/f1tray/internal/models"
	"github.com/mphartzheim/f1tray/internal/processes"
	"github.com/mphartzheim/f1tray/internal/ui/tabs"
	"github.com/mphartzheim/f1tray/internal/ui/tabs/results"
	"github.com/mphartzheim/f1tray/internal/ui/tabs/standings"

	"fyne.io/fyne/v2/container"
)

type TabResults struct {
	ScheduleTabData  models.TabData
	UpcomingTabData  models.TabData
	ResultsTabData   models.TabData
	StandingsTabData models.TabData

	ScheduleTab       *container.TabItem
	UpcomingTab       *container.TabItem
	ResultsOuterTab   *container.TabItem
	StandingsOuterTab *container.TabItem

	ResultsInnerTabs   *container.AppTabs
	StandingsInnerTabs *container.AppTabs
}

func LoadTabsConcurrently(state *models.AppState, selectedYear string) TabResults {
	type tabResult struct {
		name      string
		tabItem   *container.TabItem
		tabData   models.TabData
		innerTabs *container.AppTabs
	}

	scheduleCh := make(chan tabResult)
	upcomingCh := make(chan tabResult)
	resultsCh := make(chan tabResult)
	standingsCh := make(chan tabResult)

	go func() {
		data := tabs.CreateScheduleTableTab(processes.ParseSchedule, selectedYear)
		scheduleCh <- tabResult{name: "Schedule", tabItem: container.NewTabItem("Schedule", data.Content), tabData: data}
	}()

	go func() {
		data := tabs.CreateUpcomingTab(state, processes.ParseUpcoming, selectedYear)
		upcomingCh <- tabResult{name: "Upcoming", tabItem: container.NewTabItem("Upcoming", data.Content), tabData: data}
	}()

	go func() {
		data, inner := results.CreateResultsTab(selectedYear, "last")
		resultsCh <- tabResult{name: "Results", tabItem: container.NewTabItem("Results", data.Content), tabData: data, innerTabs: inner}
	}()

	go func() {
		data, inner := standings.CreateStandingsTab(selectedYear, selectedYear)
		standingsCh <- tabResult{name: "Standings", tabItem: container.NewTabItem("Standings", data.Content), tabData: data, innerTabs: inner}
	}()

	schedule := <-scheduleCh
	upcoming := <-upcomingCh
	results := <-resultsCh
	standings := <-standingsCh

	return TabResults{
		ScheduleTabData:    schedule.tabData,
		UpcomingTabData:    upcoming.tabData,
		ResultsTabData:     results.tabData,
		StandingsTabData:   standings.tabData,
		ScheduleTab:        schedule.tabItem,
		UpcomingTab:        upcoming.tabItem,
		ResultsOuterTab:    results.tabItem,
		StandingsOuterTab:  standings.tabItem,
		ResultsInnerTabs:   results.innerTabs,
		StandingsInnerTabs: standings.innerTabs,
	}
}
