package gui

import (
	"encoding/json"

	"f1tray/internal/api"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func BuildRaceScheduleView(mainWindow fyne.Window, backTo fyne.CanvasObject) (fyne.CanvasObject, error) {
	resp, err := api.GetRaceSchedule()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data api.MRDataContainer[api.RaceSchedule]

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	headers := []string{"Round", "Date", "Race"}

	nextRound, _ := api.GetNextRace()

	table := widget.NewTable(
		func() (int, int) { return len(data.MRData.RaceTable.Races) + 1, len(headers) },
		func() fyne.CanvasObject {
			bg := canvas.NewRectangle(nil)
			label := widget.NewLabel("")
			return container.NewStack(bg, label)
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			wrapper := obj.(*fyne.Container)
			label := wrapper.Objects[1].(*widget.Label)
			bg := wrapper.Objects[0].(*canvas.Rectangle)
			if id.Row == 0 {
				label.SetText(headers[id.Col])
				label.TextStyle = fyne.TextStyle{Bold: true}
				bg.Hide()
				return
			}
			race := data.MRData.RaceTable.Races[id.Row-1]
			if race.Round == nextRound {
				bg.FillColor = HighlightColor()
				bg.Show()
			} else {
				bg.Hide()
			}

			switch id.Col {
			case 0:
				label.SetText(race.Round)
			case 1:
				label.SetText(race.Date)
			case 2:
				label.SetText(race.RaceName)
			}
		},
	)
	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 100)
	table.SetColumnWidth(2, 200)

	title := widget.NewLabel("2025 F1 Race Schedule")
	back := BackButton("Back to Menu", mainWindow, backTo)

	return container.NewBorder(title, back, nil, nil, container.NewVScroll(table)), nil
}
