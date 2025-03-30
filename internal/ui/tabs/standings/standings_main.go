package standings

import (
	"github.com/mphartzheim/f1tray/internal/models"
	"github.com/mphartzheim/f1tray/internal/processes"

	"fyne.io/fyne/v2/container"
)

// CreateStandingsTab combines Race, Qualifying, and Sprint tabs into one unified Results tab.
func CreateStandingsTab(year string, round string) (models.TabData, *container.AppTabs) {
	// Create individual sub-tabs using your existing result-building logic
	driversTab := CreateStandingsTableTab(processes.ParseDriverStandings, year)
	constructorsTab := CreateStandingsTableTab(processes.ParseConstructorStandings, year)

	// Create internal tab bar
	nestedTabs := container.NewAppTabs(
		container.NewTabItem("Drivers", driversTab.Content),
		container.NewTabItem("Constructors", constructorsTab.Content),
	)
	nestedTabs.SetTabLocation(container.TabLocationTop)

	// Combined Results tab with refresh support
	return models.TabData{
		Content: nestedTabs,
		Refresh: func() bool {
			driversTab.Refresh()
			constructorsTab.Refresh()
			return true
		},
	}, nestedTabs
}
