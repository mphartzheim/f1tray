package tabs

import (
	"encoding/json"
	"fmt"
	"image/color"
	"time"

	"f1tray/internal/models"
	"f1tray/internal/processes"
	"f1tray/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// CreateScheduleTableTab builds a tab displaying the full race schedule with interactive circuit links and highlighted upcoming event.
func CreateScheduleTableTab(parseFunc func([]byte) (string, [][]string, error), year string) models.TabData {
	status := widget.NewLabel("Loading schedule...")
	tableContainer := container.NewStack()

	// Format URL with the season/year.
	url := fmt.Sprintf(models.ScheduleURL, year)

	refresh := func() bool {
		data, err := processes.FetchData(url)
		if err != nil {
			status.SetText("Failed to fetch schedule")
			return false
		}

		title, rows, err := parseFunc(data)
		if err != nil {
			status.SetText("Failed to parse schedule")
			return false
		}

		// Unmarshal into our struct that now includes lat/long as strings.
		var schedule models.ScheduleResponse
		err = json.Unmarshal(data, &schedule)
		if err != nil {
			status.SetText(fmt.Sprintf("Error parsing schedule JSON: %v", err))
			return false
		}

		now := time.Now()
		highlightRow := -1
		for i, race := range schedule.MRData.RaceTable.Races {
			raceDate, _ := time.Parse("2006-01-02", race.Date)
			if raceDate.After(now) || raceDate.Equal(now) {
				highlightRow = i + 1 // +1 because headers take row 0
				break
			}
		}

		// Factory function returns a container with a background rectangle and a clickable label.
		factory := func() fyne.CanvasObject {
			bg := canvas.NewRectangle(nil)
			cl := ui.NewClickableLabel("", nil)
			return container.NewStack(bg, cl)
		}

		// Update function for each cell.
		update := func(id widget.TableCellID, obj fyne.CanvasObject) {
			wrapper := obj.(*fyne.Container)
			bg := wrapper.Objects[0].(*canvas.Rectangle)
			cl := wrapper.Objects[1].(*ui.ClickableLabel)

			if id.Row == 0 {
				// Header row.
				headers := []string{"Round", "Race Name", "Circuit", "Location (Date)"}
				cl.SetText(headers[id.Col])
				cl.OnDoubleTapped = nil
				bg.Hide()
			} else {
				cl.SetText(rows[id.Row-1][id.Col])

				race := schedule.MRData.RaceTable.Races[id.Row-1]

				if id.Col == 1 {
					// Parse the race date to determine if it has occurred.
					raceDate, err := time.Parse("2006-01-02", race.Date)
					if err != nil {
						fmt.Printf("Error parsing race date: %v\n", err)
						cl.OnDoubleTapped = nil
					} else if raceDate.Before(time.Now()) {
						// Only enable double-click if the race has already happened.
						cl.OnDoubleTapped = func() {
							round := rows[id.Row-1][0]

							newResultsTab := CreateResultsTableTab(processes.ParseRaceResults, year, round)
							newQualifyingTab := CreateResultsTableTab(processes.ParseQualifyingResults, year, round)
							newSprintTab := CreateResultsTableTab(processes.ParseSprintResults, year, round)

							newResultsTab.Refresh()
							newQualifyingTab.Refresh()
							newSprintTab.Refresh()

							processes.ReloadOtherTabs(newResultsTab.Content, newQualifyingTab.Content, newSprintTab.Content)
						}
					} else {
						// Future races — disable double-click
						cl.OnDoubleTapped = nil
					}
				} else if id.Col == 2 {
					lat := race.Circuit.Location.Lat
					lon := race.Circuit.Location.Long
					mapURL := fmt.Sprintf("%s?mlat=%s&mlon=%s#map=15/%s/%s", models.MapBaseURL, lat, lon, lat, lon)
					cl.OnDoubleTapped = func() {
						if err := ui.OpenWebPage(mapURL); err != nil {
							status.SetText("Failed to open map URL")
						}
					}
				} else {
					cl.OnDoubleTapped = nil
				}

				if id.Row == highlightRow {
					bg.StrokeColor = color.NRGBA{R: 0xFF, G: 0x18, B: 0x01, A: 0xFF} // Red border.
					bg.StrokeWidth = 2
					bg.Show()
				} else {
					bg.Hide()
				}
				bg.Resize(wrapper.Size())
			}
			wrapper.Refresh()
		}

		table := widget.NewTable(
			func() (int, int) { return len(rows) + 1, 4 },
			factory,
			update,
		)

		table.SetColumnWidth(0, 60)
		table.SetColumnWidth(1, 200)
		table.SetColumnWidth(2, 280)
		table.SetColumnWidth(3, 280)
		table.Resize(fyne.NewSize(820, float32((len(rows)+1)*30)))

		tableContainer.Objects = []fyne.CanvasObject{table}
		tableContainer.Refresh()
		status.SetText(title)
		return true
	}

	refresh()

	content := container.NewBorder(nil, status, nil, nil, tableContainer)
	return models.TabData{
		Content: content,
		Refresh: refresh,
	}
}
