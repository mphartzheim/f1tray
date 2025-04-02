package preferences

import (
	"github.com/mphartzheim/f1tray/internal/appstate"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func BuildPreferencesTab(state *appstate.AppState) *container.TabItem {
	themeSelect := ThemeSelector(func(selected string) {
		ApplyThemeSelection(fyne.CurrentApp(), selected)
	})

	clockToggle := ClockToggle()
	trayToggle := TrayToggle()
	notifySlider := NotificationSlider()

	layout := container.NewVBox(
		container.NewHBox(widget.NewLabel("Theme:"), layout.NewSpacer(), themeSelect),
		clockToggle,
		trayToggle,
		container.NewHBox(widget.NewLabel("Notify X minutes before session:"), layout.NewSpacer(), notifySlider),
	)

	return container.NewTabItem("Main", layout)
}
