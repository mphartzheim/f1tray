package tabs

import (
	"encoding/json"
	"fmt"
	"image/color"
	"time"

	"github.com/mphartzheim/f1tray/internal/models"
	"github.com/mphartzheim/f1tray/internal/processes"
	"github.com/mphartzheim/f1tray/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const dtLayout = "2006-01-02 15:04:05"

// CreateUpcomingTab builds a tab showing upcoming race sessions, with a clickable label for map access and a link to F1TV.
func CreateUpcomingTab(state *models.AppState, parseFunc func([]byte) (string, [][]string, error), year string) models.TabData {

	status := widget.NewLabel("Loading upcoming races...")
	nextRaceLabel := ui.NewClickableLabel("Next Race", nil, false)
	tableContainer := container.NewStack()

	url := fmt.Sprintf(models.UpcomingURL, year)

	// Create the button for F1TV.
	button := widget.NewButton("Watch on F1TV", func() {
		if err := processes.OpenWebPage(models.F1tvURL); err != nil {
			status.SetText("Failed to open F1TV URL.")
		}
	})

	// Create a rounded rectangle overlay using canvas.NewRectangle.
	// (This overlay is still used on the watch button.)
	rect := canvas.NewRectangle(color.Transparent)
	rect.StrokeColor = theme.Current().Color(theme.ColorNamePrimary, fyne.CurrentApp().Settings().ThemeVariant())
	rect.StrokeWidth = 4
	rect.CornerRadius = 5

	// Overlay the rectangle on top of the button using container.NewStack.
	watchButton := container.NewStack(button, rect)

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
			createTableCell,
			func(id widget.TableCellID, cell fyne.CanvasObject) {
				updateTableCell(cell, rows, id)
			},
		)

		table.SetColumnWidth(0, 150)
		table.SetColumnWidth(1, 200)
		table.SetColumnWidth(2, 150)
		tableContainer.Objects = []fyne.CanvasObject{table}
		tableContainer.Refresh()
		status.SetText("Upcoming race loaded.")

		var upcomingResp models.UpcomingResponse
		err = json.Unmarshal(data, &upcomingResp)
		if err != nil {
			nextRaceLabel.Text = "Next Race (map unavailable)"
			nextRaceLabel.OnTapped = nil
			nextRaceLabel.Clickable = false
			nextRaceLabel.Refresh()
		} else {
			var upcoming []models.SessionInfo
			now := time.Now()

			for _, race := range upcomingResp.MRData.RaceTable.Races {
				// Build session datetime strings using available session times.
				sessions := make(map[string]string)
				if race.FirstPractice.Date != "" && race.FirstPractice.Time != "" {
					sessions["Practice"] = race.FirstPractice.Date + " " + race.FirstPractice.Time
				}
				if race.SecondPractice.Date != "" && race.SecondPractice.Time != "" {
					sessions["Practice 2"] = race.SecondPractice.Date + " " + race.SecondPractice.Time
				}
				if race.ThirdPractice.Date != "" && race.ThirdPractice.Time != "" {
					sessions["Practice 3"] = race.ThirdPractice.Date + " " + race.ThirdPractice.Time
				}
				if race.Qualifying.Date != "" && race.Qualifying.Time != "" {
					sessions["Qualifying"] = race.Qualifying.Date + " " + race.Qualifying.Time
				}
				if race.Sprint.Date != "" && race.Sprint.Time != "" {
					sessions["Sprint"] = race.Sprint.Date + " " + race.Sprint.Time
				}
				// Always include Race session.
				if race.Date != "" && race.Time != "" {
					sessions["Race"] = race.Date + " " + race.Time
				}

				for sessionType, datetime := range sessions {
					startTime, err := time.Parse(dtLayout, datetime)
					if err != nil {
						continue
					}
					if startTime.After(now) {
						upcoming = append(upcoming, models.SessionInfo{
							Type:      sessionType,
							StartTime: startTime,
							Label:     race.RaceName + " – " + sessionType,
						})
					}
				}
			}
			state.UpcomingSessions = upcoming

			// Configure the Next Race label.
			found := false
			for _, race := range upcomingResp.MRData.RaceTable.Races {
				raceDate, err := time.Parse("2006-01-02", race.Date)
				if err != nil {
					continue
				}
				if raceDate.After(now) || raceDate.Equal(now) {
					nextRaceLabel.Text = fmt.Sprintf("Next Race: %s (%s 🗺️)", race.RaceName, race.Circuit.CircuitName)
					// UpcomingResponse may not have lat/long, so using locality and country as a fallback.
					locality := race.Circuit.Location.Locality
					country := race.Circuit.Location.Country
					mapURL := fmt.Sprintf("%s?locality=%s&country=%s", models.MapBaseURL, locality, country)
					nextRaceLabel.OnTapped = func() {
						if err := processes.OpenWebPage(mapURL); err != nil {
							status.SetText("Failed to open map URL")
						}
					}
					nextRaceLabel.Clickable = true
					nextRaceLabel.Refresh()
					found = true
					break
				}
			}
			if !found {
				nextRaceLabel.Text = "Next Race: Not available"
				nextRaceLabel.OnTapped = nil
				nextRaceLabel.Clickable = false
				nextRaceLabel.Refresh()
			}
		}

		return true
	}

	refresh()

	// Create a horizontal container for the top row.
	// The spacer pushes the button to the far right.
	topRow := container.NewHBox(nextRaceLabel, layout.NewSpacer(), watchButton)

	// Place the status label in a separate container at the bottom.
	bottomContent := container.NewVBox(status)

	// Set up the main content with the table in the center.
	content := container.NewBorder(
		topRow,         // top
		bottomContent,  // bottom
		nil,            // left
		nil,            // right
		tableContainer, // center
	)

	return models.TabData{
		Content: content,
		Refresh: refresh,
	}
}

