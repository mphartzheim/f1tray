package appstate

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// AppState holds shared application state.
type AppState struct {
	Window       fyne.Window
	SelectedYear string
	Debug        bool
	OuterTabs    *container.AppTabs
	ResultsTabs  *container.AppTabs
}
