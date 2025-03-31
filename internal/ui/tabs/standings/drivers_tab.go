package standings

import (
	"fmt"
	"strings"

	"github.com/mphartzheim/f1tray/internal/config"
	"github.com/mphartzheim/f1tray/internal/models"
	"github.com/mphartzheim/f1tray/internal/processes"
	"github.com/mphartzheim/f1tray/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// CreateDriverStandingsTableTab builds the driver standings tab with a clickable star column.
func CreateDriverStandingsTableTab(parseFunc func([]byte) (string, [][]string, error), year string) models.TabData {
	status := widget.NewLabel("Loading standings...")
	headerLabel := widget.NewLabel("")
	tableContainer := container.NewStack() // Container that fills available space

	url := processes.BuildStandingsURL(parseFunc, year)

	// Declare refresh as a variable so it can be referenced by toggleFavorite.
	var refresh func() bool

	// toggleFavorite updates the favorites in the config.
	toggleFavorite := func(driverName string) {
		prefs := config.Get()
		favs := prefs.FavoriteDrivers
		alreadyFav := false
		for _, fav := range favs {
			if fav == driverName {
				alreadyFav = true
				break
			}
		}
		if alreadyFav {
			// Remove from favorites.
			newFavs := []string{}
			for _, fav := range favs {
				if fav != driverName {
					newFavs = append(newFavs, fav)
				}
			}
			prefs.FavoriteDrivers = newFavs
		} else {
			// Only add if there are fewer than 2 favorites.
			if len(favs) < 2 {
				prefs.FavoriteDrivers = append(favs, driverName)
			} else {
				ui.ShowNotification(models.MainWindow, "You can only select up to two favorite drivers.")
				return
			}
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

		table := widget.NewTable(
			func() (int, int) {
				if len(rows) == 0 {
					return 0, 0
				}
				// Expected columns for driver standings:
				// 0: Position, 1: Favorite star, 2: Driver Name, 3: Team, 4: Points.
				return len(rows), len(rows[0])
			},
			// Create each cell as a container that we can update later.
			func() fyne.CanvasObject {
				// Use a container holding a ClickableLabel for consistency.
				return container.NewStack(ui.NewClickableLabel("", nil, false))
			},
			// Update each cell.
			func(id widget.TableCellID, co fyne.CanvasObject) {
				cont, ok := co.(*fyne.Container)
				if !ok {
					return
				}
				cont.Objects = nil

				var cellWidget fyne.CanvasObject
				switch id.Col {
				case 1:
					// Column 1: clickable favorite star.
					driverNameRaw := rows[id.Row][2]
					driverName := driverNameRaw
					if strings.Contains(driverNameRaw, "|||") {
						parts := strings.SplitN(driverNameRaw, "|||", 2)
						driverName = parts[0]
					}
					cellWidget = processes.CreateClickableStar(driverName, toggleFavorite)
				case 2:
					// Column 2: Driver Name.
					cellWidget = processes.MakeClickableDriverCell(rows[id.Row][2])
				default:
					// Other columns (Position, Team, Points) show plain text.
					cellWidget = ui.NewClickableLabel(rows[id.Row][id.Col], nil, false)
				}
				cont.Add(cellWidget)
				cont.Refresh()
			},
		)

		// Set column widths for Driver Standings.
		table.SetColumnWidth(0, 50)  // Position
		table.SetColumnWidth(1, 50)  // Favorite star
		table.SetColumnWidth(2, 180) // Driver Name
		table.SetColumnWidth(3, 100) // Team
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
