package gui

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"f1tray/internal/api"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func isEmptyResponse(raw map[string]interface{}) bool {
	mrData, ok := raw["MRData"].(map[string]interface{})
	if !ok {
		return true
	}
	raceTable, ok := mrData["RaceTable"].(map[string]interface{})
	if !ok {
		return true
	}
	races, ok := raceTable["Races"].([]interface{})
	return !ok || len(races) == 0
}

type RowData struct {
	Pos         string
	Driver      string
	Constructor string
	TimeStatus  string
}

func BuildSessionResults(mainWindow fyne.Window, session string, backTo fyne.CanvasObject) (fyne.CanvasObject, error) {
	var err error
	var resp *http.Response
	var sessionBody []byte

	// Get the session data
	resp, err = api.GetSessionResults(session)
	if err != nil {
		return nil, err
	}
	sessionBody, err = io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	// Try decoding response once
	var raw map[string]interface{}
	_ = json.Unmarshal(sessionBody, &raw)
	if isEmptyResponse(raw) {
		// Fallback to last race endpoint
		fallbackSession := "last/" + session
		resp, err = api.GetSessionResults(fallbackSession)
		if err != nil {
			return nil, err
		}
		sessionBody, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
	}

	var title string
	headers := []string{"Pos", "Driver", "Constructor", "Time / Status"}
	var rows []RowData

	// Define a structure to decode based on session type
	switch session {
	case "qualifying":
		var data api.MRDataContainer[api.QualifyingRace]
		err = json.Unmarshal(sessionBody, &data)
		if err != nil {
			return nil, err
		}
		if len(data.MRData.RaceTable.Races) == 0 {
			dialog.ShowInformation("No Data", fmt.Sprintf("No %s data is available for the current weekend.", session), mainWindow)
			return backTo, nil
		}
		r := data.MRData.RaceTable.Races[0]
		title = "Qualifying - " + r.RaceName
		for _, result := range r.QualifyingResults {
			time := result.Q3
			if time == "" {
				time = result.Q2
			}
			if time == "" {
				time = result.Q1
			}
			rows = append(rows, RowData{
				Pos:         result.Position,
				Driver:      result.Driver.GivenName + " " + result.Driver.FamilyName,
				Constructor: result.Constructor.Name,
				TimeStatus:  time,
			})
		}

	case "sprint":
		var data api.MRDataContainer[api.SprintRace]
		err = json.Unmarshal(sessionBody, &data)
		if err != nil {
			return nil, err
		}
		if len(data.MRData.RaceTable.Races) == 0 {
			dialog.ShowInformation("No Data", fmt.Sprintf("No %s data is available for the current weekend.", session), mainWindow)
			return backTo, nil
		}
		r := data.MRData.RaceTable.Races[0]
		title = "Sprint - " + r.RaceName
		for _, result := range r.SprintResults {
			time := result.Time.Time
			if time == "" {
				time = result.Status
			}
			rows = append(rows, RowData{
				Pos:         result.Position,
				Driver:      result.Driver.GivenName + " " + result.Driver.FamilyName,
				Constructor: result.Constructor.Name,
				TimeStatus:  time,
			})
		}

	default:
		var data api.MRDataContainer[api.ResultsRace]
		err = json.Unmarshal(sessionBody, &data)
		if err != nil {
			return nil, err
		}
		if len(data.MRData.RaceTable.Races) == 0 {
			dialog.ShowInformation("No Data", fmt.Sprintf("No %s data is available for the current weekend.", session), mainWindow)
			return backTo, nil
		}
		r := data.MRData.RaceTable.Races[0]
		title = session + " - " + r.RaceName
		for _, result := range r.Results {
			time := result.Time.Time
			if time == "" {
				time = result.Status
			}
			rows = append(rows, RowData{
				Pos:         result.Position,
				Driver:      result.Driver.GivenName + " " + result.Driver.FamilyName,
				Constructor: result.Constructor.Name,
				TimeStatus:  time,
			})
		}
	}

	table := widget.NewTable(
		func() (int, int) { return len(rows) + 1, len(headers) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			if id.Row == 0 {
				label.SetText(headers[id.Col])
				return
			}
			row := rows[id.Row-1]
			switch id.Col {
			case 0:
				label.SetText(row.Pos)
			case 1:
				label.SetText(row.Driver)
			case 2:
				label.SetText(row.Constructor)
			case 3:
				label.SetText(row.TimeStatus)
			}
		},
	)
	table.SetColumnWidth(0, 40)
	table.SetColumnWidth(1, 140)
	table.SetColumnWidth(2, 100)
	table.SetColumnWidth(3, 100)

	titleLabel := widget.NewLabel(title)
	backBtn := BackButton("Back to Home", mainWindow, backTo)

	return container.NewBorder(
		titleLabel,
		backBtn,
		nil,
		nil,
		container.NewVScroll(table),
	), nil
}
