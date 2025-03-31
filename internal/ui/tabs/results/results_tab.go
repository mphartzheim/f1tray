package results

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

// CreateResultsTableTab builds a tab displaying race results fetched from a URL and parsed into a formatted table.
func CreateResultsTableTab(parseFunc func([]byte) (string, [][]string, error), year string, round string) models.TabData {
	status := widget.NewLabel("Loading results...")
	// Use ClickableLabel for the race name display (non-clickable).
	raceNameLabel := ui.NewClickableLabel("", nil, false)
	tableContainer := container.NewStack()

	url := buildResultsURL(parseFunc, year, round)

	// Declare refresh as a variable so it can be referenced inside toggleFavorite.
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
		// Save updated config.
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
			status.SetText("Failed to fetch results.")
			return false
		}

		// Get the function name for later use.
		funcName := runtime.FuncForPC(reflect.ValueOf(parseFunc).Pointer()).Name()

		// Parse the data.
		raceName, rows, err := parseFunc(data)
		// Handle missing data gracefully for sprint/qualifying events.
		if err != nil {
			if strings.HasSuffix(funcName, "ParseSprintResults") && strings.Contains(err.Error(), "no sprint data found") {
				raceNameLabel.Text = "Not a sprint race event"
				raceNameLabel.Refresh()
				tableContainer.Objects = nil
				tableContainer.Refresh()
				status.SetText("Results loaded")
				return true
			} else if strings.HasSuffix(funcName, "ParseQualifyingResults") && strings.Contains(err.Error(), "no qualifying data found") {
				raceNameLabel.Text = "No data available on Jolpica API"
				raceNameLabel.Refresh()
				tableContainer.Objects = nil
				tableContainer.Refresh()
				status.SetText("Results loaded.")
				return true
			}
			status.SetText("Failed to parse results")
			return false
		}

		raceNameLabel.Text = fmt.Sprintf("Results for: %s", raceName)
		raceNameLabel.Refresh()

		// Create a table:
		// Now each row has 5 columns:
		//   [0]: Position, [1]: Favorite, [2]: Driver, [3]: Team, [4]: Time/Status.
		table := widget.NewTable(
			func() (int, int) {
				if len(rows) == 0 {
					return 0, 0
				}
				return len(rows), len(rows[0])
			},
			// Create each cell as a container that we can update later.
			func() fyne.CanvasObject {
				// Use ClickableLabel with no tap callback.
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
				// Column 0: Position.
				case 0:
					cellWidget = ui.NewClickableLabel(rows[id.Row][0], nil, false)
				// Column 1: Favorite column.
				case 1:
					// The driver name is stored in column 2.
					driverNameRaw := rows[id.Row][2]
					driverName := driverNameRaw
					if strings.Contains(driverNameRaw, "|||") {
						parts := strings.SplitN(driverNameRaw, "|||", 2)
						driverName = parts[0]
					}
					cellWidget = processes.CreateClickableStar(driverName, toggleFavorite)
				// Column 2: Driver column.
				case 2:
					cellWidget = processes.MakeClickableDriverCell(rows[id.Row][2])
				// Column 3: Team.
				case 3:
					cellWidget = ui.NewClickableLabel(rows[id.Row][3], nil, false)
				// Column 4: Time/Status.
				case 4:
					cellWidget = ui.NewClickableLabel(rows[id.Row][4], nil, false)
				}
				cont.Add(cellWidget)
				cont.Refresh()
			},
		)

		// Set proper column widths.
		table.SetColumnWidth(0, 50)  // Position
		table.SetColumnWidth(1, 50)  // Favorite
		table.SetColumnWidth(2, 180) // Driver
		table.SetColumnWidth(3, 180) // Team
		table.SetColumnWidth(4, 300) // Time/Status

		tableContainer.Objects = []fyne.CanvasObject{table}
		tableContainer.Refresh()
		status.SetText("Results loaded")
		return true
	}

	refresh()

	content := container.NewBorder(
		container.NewVBox(raceNameLabel), // Top
		status,                           // Bottom
		nil, nil,
		tableContainer,
	)

	return models.TabData{
		Content: content,
		Refresh: refresh,
	}
}

// buildResultsURL builds the correct URL based on the parse function provided.
func buildResultsURL(parseFunc func([]byte) (string, [][]string, error), year string, round string) string {
	funcName := runtime.FuncForPC(reflect.ValueOf(parseFunc).Pointer()).Name()

	if strings.HasSuffix(funcName, "ParseRaceResults") {
		return fmt.Sprintf(models.RaceURL, year, round)
	} else if strings.HasSuffix(funcName, "ParseQualifyingResults") {
		return fmt.Sprintf(models.QualifyingURL, year, round)
	} else if strings.HasSuffix(funcName, "ParseSprintResults") {
		return fmt.Sprintf(models.SprintURL, year, round)
	}
	return ""
}
