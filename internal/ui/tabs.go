package ui

import (
	"f1tray/internal/config"
	"f1tray/internal/processes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// TabData encapsulates a tab's content and its refresh function.
type TabData struct {
	Content fyne.CanvasObject
	Refresh func()
}

func CreateResultsTableTab(url string, parseFunc func([]byte) (string, [][]string, error)) TabData {
	status := widget.NewLabel("Loading results...")
	tableContainer := container.NewStack()

	// Define refresh function to load data.
	refresh := func() {
		processes.LoadResults(url, parseFunc, status, tableContainer)
	}

	// Load data initially.
	go refresh()

	content := container.NewBorder(nil, status, nil, nil, tableContainer)
	return TabData{
		Content: content,
		Refresh: refresh,
	}
}

func CreateScheduleTableTab(url string, parseFunc func([]byte) (string, [][]string, error)) TabData {
	status := widget.NewLabel("Loading schedule...")
	tableContainer := container.NewStack()

	refresh := func() {
		processes.LoadSchedule(url, parseFunc, status, tableContainer)
	}

	go refresh()

	content := container.NewBorder(nil, status, nil, nil, tableContainer)
	return TabData{
		Content: content,
		Refresh: refresh,
	}
}

func CreateUpcomingTab(url string, parseFunc func([]byte) (string, [][]string, error)) TabData {
	status := widget.NewLabel("Loading upcoming races...")
	tableContainer := container.NewStack()

	refresh := func() {
		processes.LoadUpcoming(url, parseFunc, status, tableContainer)
	}

	go refresh()

	content := container.NewBorder(nil, status, nil, nil, tableContainer)
	return TabData{
		Content: content,
		Refresh: refresh,
	}
}

func CreatePreferencesTab(currentPrefs config.Preferences, onSave func(config.Preferences)) fyne.CanvasObject {
	isExit := currentPrefs.CloseBehavior == "exit"

	checkbox := widget.NewCheck("Close on exit?", func(checked bool) {
		if checked {
			currentPrefs.CloseBehavior = "exit"
		} else {
			currentPrefs.CloseBehavior = "minimize"
		}
		onSave(currentPrefs)
	})
	checkbox.SetChecked(isExit)

	return container.NewVBox(
		widget.NewLabel("Window Close Behavior:"),
		checkbox,
	)
}
