package tabs

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"f1tray/internal/models"
	"f1tray/internal/processes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// CreateResultsTableTab builds a tab displaying race results fetched from a URL and parsed into a formatted table.
func CreateResultsTableTab(parseFunc func([]byte) (string, [][]string, error), year string, round string) models.TabData {
	status := widget.NewLabel("Loading results...")
	raceNameLabel := widget.NewLabel("")
	tableContainer := container.NewStack()

	url := buildResultsURL(parseFunc, year, round)

	refresh := func() bool {
		data, changed, err := processes.FetchData(url)
		if err != nil {
			status.SetText("Failed to fetch results.")
			return false
		}
		if !changed {
			return false
		}

		// Get the function name for later use.
		funcName := runtime.FuncForPC(reflect.ValueOf(parseFunc).Pointer()).Name()

		// Parse the data.
		raceName, rows, err := parseFunc(data)
		// If there's an error and we're dealing with sprint results, check if it's due to no data.
		if err != nil {
			if strings.HasSuffix(funcName, "ParseSprintResults") && strings.Contains(err.Error(), "no sprint data found") {
				raceNameLabel.SetText("Not a sprint race event")
				tableContainer.Objects = nil
				tableContainer.Refresh()
				status.SetText("Results loaded.")
				return true
			}
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

		// Restore proper column widths for results view.
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

// buildResultsURL builds the correct URL based on the parse function provided.
func buildResultsURL(parseFunc func([]byte) (string, [][]string, error), year string, round string) string {
	// Get the function name using runtime and reflect.
	funcName := runtime.FuncForPC(reflect.ValueOf(parseFunc).Pointer()).Name()

	// Decide which URL to use based on the function name.
	if strings.HasSuffix(funcName, "ParseRaceResults") {
		return fmt.Sprintf(models.RaceResultsURL, year, round)
	} else if strings.HasSuffix(funcName, "ParseQualifyingResults") {
		return fmt.Sprintf(models.QualifyingURL, year, round)
	} else if strings.HasSuffix(funcName, "ParseSprintResults") {
		return fmt.Sprintf(models.SprintURL, year, round)
	}
	// Default to empty string if no match is found.
	return ""
}
