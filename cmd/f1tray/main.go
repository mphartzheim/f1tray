package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type QualifyingResultResponse struct {
	MRData struct {
		RaceTable struct {
			Races []struct {
				RaceName          string `json:"raceName"`
				QualifyingResults []struct {
					Position string `json:"position"`
					Driver   struct {
						FamilyName string `json:"familyName"`
						GivenName  string `json:"givenName"`
					} `json:"Driver"`
					Constructor struct {
						Name string `json:"name"`
					} `json:"Constructor"`
					Q1 string `json:"Q1"`
					Q2 string `json:"Q2"`
					Q3 string `json:"Q3"`
				} `json:"QualifyingResults"`
			} `json:"Races"`
		} `json:"RaceTable"`
	} `json:"MRData"`
}

type SprintResultResponse struct {
	MRData struct {
		RaceTable struct {
			Races []struct {
				RaceName      string `json:"raceName"`
				SprintResults []struct {
					Position string `json:"position"`
					Driver   struct {
						FamilyName string `json:"familyName"`
						GivenName  string `json:"givenName"`
					} `json:"Driver"`
					Constructor struct {
						Name string `json:"name"`
					} `json:"Constructor"`
					Time struct {
						Time string `json:"time"`
					} `json:"Time"`
					Status string `json:"status"`
				} `json:"SprintResults"`
			} `json:"Races"`
		} `json:"RaceTable"`
	} `json:"MRData"`
}

type RaceResultResponse struct {
	MRData struct {
		RaceTable struct {
			Races []struct {
				RaceName string `json:"raceName"`
				Results  []struct {
					Position string `json:"position"`
					Driver   struct {
						FamilyName string `json:"familyName"`
						GivenName  string `json:"givenName"`
					} `json:"Driver"`
					Constructor struct {
						Name string `json:"name"`
					} `json:"Constructor"`
					Time struct {
						Time string `json:"time"`
					} `json:"Time"`
					Status string `json:"status"`
				} `json:"Results"`
			} `json:"Races"`
		} `json:"RaceTable"`
	} `json:"MRData"`
}

type ScheduleResponse struct {
	MRData struct {
		RaceTable struct {
			Races []struct {
				Round    string `json:"round"`
				RaceName string `json:"raceName"`
				Date     string `json:"date"`
				Circuit  struct {
					CircuitName string `json:"circuitName"`
					Location    struct {
						Locality string `json:"locality"`
						Country  string `json:"country"`
					} `json:"Location"`
				} `json:"Circuit"`
			} `json:"Races"`
		} `json:"RaceTable"`
	} `json:"MRData"`
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("F1 Viewer")

	resultsTab := createResultsTableTab("https://api.jolpi.ca/ergast/f1/current/last/results.json", parseRaceResults)
	qualifyingTab := createResultsTableTab("https://api.jolpi.ca/ergast/f1/current/last/qualifying.json", parseQualifyingResults)
	sprintTab := createResultsTableTab("https://api.jolpi.ca/ergast/f1/current/last/sprint.json", parseSprintResults)
	scheduleTab := createScheduleTableTab("https://api.jolpi.ca/ergast/f1/current.json", parseSchedule)

	tabs := container.NewAppTabs(
		container.NewTabItem("Race Results", resultsTab),
		container.NewTabItem("Qualifying", qualifyingTab),
		container.NewTabItem("Sprint", sprintTab),
		container.NewTabItem("Schedule", scheduleTab),
	)

	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(900, 600))
	myWindow.ShowAndRun()
}

func createResultsTableTab(url string, parseFunc func([]byte) (string, [][]string, error)) fyne.CanvasObject {
	status := widget.NewLabel("Press 'Load Results' to fetch data.")
	tableContainer := container.NewStack()

	loadButton := widget.NewButton("Load Results", func() {
		resp, err := http.Get(url)
		if err != nil {
			status.SetText(fmt.Sprintf("Fetch error: %v", err))
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			status.SetText(fmt.Sprintf("Read error: %v", err))
			return
		}

		label, rows, err := parseFunc(body)
		if err != nil {
			status.SetText(err.Error())
			return
		}

		table := widget.NewTable(
			func() (int, int) { return len(rows) + 1, 4 },
			func() fyne.CanvasObject {
				label := widget.NewLabel("")
				return container.New(layout.NewStackLayout(), label)
			},
			func(id widget.TableCellID, o fyne.CanvasObject) {
				label := o.(*fyne.Container).Objects[0].(*widget.Label)
				if id.Row == 0 {
					headers := []string{"Pos", "Driver", "Team", "Time/Status"}
					label.SetText(headers[id.Col])
				} else {
					label.SetText(rows[id.Row-1][id.Col])
				}
			},
		)

		table.SetColumnWidth(0, 50)
		table.SetColumnWidth(1, 180)
		table.SetColumnWidth(2, 180)
		table.SetColumnWidth(3, 300)
		table.Resize(fyne.NewSize(600, float32((len(rows)+1)*30)))

		tableContainer.Objects = []fyne.CanvasObject{table}
		tableContainer.Refresh()
		status.SetText(fmt.Sprintf("Results loaded for %s", label))
	})

	return container.NewBorder(loadButton, status, nil, nil, tableContainer)
}

