package ui

import (
	"encoding/json"
	"f1tray/internal/config"
	"f1tray/internal/models"
	"f1tray/internal/processes"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func CreateResultsTableTab(url string, parseFunc func([]byte) (string, [][]string, error)) models.TabData {
	status := widget.NewLabel("Loading results...")
	raceNameLabel := widget.NewLabel("")
	tableContainer := container.NewStack()

	refresh := func() bool {
		data, changed, err := processes.FetchData(url)
		if err != nil {
			status.SetText("Failed to fetch results.")
			return false
		}
		if !changed {
			return false
		}

		raceName, rows, err := parseFunc(data)
		if err != nil {
			status.SetText("Failed to parse results.")
			return false
		}

		raceNameLabel.SetText(fmt.Sprintf("Results for: %s", raceName))
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

		// Restore proper column widths for results view
		table.SetColumnWidth(0, 50)  // Pos
		table.SetColumnWidth(1, 180) // Driver
		table.SetColumnWidth(2, 180) // Team
		table.SetColumnWidth(3, 300) // Time/Status

		tableContainer.Objects = []fyne.CanvasObject{table}
		tableContainer.Refresh()
		status.SetText("Results loaded.")
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

func CreateUpcomingTab(url string, parseFunc func([]byte) (string, [][]string, error)) models.TabData {
	status := widget.NewLabel("Loading upcoming races...")
	titleLabel := widget.NewLabel("")
	tableContainer := container.NewStack()

	// Create the "Watch on F1TV" button.
	watchButton := widget.NewButton("Watch on F1TV", func() {
		if err := OpenWebPage(models.F1tvURL); err != nil {
			status.SetText("Failed to open F1TV URL.")
		}
	})

	refresh := func() bool {
		data, changed, err := processes.FetchData(url)
		if err != nil {
			status.SetText("Failed to fetch upcoming data.")
			return false
		}
		if !changed {
			return false
		}

		title, rows, err := parseFunc(data)
		if err != nil {
			status.SetText("Failed to parse upcoming data.")
			return false
		}

		titleLabel.SetText(title)
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

		table.SetColumnWidth(0, 150)
		table.SetColumnWidth(1, 150)
		table.SetColumnWidth(2, 150)
		tableContainer.Objects = []fyne.CanvasObject{table}
		tableContainer.Refresh()
		status.SetText("Upcoming race loaded.")
		return true
	}

	refresh()

	// Combine the status label and the watch button in a vertical layout.
	bottomContent := container.NewVBox(status, watchButton)

	content := container.NewBorder(
		container.NewVBox(titleLabel), // Top
		bottomContent,                 // Bottom
		nil, nil,
		tableContainer,
	)

	return models.TabData{
		Content: content,
		Refresh: refresh,
	}
}

func CreateScheduleTableTab(url string, parseFunc func([]byte) (string, [][]string, error)) models.TabData {
	status := widget.NewLabel("Loading schedule...")
	tableContainer := container.NewStack()

	refresh := func() bool {
		data, changed, err := processes.FetchData(url)
		if err != nil {
			status.SetText("Failed to fetch schedule.")
			return false
		}
		if !changed {
			return false
		}

		title, rows, err := parseFunc(data)
		if err != nil {
			status.SetText("Failed to parse schedule.")
			return false
		}

		// Unmarshal into our struct that now includes lat/long as strings.
		var schedule models.ScheduleResponse
		err = json.Unmarshal(data, &schedule)
		if err != nil {
			status.SetText(fmt.Sprintf("Error parsing schedule JSON: %v", err))
			return false
		}

		now := time.Now()
		highlightRow := -1
		for i, race := range schedule.MRData.RaceTable.Races {
			raceDate, _ := time.Parse("2006-01-02", race.Date)
			if raceDate.After(now) || raceDate.Equal(now) {
				highlightRow = i + 1 // +1 because headers take row 0
				break
			}
		}

		// Factory function returns a container with a background rectangle and a clickable label.
		factory := func() fyne.CanvasObject {
			bg := canvas.NewRectangle(nil)
			cl := NewClickableLabel("", nil)
			return container.NewStack(bg, cl)
		}

		// Update function for each cell.
		update := func(id widget.TableCellID, obj fyne.CanvasObject) {
			wrapper := obj.(*fyne.Container)
			bg := wrapper.Objects[0].(*canvas.Rectangle)
			// The clickable label we defined in ui.
			cl := wrapper.Objects[1].(*ClickableLabel)

			if id.Row == 0 {
				// Header row.
				headers := []string{"Round", "Race Name", "Circuit", "Location (Date)"}
				cl.SetText(headers[id.Col])
				cl.OnDoubleTapped = nil // No click action for headers.
				bg.Hide()
			} else {
				// Data rows.
				cl.SetText(rows[id.Row-1][id.Col])
				// For the Circuit column (column index 2), set the click callback.
				if id.Col == 2 {
					// Get the corresponding race data.
					race := schedule.MRData.RaceTable.Races[id.Row-1]
					lat := race.Circuit.Location.Lat
					lon := race.Circuit.Location.Long
					// Build the OpenStreetMap URL. We'll use a default zoom level of 15.
					mapURL := fmt.Sprintf("%s?mlat=%s&mlon=%s#map=15/%s/%s", models.MapBaseURL, lat, lon, lat, lon)
					cl.OnDoubleTapped = func() {
						if err := OpenWebPage(mapURL); err != nil {
							status.SetText("Failed to open map URL")
						}
					}
				} else {
					// For all other columns, remove any tap handler.
					cl.OnDoubleTapped = nil
				}

				// Highlight the row if needed.
				if id.Row == highlightRow {
					bg.FillColor = theme.Color(theme.ColorNamePrimary)
					bg.Show()
				} else {
					bg.Hide()
				}
				bg.Resize(wrapper.Size())
			}
			wrapper.Refresh()
		}

		table := widget.NewTable(
			func() (int, int) { return len(rows) + 1, 4 },
			factory,
			update,
		)

		table.SetColumnWidth(0, 60)
		table.SetColumnWidth(1, 200)
		table.SetColumnWidth(2, 280)
		table.SetColumnWidth(3, 280)
		table.Resize(fyne.NewSize(820, float32((len(rows)+1)*30)))

		tableContainer.Objects = []fyne.CanvasObject{table}
		tableContainer.Refresh()
		status.SetText(fmt.Sprintf("%s loaded", title))
		return true
	}

	refresh()

	content := container.NewBorder(nil, status, nil, nil, tableContainer)
	return models.TabData{
		Content: content,
		Refresh: refresh,
	}
}

func CreatePreferencesTab(currentPrefs config.Preferences, onSave func(config.Preferences)) fyne.CanvasObject {
	isExit := currentPrefs.CloseBehavior == "exit"
	closeCheckbox := widget.NewCheck("Close on exit?", func(checked bool) {
		if checked {
			currentPrefs.CloseBehavior = "exit"
		} else {
			currentPrefs.CloseBehavior = "minimize"
		}
		onSave(currentPrefs)
	})
	closeCheckbox.SetChecked(isExit)

	hideCheckbox := widget.NewCheck("Hide on open?", func(checked bool) {
		currentPrefs.HideOnOpen = checked
		onSave(currentPrefs)
	})
	hideCheckbox.SetChecked(currentPrefs.HideOnOpen)

	soundCheckbox := widget.NewCheck("Enable sounds?", func(checked bool) {
		currentPrefs.EnableSound = checked
		onSave(currentPrefs)
	})
	soundCheckbox.SetChecked(currentPrefs.EnableSound)

	testButton := widget.NewButton("Test", func() {
		processes.PlayNotificationSound()
	})

	soundRow := container.NewHBox(
		soundCheckbox,
		testButton,
	)

	debugCheckbox := widget.NewCheck("Debug Mode?", func(checked bool) {
		currentPrefs.DebugMode = checked
		onSave(currentPrefs)
	})
	debugCheckbox.SetChecked(currentPrefs.DebugMode)

	return container.NewVBox(
		closeCheckbox,
		hideCheckbox,
		soundRow,
		debugCheckbox,
	)
}
