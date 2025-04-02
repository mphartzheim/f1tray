package layout

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func CreateResultsTabs(race, qual, sprint *widget.Label) *container.AppTabs {
	return container.NewAppTabs(
		container.NewTabItem("Race", race),
		container.NewTabItem("Qualifying", qual),
		container.NewTabItem("Sprint", sprint),
	)
}

func CreateStandingsTabs(driver, constructor *widget.Label) *container.AppTabs {
	return container.NewAppTabs(
		container.NewTabItem("Driver", driver),
		container.NewTabItem("Constructor", constructor),
	)
}

func CreatePreferencesTabs(main, notifications *widget.Label) *container.AppTabs {
	return container.NewAppTabs(
		container.NewTabItem("Main", main),
		container.NewTabItem("Notifications", notifications),
	)
}
