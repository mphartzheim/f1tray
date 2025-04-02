package schedule

import (
	"fmt"

	"github.com/mphartzheim/f1tray/internal/appstate"
	"github.com/mphartzheim/f1tray/internal/components"
	"github.com/mphartzheim/f1tray/internal/models"
	"github.com/mphartzheim/f1tray/internal/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func CreateScheduleTable(state *appstate.AppState, races []ScheduledRace) *widget.Table {
	var nextRaceIdx = -1
	now := util.GetNow(state)
	for i, race := range races {
		localTimeStr := util.FormatToLocal(race.Date, race.Time)
		localTime, err := util.ParseDateTime(localTimeStr)
		if err != nil {
			continue
		}
		if localTime.After(now) {
			nextRaceIdx = i
			break
		}
	}

	table := widget.NewTable(
		func() (int, int) {
			return len(races) + 1, 5
		},
		func() fyne.CanvasObject {
			return components.NewClickableCell()
		},
		func(id widget.TableCellID, object fyne.CanvasObject) {
			cell := object.(*components.ClickableCell)
			cell.Label.TextStyle = fyne.TextStyle{}
			cell.OnTapped = nil
			cell.OnMouseIn = nil
			cell.OnMouseOut = nil
			cell.IsPointer = false

			if id.Row == 0 {
				headers := []string{"Round", "Race", "Circuit", "Date", "Time"}
				cell.Label.SetText(headers[id.Col])
				return
			}

			race := races[id.Row-1]
			var text string

			switch id.Col {
			case 0:
				if id.Row-1 == nextRaceIdx {
					cell.Label.TextStyle = fyne.TextStyle{Bold: true}
					text = fmt.Sprintf("%s → Next", race.Round)
				} else {
					text = race.Round
				}
			case 1:
				text = race.RaceName
			case 2:
				text = race.Circuit.CircuitName + " 📍"
				cell.IsPointer = true
				lat := race.Circuit.Location.Lat
				lon := race.Circuit.Location.Long
				if lat != "" && lon != "" {
					cell.OnTapped = func() {
						url := fmt.Sprintf("%s?mlat=%s&mlon=%s#map=17/%s/%s",
							models.MapBaseURL, lat, lon, lat, lon)
						util.OpenWebPage(url)
					}
				}
			case 3:
				text = util.FormatToLocal(race.Date, race.Time)[:10] // Date
			case 4:
				text = util.FormatToLocal(race.Date, race.Time)[11:] // Time
			}

			cell.Label.SetText(text)

			if cell.OnTapped == nil {
				cell.OnTapped = func() {
					if state.Debug {
						fmt.Printf("Clicked: row=%d, col=%d, value=%s\n", id.Row, id.Col, text)
					}
				}
			}
			cell.OnMouseIn = func() {
				if state.Debug {
					fmt.Printf("Mouse in: row=%d, col=%d, value=%s\n", id.Row, id.Col, text)
				}
			}
			cell.OnMouseOut = func() {
				if state.Debug {
					fmt.Printf("Mouse out: row=%d, col=%d, value=%s\n", id.Row, id.Col, text)
				}
			}
		},
	)

	table.SetColumnWidth(0, 80)  // Round
	table.SetColumnWidth(1, 200) // Race
	table.SetColumnWidth(2, 300) // Circuit
	table.SetColumnWidth(3, 100) // Date
	table.SetColumnWidth(4, 100) // Time

	return table
}
