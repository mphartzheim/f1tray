package results

import (
	"fmt"

	"github.com/mphartzheim/f1tray/internal/appstate"
	"github.com/mphartzheim/f1tray/internal/components"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func CreateRaceResultsTable(state *appstate.AppState, race *RaceResultsEvent) *widget.Table {
	rows := race.Results

	table := widget.NewTable(
		func() (int, int) { return len(rows) + 1, 6 },
		func() fyne.CanvasObject {
			return components.NewClickableCell()
		},
		func(id widget.TableCellID, object fyne.CanvasObject) {
			cell := object.(*components.ClickableCell)

			// Header row
			if id.Row == 0 {
				headers := []string{"Pos", "Driver", "Constructor", "Grid", "Laps", "Time/Status"}
				cell.Label.SetText(headers[id.Col])
				cell.OnTapped = nil
				return
			}

			r := rows[id.Row-1]
			var text string
			switch id.Col {
			case 0:
				text = r.Position
			case 1:
				text = fmt.Sprintf("%s %s", r.Driver.GivenName, r.Driver.FamilyName)
			case 2:
				text = r.Constructor.Name
			case 3:
				text = r.Grid
			case 4:
				text = r.Laps
			case 5:
				if r.Time != nil {
					text = r.Time.Time
				} else {
					text = r.Status
				}
			}
			cell.Label.SetText(text)
		},
	)

	table.SetColumnWidth(0, 60)  // Pos
	table.SetColumnWidth(1, 160) // Driver
	table.SetColumnWidth(2, 160) // Constructor
	table.SetColumnWidth(3, 60)  // Grid
	table.SetColumnWidth(4, 60)  // Laps
	table.SetColumnWidth(5, 160) // Time/Status

	return table
}

func CreateQualifyingResultsTable(state *appstate.AppState, event *QualifyingEvent) *widget.Table {
	rows := event.Results

	table := widget.NewTable(
		func() (int, int) { return len(rows) + 1, 6 },
		func() fyne.CanvasObject {
			return components.NewClickableCell()
		},
		func(id widget.TableCellID, object fyne.CanvasObject) {
			cell := object.(*components.ClickableCell)

			if id.Row == 0 {
				headers := []string{"Pos", "Driver", "Constructor", "Q1", "Q2", "Q3"}
				cell.Label.SetText(headers[id.Col])
				cell.OnTapped = nil
				return
			}

			r := rows[id.Row-1]
			var text string
			switch id.Col {
			case 0:
				text = r.Position
			case 1:
				text = fmt.Sprintf("%s %s", r.Driver.GivenName, r.Driver.FamilyName)
			case 2:
				text = r.Constructor.Name
			case 3:
				text = r.Q1
			case 4:
				text = r.Q2
			case 5:
				text = r.Q3
			}
			cell.Label.SetText(text)
		},
	)

	table.SetColumnWidth(0, 60)  // Pos
	table.SetColumnWidth(1, 160) // Driver
	table.SetColumnWidth(2, 160) // Constructor
	table.SetColumnWidth(3, 100) // Q1
	table.SetColumnWidth(4, 100) // Q2
	table.SetColumnWidth(5, 100) // Q3

	return table
}

func CreateSprintResultsTable(state *appstate.AppState, event *SprintEvent) *widget.Table {
	rows := event.SprintResults

	table := widget.NewTable(
		func() (int, int) { return len(rows) + 1, 6 },
		func() fyne.CanvasObject {
			return components.NewClickableCell()
		},
		func(id widget.TableCellID, object fyne.CanvasObject) {
			cell := object.(*components.ClickableCell)

			if id.Row == 0 {
				headers := []string{"Pos", "Driver", "Constructor", "Grid", "Laps", "Time/Status"}
				cell.Label.SetText(headers[id.Col])
				cell.OnTapped = nil
				return
			}

			r := rows[id.Row-1]
			var text string
			switch id.Col {
			case 0:
				text = r.Position
			case 1:
				text = fmt.Sprintf("%s %s", r.Driver.GivenName, r.Driver.FamilyName)
			case 2:
				text = r.Constructor.Name
			case 3:
				text = r.Grid
			case 4:
				text = r.Laps
			case 5:
				if r.Time != nil {
					text = r.Time.Time
				} else {
					text = r.Status
				}
			}
			cell.Label.SetText(text)
		},
	)

	table.SetColumnWidth(0, 60)  // Pos
	table.SetColumnWidth(1, 160) // Driver
	table.SetColumnWidth(2, 160) // Constructor
	table.SetColumnWidth(3, 60)  // Grid
	table.SetColumnWidth(4, 60)  // Laps
	table.SetColumnWidth(5, 160) // Time/Status

	return table
}
