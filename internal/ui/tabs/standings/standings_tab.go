package standings

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/mphartzheim/f1tray/internal/config"
	"github.com/mphartzheim/f1tray/internal/models"
	"github.com/mphartzheim/f1tray/internal/processes"
	"github.com/mphartzheim/f1tray/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// CreateStandingsTableTab builds a tab displaying standings fetched from a URL and parsed into a formatted table.
func CreateStandingsTableTab(parseFunc func([]byte) (string, [][]string, error), year string) models.TabData {
	status := widget.NewLabel("Loading standings...")
	headerLabel := widget.NewLabel("")
	tableContainer := container.NewStack() // Use a container that fills available space

	url := buildStandingsURL(parseFunc, year)

	// Helper function to check whether a driver is a favorite.
	isFavorite := func(favs []string, driverName string) bool {
		for _, fav := range favs {
			if fav == driverName {
				return true
			}
		}
		return false
	}

	// Declare refresh as a variable so it can be referenced by toggleFavorite.
	var refresh func() bool

	// toggleFavorite updates the favorites in the config.
	toggleFavorite := func(driverName string) {
		prefs := config.Get() // retrieve current preferences
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
		// Save the updated config.
		if err := config.SaveConfig(*prefs); err != nil {
			ui.ShowNotification(models.MainWindow, "Failed to save config.")
			return
		}
		// Refresh the table UI.
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
				// Now rows have 5 columns.
				return len(rows), len(rows[0])
			},
			// Create each cell as a container that we can update later.
			func() fyne.CanvasObject {
				return container.NewStack(widget.NewLabel(""))
			},
			// Update each cell based on its row and column.
			func(id widget.TableCellID, co fyne.CanvasObject) {
				cont, ok := co.(*fyne.Container)
				if !ok {
					return
				}
				cont.Objects = nil

				var cellWidget fyne.CanvasObject
				// Column 1 is the Favorite column.
				if id.Col == 1 {
					// The driver name is in column 2.
					driverNameRaw := rows[id.Row][2]
					driverName := driverNameRaw
					if strings.Contains(driverNameRaw, "|||") {
						parts := strings.SplitN(driverNameRaw, "|||", 2)
						driverName = parts[0]
					}
					star := "☆"
					if isFavorite(config.Get().FavoriteDrivers, driverName) {
						star = "★"
					}
					cellWidget = ui.NewClickableLabel(star, func() {
						toggleFavorite(driverName)
					}, true)
				} else if id.Col == 2 {
					// Driver Name column.
					text := rows[id.Row][2]
					// Check for clickable driver URL indicator.
					if strings.Contains(text, "|||") {
						parts := strings.SplitN(text, "|||", 2)
						displayName := parts[0]
						fallback := parts[1]
						fallback = strings.TrimSuffix(fallback, " 👤")
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
				} else {
					// For other columns, simply create a label.
					cellWidget = widget.NewLabel(rows[id.Row][id.Col])
				}

				cont.Add(cellWidget)
				cont.Refresh()
			},
		)

		// Set updated column widths.
		table.SetColumnWidth(0, 50)  // Position
		table.SetColumnWidth(1, 50)  // Favorite
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

// buildStandingsURL builds the URL for standings data based on the provided parse function.
func buildStandingsURL(parseFunc func([]byte) (string, [][]string, error), year string) string {
	funcName := runtime.FuncForPC(reflect.ValueOf(parseFunc).Pointer()).Name()

	if strings.HasSuffix(funcName, "ParseDriverStandings") {
		return fmt.Sprintf(models.DriversStandingsURL, year)
	} else if strings.HasSuffix(funcName, "ParseConstructorStandings") {
		return fmt.Sprintf(models.ConstructorsStandingsURL, year)
	}
	return ""
}
