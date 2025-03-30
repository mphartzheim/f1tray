package results

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/mphartzheim/f1tray/internal/models"
	"github.com/mphartzheim/f1tray/internal/processes"
	"github.com/mphartzheim/f1tray/internal/ui"

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
		data, err := processes.FetchData(url)
		if err != nil {
			status.SetText("Failed to fetch results.")
			return false
		}

		// Get the function name for later use.
		funcName := runtime.FuncForPC(reflect.ValueOf(parseFunc).Pointer()).Name()

		// Parse the data.
		raceName, rows, err := parseFunc(data)
		// If there's an error and we're dealing with sprint or qualifying results, handle it gracefully.
		if err != nil {
			if strings.HasSuffix(funcName, "ParseSprintResults") && strings.Contains(err.Error(), "no sprint data found") {
				raceNameLabel.SetText("Not a sprint race event")
				tableContainer.Objects = nil
				tableContainer.Refresh()
				status.SetText("Results loaded")
				return true
			} else if strings.HasSuffix(funcName, "ParseQualifyingResults") && strings.Contains(err.Error(), "no qualifying data found") {
				raceNameLabel.SetText("No data available on Jolpica API")
				tableContainer.Objects = nil
				tableContainer.Refresh()
				status.SetText("Results loaded.")
				return true
			}
			status.SetText("Failed to parse results")
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
			// Create each cell as a container so we can swap its content if needed.
			func() fyne.CanvasObject {
				return container.NewStack(widget.NewLabel(""))
			},
			// Update each cell.
			func(id widget.TableCellID, co fyne.CanvasObject) {
				cont, ok := co.(*fyne.Container)
				if !ok {
					return
				}
				cont.Objects = nil
				text := rows[id.Row][id.Col]
				var cellWidget fyne.CanvasObject

				// Apply clickable logic for the Driver column (index 1)
				if id.Col == 1 {
					// Check for our delimiter indicating a clickable cell.
					if strings.Contains(text, "|||") {
						parts := strings.SplitN(text, "|||", 2)
						displayName := parts[0]
						fallback := strings.TrimSuffix(parts[1], " 👤")
						clickableText := fmt.Sprintf("%s 👤", displayName)
						// Use custom URL if available, otherwise fallback to the API URL.
						if slug, ok := models.DriverURLMap[displayName]; ok {
							url := fmt.Sprintf(models.F1DriverBioURL, slug)
							cellWidget = ui.NewClickableLabel(clickableText, func() {
								processes.OpenWebPage(url)
							}, true)
						} else {
							cellWidget = ui.NewClickableLabel(clickableText, func() {
								processes.OpenWebPage(fallback)
							}, true)
						}
					} else {
						cellWidget = widget.NewLabel(text)
					}
				} else {
					cellWidget = widget.NewLabel(text)
				}
				cont.Add(cellWidget)
				cont.Refresh()
			},
		)

		// Set proper column widths.
		table.SetColumnWidth(0, 50)  // Position
		table.SetColumnWidth(1, 180) // Driver
		table.SetColumnWidth(2, 180) // Team
		table.SetColumnWidth(3, 300) // Time/Status

		tableContainer.Objects = []fyne.CanvasObject{table}
		tableContainer.Refresh()
		status.SetText("Results loaded")
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
	funcName := runtime.FuncForPC(reflect.ValueOf(parseFunc).Pointer()).Name()

	if strings.HasSuffix(funcName, "ParseRaceResults") {
		return fmt.Sprintf(models.RaceURL, year, round)
	} else if strings.HasSuffix(funcName, "ParseQualifyingResults") {
		return fmt.Sprintf(models.QualifyingURL, year, round)
	} else if strings.HasSuffix(funcName, "ParseSprintResults") {
		return fmt.Sprintf(models.SprintURL, year, round)
	}
	return ""
}
