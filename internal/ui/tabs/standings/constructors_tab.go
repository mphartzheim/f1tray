package standings

import (
	"fmt"

	"github.com/mphartzheim/f1tray/internal/models"
	"github.com/mphartzheim/f1tray/internal/processes"
	"github.com/mphartzheim/f1tray/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// CreateConstructorStandingsTableTab builds the constructor standings tab.
func CreateConstructorStandingsTableTab(parseFunc func([]byte) (string, [][]string, error), year string) models.TabData {
	status := widget.NewLabel("Loading standings...")
	headerLabel := widget.NewLabel("")
	tableContainer := container.NewStack()

	url := processes.BuildStandingsURL(parseFunc, year)

	// refresh is specific to the constructor tab.
	refresh := func() bool {
		data, err := processes.FetchData(url)
		if err != nil {
			status.SetText("Failed to fetch standings.")
			return false
		}

		standingsTitle, rows, err := parseFunc(data)
		if err != nil {
			status.SetText("Failed to parse standings")
			return false
		}

		headerLabel.SetText(fmt.Sprintf("Standings for: %s", standingsTitle))

		table := widget.NewTable(
			func() (int, int) {
				if len(rows) == 0 {
					return 0, 0
				}
				// Expected columns for constructor standings:
				// 0: Position, 1: Team Name, 2: Team Nationality, 3: Points.
				return len(rows), len(rows[0])
			},
			// Create each cell as a container with a ClickableLabel.
			func() fyne.CanvasObject {
				return container.NewStack(
					// Create a non-clickable label with empty text.
					// This ensures a consistent cell structure.
					// (Assuming ui.NewClickableLabel is imported from your internal/ui package.)
					ui.NewClickableLabel("", nil, false),
				)
			},
			// Update each cell.
			func(id widget.TableCellID, co fyne.CanvasObject) {
				cont, ok := co.(*fyne.Container)
				if !ok {
					return
				}
				cont.Objects = nil

				var cellWidget fyne.CanvasObject
				if id.Col == 1 {
					// Column 1: Team Name.
					cellWidget = ui.NewClickableLabel(rows[id.Row][1], nil, false)
				} else {
					// Other columns (Position, Nationality, Points) use the corresponding row text.
					cellWidget = ui.NewClickableLabel(rows[id.Row][id.Col], nil, false)
				}

				cont.Add(cellWidget)
				cont.Refresh()
			},
		)

		// Set column widths for Constructor Standings.
		table.SetColumnWidth(0, 50)  // Position
		table.SetColumnWidth(1, 180) // Team Name
		table.SetColumnWidth(2, 100) // Team Nationality
		table.SetColumnWidth(3, 80)  // Points

		tableContainer.Objects = []fyne.CanvasObject{table}
		tableContainer.Refresh()
		status.SetText("Standings loaded")
		return true
	}

	refresh()

	content := container.NewBorder(
		container.NewVBox(headerLabel),
		status,
		nil, nil,
		tableContainer,
	)

	return models.TabData{
		Content: content,
		Refresh: refresh,
	}
}
