package tabs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mphartzheim/f1tray/internal/models"
	"github.com/mphartzheim/f1tray/internal/processes"
	"github.com/mphartzheim/f1tray/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CreateScheduleTableTab builds a race schedule tab with interactive circuit links and highlights the upcoming event.
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
				highlightRow = i + 1 // +1 because header is row 0
				break
			}
		}

		// Factory function returns a container with a background rectangle and a clickable label.
		factory := func() fyne.CanvasObject {
			bg := canvas.NewRectangle(nil)
			// Initialize ClickableLabel with Clickable false.
			cl := ui.NewClickableLabel("", nil, false)
			return container.NewStack(bg, cl)
		}

		// Update function for each cell.
		update := func(id widget.TableCellID, obj fyne.CanvasObject) {
			wrapper := obj.(*fyne.Container)
			bg := wrapper.Objects[0].(*canvas.Rectangle)
			cl := wrapper.Objects[1].(*ui.ClickableLabel)

			if id.Row == 0 {
				// Header row: set header text, disable click, and hide background.
				headers := []string{"Round", "Race Name", "Circuit", "Location (Date)"}
				cl.SetText(headers[id.Col])
				cl.OnTapped = nil
				cl.Clickable = false
				bg.Hide()
			} else {
				// Set base text from the row data.
				baseText := rows[id.Row-1][id.Col]
				cl.SetText(baseText)
				// By default, disable clickability.
				cl.OnTapped = nil
				cl.Clickable = false
				// Hide background initially.
				bg.Hide()

				race := schedule.MRData.RaceTable.Races[id.Row-1]
				if id.Col == 1 {
					// Column 1: Race Name. Enable single-click for past races.
					raceDate, err := time.Parse("2006-01-02", race.Date)
					if err != nil {
						fmt.Printf("Error parsing race date: %v\n", err)
					} else if raceDate.Before(now) {
						// Append checkered flag emoji.
						cl.SetText(baseText + " 🏁")
						// Set callback for single-click.
						cl.OnTapped = func() {
							round := rows[id.Row-1][0]
							newResultsTab := CreateResultsTableTab(processes.ParseRaceResults, year, round)
							newQualifyingTab := CreateResultsTableTab(processes.ParseQualifyingResults, year, round)
							newSprintTab := CreateResultsTableTab(processes.ParseSprintResults, year, round)
							newResultsTab.Refresh()
							newQualifyingTab.Refresh()
							newSprintTab.Refresh()
							processes.ReloadOtherTabs(newResultsTab.Content, newQualifyingTab.Content, newSprintTab.Content)
						}
						cl.Clickable = true
					}
				} else if id.Col == 2 {
					// Column 2: Circuit name. Append map emoji and enable single-click.
					cl.SetText(baseText + " 🗺️")
					lat := race.Circuit.Location.Lat
					lon := race.Circuit.Location.Long
					mapURL := fmt.Sprintf("%s?mlat=%s&mlon=%s#map=15/%s/%s", models.MapBaseURL, lat, lon, lat, lon)
					cl.OnTapped = func() {
						if err := ui.OpenWebPage(mapURL); err != nil {
							status.SetText("Failed to open map URL")
						}
					}
					cl.Clickable = true
				}
				// For other columns, leave as non-clickable.

				// Additionally, if this row is the highlighted row, add a red border.
				if id.Row == highlightRow {
					bg.StrokeColor = theme.Current().Color(theme.ColorNamePrimary, fyne.CurrentApp().Settings().ThemeVariant())
					bg.StrokeWidth = 2
					bg.Show()
				} else {
					// Remove border for non-highlighted rows.
					bg.StrokeWidth = 0
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
