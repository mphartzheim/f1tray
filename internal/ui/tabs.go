package ui

import (
	"f1tray/internal/config"
	"f1tray/internal/models"
	"f1tray/internal/processes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func CreateResultsTableTab(url string, parseFunc func([]byte) (string, [][]string, error)) models.TabData {
	status := widget.NewLabel("Loading results...")
	tableContainer := container.NewStack()

	// Define refresh function to load data.
	refresh := func() bool {
		return processes.LoadResults(url, parseFunc, status, tableContainer)
	}

	// Load data initially.
	go func() { refresh() }()

	content := container.NewBorder(nil, status, nil, nil, tableContainer)
	return models.TabData{
		Content: content,
		Refresh: refresh,
	}
}

func CreateScheduleTableTab(url string, parseFunc func([]byte) (string, [][]string, error)) models.TabData {
	status := widget.NewLabel("Loading schedule...")
	tableContainer := container.NewStack()

	refresh := func() bool {
		return processes.LoadSchedule(url, parseFunc, status, tableContainer)
	}

	go func() { refresh() }()

	content := container.NewBorder(nil, status, nil, nil, tableContainer)
	return models.TabData{
		Content: content,
		Refresh: refresh,
	}
}

func CreateUpcomingTab(url string, parseFunc func([]byte) (string, [][]string, error)) models.TabData {
	status := widget.NewLabel("Loading upcoming races...")
	tableContainer := container.NewStack()

	refresh := func() bool {
		return processes.LoadUpcoming(url, parseFunc, status, tableContainer)
	}

	go func() { refresh() }()

	content := container.NewBorder(nil, status, nil, nil, tableContainer)
	return models.TabData{
		Content: content,
		Refresh: refresh,
	}
}

func CreatePreferencesTab(currentPrefs config.Preferences, onSave func(config.Preferences)) fyne.CanvasObject {
	// Checkbox for Close Behavior.
	isExit := currentPrefs.CloseBehavior == "exit"
	closeCheckbox := widget.NewCheck("Close on exit?", func(checked bool) {
		if checked {
			currentPrefs.CloseBehavior = "exit"
		} else {
			currentPrefs.CloseBehavior = "minimize"
		}
		onSave(currentPrefs)
	})
	closeCheckbox.SetChecked(isExit)

	// Checkbox for Hide on Open.
	hideCheckbox := widget.NewCheck("Hide on open?", func(checked bool) {
		currentPrefs.HideOnOpen = checked
		onSave(currentPrefs)
	})
	hideCheckbox.SetChecked(currentPrefs.HideOnOpen)

	return container.NewVBox(
		widget.NewLabel("Window Close Behavior:"),
		closeCheckbox,
		widget.NewLabel("Window Open Behavior:"),
		hideCheckbox,
	)
}