// createTableCell returns a cell that is a stack containing a background rectangle and a ClickableLabel.
func createTableCell() fyne.CanvasObject {
	cl := ui.NewClickableLabel("", nil, false)
	return container.NewVBox(cl)
}

// updateTableCell sets the cell text and, for certain columns, applies formatting or interactivity.
func updateTableCell(cell fyne.CanvasObject, rows [][]string, id widget.TableCellID) {
	cont, ok := cell.(*fyne.Container)
	if !ok {
		return
	}
	cont.Objects = nil

	var newLabel *ui.ClickableLabel
	switch id.Col {
	case 1:
		// For column 1, reformat the date to a full format.
		text := formatFullDate(rows[id.Row][id.Col])
		newLabel = ui.NewClickableLabel(text, nil, false)
	case 2:
		// Column 2: For the session time.
		text := rows[id.Row][id.Col]
		if processes.IsSessionInProgress(rows[id.Row][0], rows[id.Row][1]) {
			clickableText := fmt.Sprintf("%s 🔴", text)
			newLabel = ui.NewClickableLabel(clickableText, func() {
				processes.OpenWebPage(models.F1tvURL)
			}, true)
			newLabel.SetTextColor(theme.Current().Color(theme.ColorNamePrimary, fyne.CurrentApp().Settings().ThemeVariant()))
		} else {
			newLabel = ui.NewClickableLabel(text, nil, false)
		}
	default:
		newLabel = ui.NewClickableLabel(rows[id.Row][id.Col], nil, false)
	}
	cont.Add(newLabel)
	cont.Refresh()
}

// formatFullDate converts a YYYY-MM-DD string to a full date like "Friday, April 4th 2025".
func formatFullDate(dateStr string) string {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr // fallback to original if parsing fails
	}

	day := t.Day()
	suffix := "th"
	if day%10 == 1 && day != 11 {
		suffix = "st"
	} else if day%10 == 2 && day != 12 {
		suffix = "nd"
	} else if day%10 == 3 && day != 13 {
		suffix = "rd"
	}

	return fmt.Sprintf("%s, %s %d%s %d", t.Weekday(), t.Format("January"), day, suffix, t.Year())
}
