package standings

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/mphartzheim/f1tray/internal/models"
	"github.com/mphartzheim/f1tray/internal/processes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// CreateStandingsTableTab builds a tab displaying standings fetched from a URL and parsed into a formatted table.
func CreateStandingsTableTab(parseFunc func([]byte) (string, [][]string, error), year string) models.TabData {
	status := widget.NewLabel("Loading standings...")
	headerLabel := widget.NewLabel("")
	tableContainer := container.NewStack()

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
			func() fyne.CanvasObject {
				return widget.NewLabel("")
			},
			func(id widget.TableCellID, cell fyne.CanvasObject) {
				cell.(*widget.Label).SetText(rows[id.Row][id.Col])
			},
		)

		// Set default column widths (adjust as needed).
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
