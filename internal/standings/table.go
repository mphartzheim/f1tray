package standings

import (
	"fmt"

	"github.com/mphartzheim/f1tray/internal/appstate"
	"github.com/mphartzheim/f1tray/internal/components"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func CreateDriverStandingsTable(state *appstate.AppState, standings []DriverStandingItem) *widget.Table {
	table := widget.NewTable(
		func() (int, int) {
			return len(standings) + 1, 5 // header + rows, 5 columns
		},
		func() fyne.CanvasObject {
			return components.NewClickableCell()
		},
		func(id widget.TableCellID, object fyne.CanvasObject) {
			cell := object.(*components.ClickableCell)

			if id.Row == 0 {
				headers := []string{"Pos", "Driver", "Constructor", "Points", "Wins"}
				cell.Label.SetText(headers[id.Col])
				cell.OnTapped = nil
				return
			}

			driver := standings[id.Row-1]
			var text string
			switch id.Col {
			case 0:
				text = driver.Position
			case 1:
				text = fmt.Sprintf("%s %s", driver.Driver.GivenName, driver.Driver.FamilyName)
			case 2:
				if len(driver.Constructors) > 0 {
					text = driver.Constructors[0].Name
				} else {
					text = "N/A"
				}
			case 3:
				text = driver.Points
			case 4:
				text = driver.Wins
			}

			cell.Label.SetText(text)

			cell.OnTapped = func() {
				fmt.Printf("Driver standing cell clicked: %s\n", text)
			}
		},
	)

	table.SetColumnWidth(0, 60)  // Pos
	table.SetColumnWidth(1, 180) // Driver
	table.SetColumnWidth(2, 160) // Constructor
	table.SetColumnWidth(3, 80)  // Points
	table.SetColumnWidth(4, 80)  // Wins
	table.Resize(fyne.NewSize(900, 700))

	return table
}

func CreateConstructorStandingsTable(state *appstate.AppState, standings []ConstructorStandingPosition) *widget.Table {
	table := widget.NewTable(
		func() (int, int) {
			return len(standings) + 1, 4
		},
		func() fyne.CanvasObject {
			return components.NewClickableCell()
		},
		func(id widget.TableCellID, object fyne.CanvasObject) {
			cell := object.(*components.ClickableCell)

			if id.Row == 0 {
				headers := []string{"Pos", "Constructor", "Nationality", "Points"}
				cell.Label.SetText(headers[id.Col])
				cell.OnTapped = nil
				return
			}

			team := standings[id.Row-1]
			var text string
			switch id.Col {
			case 0:
				text = team.Position
			case 1:
				text = team.Constructor.Name
			case 2:
				text = team.Constructor.Nationality
			case 3:
				text = team.Points
			}

			cell.Label.SetText(text)
			cell.OnTapped = func() {
				fmt.Printf("Constructor standings clicked: %s\n", text)
			}
		},
	)

	table.SetColumnWidth(0, 60)  // Pos
	table.SetColumnWidth(1, 200) // Constructor
	table.SetColumnWidth(2, 120) // Nationality
	table.SetColumnWidth(3, 100) // Points

	table.Resize(fyne.NewSize(900, 700))
	return table
}
