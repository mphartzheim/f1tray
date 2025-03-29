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
	nextRaceLabel := ui.NewClickableLabel("Next Race", nil, false)
	tableContainer := container.NewStack()

	url := fmt.Sprintf(models.UpcomingURL, year)

	watchButton := widget.NewButton("Watch on F1TV", func() {
		if err := ui.OpenWebPage(models.F1tvURL); err != nil {
			status.SetText("Failed to open F1TV URL.")
		}
	})

	refresh := func() bool {
		data, err := processes.FetchData(url)
		if err != nil {
			status.SetText("Failed to fetch upcoming data.")
			return false
		}

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

		var schedule models.ScheduleResponse
		err = json.Unmarshal(data, &schedule)
		if err != nil {
			nextRaceLabel.SetText("Next Race (map unavailable)")
			nextRaceLabel.OnDoubleTapped = nil
		} else {
			now := time.Now()
			found := false
			for _, race := range schedule.MRData.RaceTable.Races {
				raceDate, err := time.Parse("2006-01-02", race.Date)
				if err != nil {
					continue
				}
				if raceDate.After(now) || raceDate.Equal(now) {
					nextRaceLabel.SetText(fmt.Sprintf("Next Race: %s (%s 🗺️)", race.RaceName, race.Circuit.CircuitName))
					lat := race.Circuit.Location.Lat
					lon := race.Circuit.Location.Long
					mapURL := fmt.Sprintf("%s?mlat=%s&mlon=%s#map=15/%s/%s", models.MapBaseURL, lat, lon, lat, lon)
					nextRaceLabel.OnDoubleTapped = func() {
						if err := ui.OpenWebPage(mapURL); err != nil {
							status.SetText("Failed to open map URL")
						}
					}
					nextRaceLabel.Clickable = true
					found = true
					break
				}
			}
			if !found {
				nextRaceLabel.SetText("Next Race: Not available")
				nextRaceLabel.OnDoubleTapped = nil
				nextRaceLabel.Clickable = false
			}
		}

		return true
	}

	refresh()

	topContent := container.NewVBox(nextRaceLabel)
	bottomContent := container.NewVBox(status, watchButton)

	content := container.NewBorder(
		topContent,
		bottomContent,
		nil, nil,
		tableContainer,
	)

	return models.TabData{
		Content: content,
		Refresh: refresh,
	}
}
