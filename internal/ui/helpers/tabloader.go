package helpers

import (
	"github.com/mphartzheim/f1tray/internal/models"
	"github.com/mphartzheim/f1tray/internal/processes"
	"github.com/mphartzheim/f1tray/internal/ui/tabs"
	"github.com/mphartzheim/f1tray/internal/ui/tabs/results"
	"github.com/mphartzheim/f1tray/internal/ui/tabs/standings"

	"fyne.io/fyne/v2/container"
)

type TabsLoadResult struct {
	ScheduleTab        *container.TabItem
	ScheduleTabData    models.TabData
	UpcomingTab        *container.TabItem
	UpcomingTabData    models.TabData
	ResultsOuterTab    *container.TabItem
	ResultsTabData     models.TabData
	ResultsInnerTabs   *container.AppTabs
	StandingsOuterTab  *container.TabItem
	StandingsTabData   models.TabData
	StandingsInnerTabs *container.AppTabs
}

func LoadTabsConcurrently(state *models.AppState, selectedYear string) TabsLoadResult {
	type tabResult struct {
		name      string
		tab       *container.TabItem
		data      models.TabData
		innerTabs *container.AppTabs
	}

	tabResults := make(chan tabResult, 4)

	// Schedule
	go func() {
		data := tabs.CreateScheduleTableTab(processes.ParseSchedule, selectedYear)
		tabResults <- tabResult{name: "schedule", tab: container.NewTabItem("Schedule", data.Content), data: data}
	}()

	// Upcoming
	go func() {
		data := tabs.CreateUpcomingTab(state, processes.ParseUpcoming, selectedYear)
		tabResults <- tabResult{name: "upcoming", tab: container.NewTabItem("Upcoming", data.Content), data: data}
	}()

	// Results
	go func() {
		data, inner := results.CreateResultsTab(selectedYear, "last")
		tabResults <- tabResult{name: "results", tab: container.NewTabItem("Results", data.Content), data: data, innerTabs: inner}
	}()

	// Standings
	go func() {
		data, inner := standings.CreateStandingsTab(selectedYear, "last")
		tabResults <- tabResult{name: "standings", tab: container.NewTabItem("Standings", data.Content), data: data, innerTabs: inner}
	}()

	var res TabsLoadResult
	for i := 0; i < 4; i++ {
		r := <-tabResults
		switch r.name {
		case "schedule":
			res.ScheduleTab = r.tab
			res.ScheduleTabData = r.data
		case "upcoming":
			res.UpcomingTab = r.tab
			res.UpcomingTabData = r.data
		case "results":
			res.ResultsOuterTab = r.tab
			res.ResultsTabData = r.data
			res.ResultsInnerTabs = r.innerTabs
		case "standings":
			res.StandingsOuterTab = r.tab
			res.StandingsTabData = r.data
			res.StandingsInnerTabs = r.innerTabs
		}
	}

	return res
}
