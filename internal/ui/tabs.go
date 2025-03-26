package ui

import (
	"f1tray/internal/processes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func CreateResultsTableTab(url string, parseFunc func([]byte) (string, [][]string, error)) fyne.CanvasObject {
	status := widget.NewLabel("Loading results...")
	tableContainer := container.NewStack()

	go processes.LoadResults(url, parseFunc, status, tableContainer)

	return container.NewBorder(nil, status, nil, nil, tableContainer)
}

func CreateScheduleTableTab(url string, parseFunc func([]byte) (string, [][]string, error)) fyne.CanvasObject {
	status := widget.NewLabel("Press 'Load Schedule' to fetch data.")
	tableContainer := container.NewStack()

	go processes.LoadSchedule(url, parseFunc, status, tableContainer)

	return container.NewBorder(nil, status, nil, nil, tableContainer)
}

func CreateUpcomingTab(url string, parseFunc func([]byte) (string, [][]string, error)) fyne.CanvasObject {
	status := widget.NewLabel("Press 'Load Upcoming' to fetch data.")
	tableContainer := container.NewStack()

	go processes.LoadUpcoming(url, parseFunc, status, tableContainer)

	return container.NewBorder(nil, status, nil, nil, tableContainer)

}
