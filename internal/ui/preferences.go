package ui

import (
	"f1tray/internal/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func CreatePreferencesTab(currentPrefs config.Preferences, onSave func(config.Preferences)) fyne.CanvasObject {
	options := []string{"minimize", "exit"}
	selectWidget := widget.NewSelect(options, func(selected string) {
		currentPrefs.CloseBehavior = selected
		onSave(currentPrefs)
	})

	selectWidget.SetSelected(currentPrefs.CloseBehavior)

	return container.NewVBox(
		widget.NewLabel("Window Close Behavior:"),
		selectWidget,
	)
}
