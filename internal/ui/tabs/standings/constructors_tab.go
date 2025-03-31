package standings

import (
	"fmt"

	"github.com/mphartzheim/f1tray/internal/config"
	"github.com/mphartzheim/f1tray/internal/models"
	"github.com/mphartzheim/f1tray/internal/processes"
	"github.com/mphartzheim/f1tray/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CreateConstructorStandingsTableTab builds the constructor standings tab with a favorite column.
func CreateConstructorStandingsTableTab(parseFunc func([]byte) (string, [][]string, error), year string) models.TabData {
	status := widget.NewLabel("Loading standings...")
	headerLabel := widget.NewLabel("")
	tableContainer := container.NewStack()

	url := processes.BuildStandingsURL(parseFunc, year)

	// Declare refresh as a variable so it can be referenced by toggleFavorite.
	var refresh func() bool

	// toggleFavoriteConstructor updates the favorite constructor in the config.
	toggleFavoriteConstructor := func(constructorName string) {
		prefs := config.Get()
		if prefs.FavoriteConstructor == constructorName {
			// Deselect favorite.
			prefs.FavoriteConstructor = ""
		} else {
			// Only one favorite is allowed.
			if prefs.FavoriteConstructor != "" {
				ui.ShowNotification(models.MainWindow, "You can only select one favorite constructor.")
				return
			}
			prefs.FavoriteConstructor = constructorName
		}
		if err := config.SaveConfig(*prefs); err != nil {
			ui.ShowNotification(models.MainWindow, "Failed to save config.")
			return
		}
		refresh()
	}

	refresh = func() bool {
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

		// The original rows have 4 columns:
		//   0: Position, 1: Team Name, 2: Team Nationality, 3: Points.
		// We insert a new favorite column at index 1 so total columns become 5.
		table := widget.NewTable(
			func() (int, int) {
				if len(rows) == 0 {
					return 0, 0
				}
				return len(rows), len(rows[0]) + 1
			},
			// Factory: Create each cell as a container with a ClickableLabel.
			func() fyne.CanvasObject {
				return container.NewStack(ui.NewClickableLabel("", nil, false))
			},
			// Update function.
			func(id widget.TableCellID, co fyne.CanvasObject) {
				cont, ok := co.(*fyne.Container)
				if !ok {
					return
				}
				cont.Objects = nil

				var cellWidget fyne.CanvasObject
				switch id.Col {
				case 0:
					// Column 0: Position (from original rows[i][0]).
					cellWidget = ui.NewClickableLabel(rows[id.Row][0], nil, false)
				case 1:
					// Column 1: Favorite star.
					// Use original team name from rows[i][1] as the constructor name.
					constructorName := rows[id.Row][1]
					star := "☆"
					prefs := config.Get()
					if prefs.FavoriteConstructor == constructorName {
						star = "★"
					}
					cellWidget = ui.NewClickableLabel(star, func() {
						toggleFavoriteConstructor(constructorName)
					}, true)
					cellWidget.(*ui.ClickableLabel).SetTextColor(
						theme.Current().Color(theme.ColorNamePrimary, fyne.CurrentApp().Settings().ThemeVariant()),
					)
				default:
					// For columns 2-4, shift index by -1.
					// Column 2: Originally Team Name (rows[i][1]),
					// Column 3: Originally Team Nationality (rows[i][2]),
					// Column 4: Originally Points (rows[i][3]).
					cellWidget = ui.NewClickableLabel(rows[id.Row][id.Col-1], nil, false)
				}
				cont.Add(cellWidget)
				cont.Refresh()
			},
		)

		// Set column widths.
		table.SetColumnWidth(0, 50)  // Position
		table.SetColumnWidth(1, 50)  // Favorite star
		table.SetColumnWidth(2, 180) // Team Name
		table.SetColumnWidth(3, 100) // Team Nationality
		table.SetColumnWidth(4, 80)  // Points

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