func createScheduleTableTab(url string, parseFunc func([]byte) (string, [][]string, error)) fyne.CanvasObject {
	status := widget.NewLabel("Press 'Load Schedule' to fetch data.")
	tableContainer := container.NewStack()

	var highlightRow int

	loadButton := widget.NewButton("Load Schedule", func() {
		resp, err := http.Get(url)
		if err != nil {
			status.SetText(fmt.Sprintf("Fetch error: %v", err))
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			status.SetText(fmt.Sprintf("Read error: %v", err))
			return
		}

		title, rows, err := parseFunc(body)
		if err != nil {
			status.SetText(err.Error())
			return
		}

		// Determine the row for the current race weekend
		highlightRow = -1
		var schedule ScheduleResponse
		err = json.Unmarshal(body, &schedule)
		if err != nil {
			status.SetText(fmt.Sprintf("Error parsing schedule: %v", err))
			return
		}

		now := time.Now()
		for i, race := range schedule.MRData.RaceTable.Races {
			raceDate, _ := time.Parse("2006-01-02", race.Date)
			if raceDate.After(now) || raceDate.Equal(now) {
				highlightRow = i + 1 // +1 to account for the header row
				break
			}
		}

		table := widget.NewTable(
			func() (int, int) { return len(rows) + 1, 4 },
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
					headers := []string{"Round", "Race Name", "Circuit", "Location (Date)"}
					label.SetText(headers[id.Col])
					bg.Hide()
				} else {
					label.SetText(rows[id.Row-1][id.Col])
					if id.Row == highlightRow {
						bg.FillColor = theme.Color(theme.ColorNamePrimary)
						bg.Show()
					}
					bg.Resize(wrapper.Size())
				}
				wrapper.Refresh()
			},
		)

		table.SetColumnWidth(0, 60)
		table.SetColumnWidth(1, 200)
		table.SetColumnWidth(2, 280)
		table.SetColumnWidth(3, 280)
		table.Resize(fyne.NewSize(820, float32((len(rows)+1)*30)))

		tableContainer.Objects = []fyne.CanvasObject{table}
		tableContainer.Refresh()
		status.SetText(fmt.Sprintf("%s loaded", title))
	})

	return container.NewBorder(loadButton, status, nil, nil, tableContainer)
}

func parseRaceResults(body []byte) (string, [][]string, error) {
	var result RaceResultResponse
	err := json.Unmarshal(body, &result)
	if err != nil {
		return "", nil, fmt.Errorf("JSON error: %v", err)
	}

	if len(result.MRData.RaceTable.Races) == 0 {
		return "", nil, fmt.Errorf("no race data found")
	}

	race := result.MRData.RaceTable.Races[0]
	rows := make([][]string, len(race.Results))
	for i, res := range race.Results {
		timeOrStatus := res.Status
		if res.Time.Time != "" {
			timeOrStatus = res.Time.Time
		}
		rows[i] = []string{
			res.Position,
			fmt.Sprintf("%s %s", res.Driver.GivenName, res.Driver.FamilyName),
			res.Constructor.Name,
			timeOrStatus,
		}
	}
	return race.RaceName, rows, nil
}

func parseSprintResults(body []byte) (string, [][]string, error) {
	var result SprintResultResponse
	err := json.Unmarshal(body, &result)
	if err != nil {
		return "", nil, fmt.Errorf("JSON error: %v", err)
	}

	if len(result.MRData.RaceTable.Races) == 0 {
		return "", nil, fmt.Errorf("no sprint data found")
	}

	race := result.MRData.RaceTable.Races[0]
	rows := make([][]string, len(race.SprintResults))
	for i, res := range race.SprintResults {
		timeOrStatus := res.Status
		if res.Time.Time != "" {
			timeOrStatus = res.Time.Time
		}
		rows[i] = []string{
			res.Position,
			fmt.Sprintf("%s %s", res.Driver.GivenName, res.Driver.FamilyName),
			res.Constructor.Name,
			timeOrStatus,
		}
	}
	return race.RaceName, rows, nil
}

func parseQualifyingResults(body []byte) (string, [][]string, error) {
	var result QualifyingResultResponse
	err := json.Unmarshal(body, &result)
	if err != nil {
		return "", nil, fmt.Errorf("JSON error: %v", err)
	}

	if len(result.MRData.RaceTable.Races) == 0 {
		return "", nil, fmt.Errorf("no qualifying data found")
	}

	race := result.MRData.RaceTable.Races[0]
	rows := make([][]string, len(race.QualifyingResults))
	for i, res := range race.QualifyingResults {
		bestTime := res.Q3
		if bestTime == "" {
			bestTime = res.Q2
		}
		if bestTime == "" {
			bestTime = res.Q1
		}
		rows[i] = []string{
			res.Position,
			fmt.Sprintf("%s %s", res.Driver.GivenName, res.Driver.FamilyName),
			res.Constructor.Name,
			bestTime,
		}
	}
	return race.RaceName, rows, nil
}

func parseSchedule(body []byte) (string, [][]string, error) {
	var result ScheduleResponse
	err := json.Unmarshal(body, &result)
	if err != nil {
		return "", nil, fmt.Errorf("JSON error: %v", err)
	}

	races := result.MRData.RaceTable.Races
	rows := make([][]string, len(races))
	for i, race := range races {
		rows[i] = []string{
			race.Round,
			race.RaceName,
			race.Circuit.CircuitName,
			fmt.Sprintf("%s, %s (%s)", race.Circuit.Location.Locality, race.Circuit.Location.Country, race.Date),
		}
	}
	return "Current Season Schedule", rows, nil
}
