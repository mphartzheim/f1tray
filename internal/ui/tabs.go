package ui

import (
	"f1tray/internal/processes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func CreateResultsTableTab(url string, parseFunc func([]byte) (string, [][]string, error)) fyne.CanvasObject {
	status := widget.NewLabel("Press 'Load Results' to fetch data.")
	tableContainer := container.NewStack()

	loadButton := widget.NewButton("Load Results", func() {
		processes.LoadResults(url, parseFunc, status, tableContainer)
	})

	return container.NewBorder(loadButton, status, nil, nil, tableContainer)
}

func CreateScheduleTableTab(url string, parseFunc func([]byte) (string, [][]string, error)) fyne.CanvasObject {
	status := widget.NewLabel("Press 'Load Schedule' to fetch data.")
	tableContainer := container.NewStack()

	loadButton := widget.NewButton("Load Schedule", func() {
		processes.LoadSchedule(url, parseFunc, status, tableContainer)
	})

	return container.NewBorder(loadButton, status, nil, nil, tableContainer)
}

func CreateUpcomingTab(url string, parseFunc func([]byte) (string, [][]string, error)) fyne.CanvasObject {
	status := widget.NewLabel("Press 'Load Upcoming' to fetch data.")
	tableContainer := container.NewStack()

	loadButton := widget.NewButton("Load Upcoming", func() {
		processes.LoadUpcoming(url, parseFunc, status, tableContainer)
	})

	return container.NewBorder(loadButton, status, nil, nil, tableContainer)

}
