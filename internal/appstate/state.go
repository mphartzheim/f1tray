package appstate

import "fyne.io/fyne/v2"

// AppState holds shared application state.
type AppState struct {
	Window       fyne.Window
	SelectedYear string
	Debug        bool
}
