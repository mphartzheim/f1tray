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
	"fyne.io/fyne/v2/theme"
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

// createTableCell returns a cell that is a stack containing a background rectangle and a label.
func createTableCell() fyne.CanvasObject {
	bg := canvas.NewRectangle(color.Transparent)
	lbl := widget.NewLabel("")
	return container.NewStack(bg, lbl)
}

// updateTableCell sets the cell text and, for the date cell (column 2), reformats the date.
func updateTableCell(cell fyne.CanvasObject, rows [][]string, id widget.TableCellID) {
	cont := cell.(*fyne.Container)
	bg := cont.Objects[0].(*canvas.Rectangle)
	lbl := cont.Objects[1].(*widget.Label)

	// Get the original text from the row.
	text := rows[id.Row][id.Col]

	// For Column 2, reformat the date to a full format.
	if id.Col == 1 {
		text = formatFullDate(text)
	} else if id.Col == 2 {
		// If the session is active, add the 🔴 icon.
		if processes.IsSessionInProgress(rows[id.Row][0], rows[id.Row][1]) {
			text += " 🔴"
			bg.StrokeColor = theme.Current().Color(theme.ColorNamePrimary, fyne.CurrentApp().Settings().ThemeVariant())
			bg.StrokeWidth = 2
			bg.Show()
		} else {
			bg.StrokeWidth = 0
			bg.Hide()
		}
	}

	lbl.SetText(text)
	bg.Refresh()
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
