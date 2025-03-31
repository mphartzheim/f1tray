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
			func() fyne.CanvasObject {
				return container.NewStack(widget.NewLabel(""))
			},
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
					text := rows[id.Row][2]
					if strings.Contains(text, "|||") {
						parts := strings.SplitN(text, "|||", 2)
						displayName := parts[0]
						fallback := strings.TrimSuffix(parts[1], " 👤")
						clickableText := fmt.Sprintf("%s 👤", displayName)
						if slug, ok := models.DriverURLMap[displayName]; ok {
							url := fmt.Sprintf(models.F1DriverBioURL, slug)
							cellWidget = ui.NewClickableLabel(clickableText, func() {
								processes.OpenWebPage(url)
							}, true)
						} else {
							cellWidget = ui.NewClickableLabel(clickableText, func() {
								processes.OpenWebPage(fallback)
							}, true)
						}
					} else {
						cellWidget = widget.NewLabel(text)
					}
				default:
					// Other columns (Position, Team, Points) show plain text.
					cellWidget = widget.NewLabel(rows[id.Row][id.Col])
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
