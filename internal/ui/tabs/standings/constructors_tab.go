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
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CreateConstructorStandingsTableTab builds the constructor standings tab with a favorite column.
func CreateConstructorStandingsTableTab(parseFunc func([]byte) (string, [][]string, error), year string) models.TabData {
	status := widget.NewLabel("Loading standings...")
	headerLabel := widget.NewLabel("")
	tableContainer := container.NewStack()

	url := processes.BuildStandingsURL(parseFunc, year)

	var refresh func() bool

	// toggleFavoriteConstructor updates the favorite constructor in the config.
	toggleFavoriteConstructor := func(constructorName string) {
		prefs := config.Get()
		if prefs.FavoriteConstructor == constructorName {
			prefs.FavoriteConstructor = ""
		} else {
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

		table := widget.NewTable(
			func() (int, int) {
				if len(rows) == 0 {
					return 0, 0
				}
				return len(rows), len(rows[0])
			},
			func() fyne.CanvasObject {
				return container.NewStack(ui.NewClickableLabel("", nil, false))
			},
			func(id widget.TableCellID, co fyne.CanvasObject) {
				cont, ok := co.(*fyne.Container)
				if !ok {
					return
				}
				cont.Objects = nil

				var cellWidget fyne.CanvasObject
				switch id.Col {
				case 0:
					// Column 0: Position
					cellWidget = ui.NewClickableLabel(rows[id.Row][0], nil, false)
				case 1:
					// Column 1: Favorite star
					rawText := rows[id.Row][2] // linked constructor name
					constructorName := rawText
					if strings.Contains(rawText, "|||") {
						parts := strings.SplitN(rawText, "|||", 2)
						constructorName = parts[0]
					}

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
				case 2:
					// Column 2: Constructor Name with 🌐 icon and link
					cellWidget = processes.MakeClickableConstructorCell(rows[id.Row][2])
				default:
					// Column 3+: Copy remaining columns (e.g., Points)
					cellWidget = ui.NewClickableLabel(rows[id.Row][id.Col], nil, false)
				}
				cont.Add(cellWidget)
				cont.Refresh()
			},
		)

		// Set column widths.
		table.SetColumnWidth(0, 50)  // Position
		table.SetColumnWidth(1, 50)  // Favorite star
		table.SetColumnWidth(2, 180) // Team Name
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
