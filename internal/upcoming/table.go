package upcoming

import (
	"strings"
	"time"

	"github.com/mphartzheim/f1tray/internal/appstate"
	"github.com/mphartzheim/f1tray/internal/components"
	"github.com/mphartzheim/f1tray/internal/models"
	"github.com/mphartzheim/f1tray/internal/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func CreateUpcomingTable(state *appstate.AppState, race *NextRace) *widget.Table {
	rows := [][]string{
		{"Round", race.Round},
		{"Race", race.RaceName},
		{"Circuit", race.Circuit.CircuitName},
		{"Locality", race.Circuit.Location.Locality},
		{"Country", race.Circuit.Location.Country},
	}

	if race.FirstPractice != nil {
		rows = append(rows, []string{"First Practice", util.FormatToLocal(race.FirstPractice.Date, race.FirstPractice.Time)})
	}
	if race.SecondPractice != nil {
		rows = append(rows, []string{"Second Practice", util.FormatToLocal(race.SecondPractice.Date, race.SecondPractice.Time)})
	}
	if race.ThirdPractice != nil {
		rows = append(rows, []string{"Third Practice", util.FormatToLocal(race.ThirdPractice.Date, race.ThirdPractice.Time)})
	}
	if race.Sprint != nil {
		rows = append(rows, []string{"Sprint", util.FormatToLocal(race.Sprint.Date, race.Sprint.Time)})
	}
	if race.Qualifying != nil {
		rows = append(rows, []string{"Qualifying", util.FormatToLocal(race.Qualifying.Date, race.Qualifying.Time)})
	}

	rows = append(rows, []string{"Race", util.FormatToLocal(race.Date, race.Time)})

	table := widget.NewTable(
		func() (int, int) { return len(rows), 2 },
		func() fyne.CanvasObject {
			return components.NewClickableCell()
		},
		func(id widget.TableCellID, object fyne.CanvasObject) {
			cell := object.(*components.ClickableCell)
			label := rows[id.Row][0]
			value := rows[id.Row][1]

			displayText := label
			cell.IsPointer = false

			if id.Col == 1 {
				displayText = value

				if isSessionRow(label) {
					if sessionTime, err := util.ParseDateTime(value); err == nil {
						now := util.GetNow(state)
						diff := sessionTime.Sub(now)

						switch {
						case !now.Before(sessionTime) && diff > -1*time.Hour:
							displayText += "    🔴 LIVE"
							cell.Label.TextStyle = fyne.TextStyle{Bold: true}
							cell.IsPointer = true
						case diff <= 30*time.Minute && diff > 0:
							cell.Label.TextStyle = fyne.TextStyle{Italic: true}
						}
					}
				}
			}
			cell.Label.SetText(displayText)
			cell.OnTapped = func() {
				if id.Col == 1 && isSessionRow(label) && strings.Contains(displayText, "🔴 LIVE") {
					if models.F1tvURL != "" {
						util.OpenWebPage(models.F1tvURL)
					}
				}
			}
		},
	)

	table.SetColumnWidth(0, 150)
	table.SetColumnWidth(1, 400)

	return table
}

func isSessionRow(label string) bool {
	switch label {
	case "First Practice", "Second Practice", "Third Practice", "Sprint", "Qualifying", "Race":
		return true
	default:
		return false
	}
}
