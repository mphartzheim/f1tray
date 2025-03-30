package standings

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

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

				text := rows[id.Row][id.Col]
				var cellWidget fyne.CanvasObject

				// Check if this cell is in the Driver Name column (index 1).
				if id.Col == 1 {
					// Look for the delimiter indicating a clickable cell.
					if strings.Contains(text, "|||") {
						parts := strings.SplitN(text, "|||", 2)
						displayName := parts[0]
						fallback := parts[1]
						// Remove the trailing emoji from the fallback URL.
						fallback = strings.TrimSuffix(fallback, " 👤")
						// Rebuild the clickable display text.
						clickableText := fmt.Sprintf("%s 👤", displayName)
						// If the driver exists in our custom mapping, use that URL.
						if slug, ok := models.DriverURLMap[displayName]; ok {
							url := fmt.Sprintf(models.F1DriverBioURL, slug)
							cellWidget = ui.NewClickableLabel(clickableText, func() {
								processes.OpenWebPage(url)
							}, true)
						} else {
							// Otherwise, fallback to the API-provided URL.
							cellWidget = ui.NewClickableLabel(clickableText, func() {
								processes.OpenWebPage(fallback)
							}, true)
						}
					} else {
						cellWidget = widget.NewLabel(text)
					}
				} else {
					cellWidget = widget.NewLabel(text)
				}

				cont.Add(cellWidget)
				cont.Refresh()
			},
		)

		table.SetColumnWidth(0, 50)
		table.SetColumnWidth(1, 180)
		table.SetColumnWidth(2, 100)
		if len(rows) > 0 && len(rows[0]) > 3 {
			table.SetColumnWidth(3, 180)
		}

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
