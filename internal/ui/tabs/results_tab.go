package tabs

import (
	"fmt"

	"f1tray/internal/models"
	"f1tray/internal/processes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// CreateResultsTableTab builds a tab displaying race results fetched from a URL and parsed into a formatted table.
func CreateResultsTableTab(url string, parseFunc func([]byte) (string, [][]string, error), year string, round string) models.TabData {
	status := widget.NewLabel("Loading results...")
	raceNameLabel := widget.NewLabel("")
	tableContainer := container.NewStack()

	url = fmt.Sprintf(url, year)

	refresh := func() bool {
		data, changed, err := processes.FetchData(url)
		if err != nil {
			status.SetText("Failed to fetch results.")
			return false
		}
		if !changed {
			return false
		}

		raceName, rows, err := parseFunc(data)
		if err != nil {
			status.SetText("Failed to parse results.")
			return false
		}

		raceNameLabel.SetText(fmt.Sprintf("Results for: %s", raceName))
		table := widget.NewTable(
			func() (int, int) {
				if len(rows) == 0 {
					return 0, 0
				}
				return len(rows), len(rows[0])
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("")
			},
			func(id widget.TableCellID, cell fyne.CanvasObject) {
				cell.(*widget.Label).SetText(rows[id.Row][id.Col])
			},
		)

		// Restore proper column widths for results view
		table.SetColumnWidth(0, 50)  // Pos
		table.SetColumnWidth(1, 180) // Driver
		table.SetColumnWidth(2, 180) // Team
		table.SetColumnWidth(3, 300) // Time/Status

		tableContainer.Objects = []fyne.CanvasObject{table}
		tableContainer.Refresh()
		status.SetText("Results loaded.")
		return true
	}

	refresh()

	content := container.NewBorder(
		container.NewVBox(raceNameLabel), // Top
		status,                           // Bottom
		nil, nil,
		tableContainer,
	)

	return models.TabData{
		Content: content,
		Refresh: refresh,
	}
}
