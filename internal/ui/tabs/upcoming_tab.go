package tabs

import (
	"encoding/json"
	"fmt"
	"time"

	"f1tray/internal/models"
	"f1tray/internal/processes"
	"f1tray/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// CreateUpcomingTab builds a tab showing upcoming race sessions, with a clickable label for map access and a link to F1TV.
func CreateUpcomingTab(parseFunc func([]byte) (string, [][]string, error), year string) models.TabData {
	status := widget.NewLabel("Loading upcoming races...")
	// New double-clickable label for "Next Race"
	nextRaceLabel := ui.NewClickableLabel("Next Race", nil)
	tableContainer := container.NewStack()

	url := fmt.Sprintf(models.UpcomingURL, year)

	// Create the "Watch on F1TV" button.
	watchButton := widget.NewButton("Watch on F1TV", func() {
		if err := ui.OpenWebPage(models.F1tvURL); err != nil {
			status.SetText("Failed to open F1TV URL.")
		}
	})

	refresh := func() bool {
		data, changed, err := processes.FetchData(url)
		if err != nil {
			status.SetText("Failed to fetch upcoming data.")
			return false
		}
		if !changed {
			return false
		}

		// We still call the parse function to potentially update other parts of the UI,
		// but we no longer use its title output.
		_, rows, err := parseFunc(data)
		if err != nil {
			status.SetText("Failed to parse upcoming data.")
			return false
		}

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

		table.SetColumnWidth(0, 150)
		table.SetColumnWidth(1, 150)
		table.SetColumnWidth(2, 150)
		tableContainer.Objects = []fyne.CanvasObject{table}
		tableContainer.Refresh()
		status.SetText("Upcoming race loaded.")

		// Unmarshal into our struct (assuming upcoming data uses similar structure)
		var schedule models.ScheduleResponse
		err = json.Unmarshal(data, &schedule)
		if err != nil {
			nextRaceLabel.SetText("Next Race (map unavailable)")
			nextRaceLabel.OnDoubleTapped = nil
		} else {
			// Determine the next upcoming race
			now := time.Now()
			found := false
			for _, race := range schedule.MRData.RaceTable.Races {
				raceDate, err := time.Parse("2006-01-02", race.Date)
				if err != nil {
					continue
				}
				if raceDate.After(now) || raceDate.Equal(now) {
					// Update the label to display both the race name and circuit.
					nextRaceLabel.SetText(fmt.Sprintf("Next Race: %s (%s)", race.RaceName, race.Circuit.CircuitName))
					lat := race.Circuit.Location.Lat
					lon := race.Circuit.Location.Long
					// Build the OpenStreetMap URL with a default zoom level of 15.
					mapURL := fmt.Sprintf("%s?mlat=%s&mlon=%s#map=15/%s/%s", models.MapBaseURL, lat, lon, lat, lon)
					nextRaceLabel.OnDoubleTapped = func() {
						if err := ui.OpenWebPage(mapURL); err != nil {
							status.SetText("Failed to open map URL")
						}
					}
					found = true
					break
				}
			}
			if !found {
				nextRaceLabel.SetText("Next Race: Not available")
				nextRaceLabel.OnDoubleTapped = nil
			}
		}

		return true
	}

	refresh()

	// Use only the clickable label in the top layout.
	topContent := container.NewVBox(nextRaceLabel)
	// Combine the status label and the watch button in a vertical layout at the bottom.
	bottomContent := container.NewVBox(status, watchButton)

	content := container.NewBorder(
		topContent,    // Top
		bottomContent, // Bottom
		nil, nil,
		tableContainer,
	)

	return models.TabData{
		Content: content,
		Refresh: refresh,
	}
}
