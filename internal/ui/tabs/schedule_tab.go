package tabs

import (
	"encoding/json"
	"fmt"
	"time"

	"f1tray/internal/models"
	"f1tray/internal/processes"
	"f1tray/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CreateScheduleTableTab builds a tab displaying the full race schedule with interactive circuit links and highlighted upcoming event.
func CreateScheduleTableTab(url string, parseFunc func([]byte) (string, [][]string, error), year string) models.TabData {
	status := widget.NewLabel("Loading schedule...")
	tableContainer := container.NewStack()

	// Format URL with the season/year.
	url = fmt.Sprintf(url, year)

	refresh := func() bool {
		data, _, err := processes.FetchData(url)
		if err != nil {
			status.SetText("Failed to fetch schedule.")
			return false
		}

		title, rows, err := parseFunc(data)
		if err != nil {
			status.SetText("Failed to parse schedule.")
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
				if id.Col == 1 {
					// For the Race Name column: extract the round from column 0 and reload other tabs.
					cl.OnDoubleTapped = func() {
						round := rows[id.Row-1][0]
						fmt.Printf("Reloading tabs for season %s and round %s\n", year, round)

						// Build new endpoints using both year and round.
						// (Ensure your URL strings in models are formatted to accept two parameters.)
						resultsURL := fmt.Sprintf(models.RaceResultsURL, year, round)
						qualifyingURL := fmt.Sprintf(models.QualifyingURL, year, round)
						sprintURL := fmt.Sprintf(models.SprintURL, year, round)

						// Debug print the constructed URLs.
						fmt.Printf("Results URL: %s\n", resultsURL)
						fmt.Printf("Qualifying URL: %s\n", qualifyingURL)
						fmt.Printf("Sprint URL: %s\n", sprintURL)

						// Create new tab data for each tab.
						newResultsTab := CreateResultsTableTab(resultsURL, processes.ParseRaceResults, year, round)
						newQualifyingTab := CreateResultsTableTab(qualifyingURL, processes.ParseQualifyingResults, year, round)
						newSprintTab := CreateResultsTableTab(sprintURL, processes.ParseSprintResults, year, round)

						// Optionally trigger the refresh if needed.
						newResultsTab.Refresh()
						newQualifyingTab.Refresh()
						newSprintTab.Refresh()

						// Pass the new content to the UpdateTabs callback.
						processes.ReloadOtherTabs(newResultsTab.Content, newQualifyingTab.Content, newSprintTab.Content)
					}
				} else if id.Col == 2 {
					// For the Circuit column.
					race := schedule.MRData.RaceTable.Races[id.Row-1]
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
					bg.FillColor = theme.Color(theme.ColorNamePrimary)
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
		status.SetText(fmt.Sprintf("%s loaded", title))
		return true
	}

	refresh()

	content := container.NewBorder(nil, status, nil, nil, tableContainer)
	return models.TabData{
		Content: content,
		Refresh: refresh,
	}
}
