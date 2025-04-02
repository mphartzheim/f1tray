package schedule

import (
	"fmt"

	"github.com/mphartzheim/f1tray/internal/appstate"
	"github.com/mphartzheim/f1tray/internal/components"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func CreateScheduleTable(state *appstate.AppState, races []ScheduledRace) *widget.Table {
	table := widget.NewTable(
		func() (int, int) {
			return len(races) + 1, 5
		},
		func() fyne.CanvasObject {
			return components.NewClickableCell()
		},
		func(id widget.TableCellID, object fyne.CanvasObject) {
			cell := object.(*components.ClickableCell)
			if id.Row == 0 {
				headers := []string{"Round", "Race", "Circuit", "Date", "Time"}
				cell.Label.SetText(headers[id.Col])
				cell.OnTapped, cell.OnMouseIn, cell.OnMouseOut = nil, nil, nil
				return
			}

			race := races[id.Row-1]
			var text string
			switch id.Col {
			case 0:
				text = race.Round
			case 1:
				text = race.RaceName
			case 2:
				text = race.Circuit.CircuitName
			case 3:
				text = race.Date
			case 4:
				text = race.Time
			}
			cell.Label.SetText(text)

			cell.OnTapped = func() {
				if state.Debug {
					fmt.Printf("Clicked: row=%d, col=%d, value=%s\n", id.Row, id.Col, text)
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
