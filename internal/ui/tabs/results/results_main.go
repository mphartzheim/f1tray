package results

import (
	"github.com/mphartzheim/f1tray/internal/models"
	"github.com/mphartzheim/f1tray/internal/processes"

	"fyne.io/fyne/v2/container"
)

// CreateResultsTab combines Race, Qualifying, and Sprint tabs into one unified Results tab.
func CreateResultsTab(year string, round string) (models.TabData, *container.AppTabs) {
	// Create individual sub-tabs using your existing result-building logic
	raceTab := CreateResultsTableTab(processes.ParseRaceResults, year, round)
	qualifyingTab := CreateResultsTableTab(processes.ParseQualifyingResults, year, round)
	sprintTab := CreateResultsTableTab(processes.ParseSprintResults, year, round)

	// Create internal tab bar
	nestedTabs := container.NewAppTabs(
		container.NewTabItem("Race", raceTab.Content),
		container.NewTabItem("Qualifying", qualifyingTab.Content),
		container.NewTabItem("Sprint", sprintTab.Content),
	)
	nestedTabs.SetTabLocation(container.TabLocationTop)

	// Combined Results tab with refresh support
	return models.TabData{
		Content: nestedTabs,
		Refresh: func() bool {
			raceTab.Refresh()
			qualifyingTab.Refresh()
			sprintTab.Refresh()
			return true
		},
	}, nestedTabs
}
