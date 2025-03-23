package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var debugMode bool

func sessionButton(label, session string, win fyne.Window) *widget.Button {
	return widget.NewButton(label, func() {
		view, err := getSessionResults(win, session)
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		win.SetContent(view)
	})
}

func buildMainMenu(mainWindow fyne.Window) fyne.CanvasObject {
	label := widget.NewLabel("Welcome to F1 Tray!")

	content := container.NewVBox(
		label,
		sessionButton("Race Results", "results", mainWindow),
		sessionButton("Qualifying Results", "qualifying", mainWindow),
		sessionButton("Sprint Qualifying Results", "sprint-qualifying", mainWindow),
		sessionButton("Sprint Results", "sprint", mainWindow),
		sessionButton("Practice 1 Results", "practice/1", mainWindow),
		sessionButton("Practice 2 Results", "practice/2", mainWindow),
		sessionButton("Practice 3 Results", "practice/3", mainWindow),
	)

	if debugMode {
		content.Add(widget.NewLabel("[DEBUG MODE ENABLED]"))
		content.Add(widget.NewButton("Run Debug Test", func() {
			fmt.Println("Debug button clicked")
		}))
	}

	content.Add(widget.NewButton("Quit", func() {
		fyne.CurrentApp().Quit()
	}))

	return content
}

func getSessionResults(mainWindow fyne.Window, session string) (fyne.CanvasObject, error) {
	url := fmt.Sprintf("https://api.jolpi.ca/ergast/f1/current/%s.json", session)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var title string
	headers := []string{"Pos", "Driver", "Constructor", "Time / Status"}
	type RowData struct {
		Pos         string
		Driver      string
		Constructor string
		TimeStatus  string
	}
	var rows []RowData

	if session == "qualifying" {
		var data struct {
			MRData struct {
				RaceTable struct {
					Races []struct {
						RaceName          string `json:"raceName"`
						QualifyingResults []struct {
							Position string `json:"position"`
							Driver   struct {
								GivenName  string `json:"givenName"`
								FamilyName string `json:"familyName"`
							} `json:"Driver"`
							Constructor struct {
								Name string `json:"name"`
							} `json:"Constructor"`
							Q3 string `json:"Q3"`
							Q2 string `json:"Q2"`
							Q1 string `json:"Q1"`
						} `json:"QualifyingResults"`
					} `json:"Races"`
				} `json:"RaceTable"`
			} `json:"MRData"`
		}
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return nil, err
		}
		if len(data.MRData.RaceTable.Races) == 0 {
			dialog.ShowInformation("No Data", fmt.Sprintf("No %s data is available for the current weekend.", session), mainWindow)
			return buildMainMenu(mainWindow), nil
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
	} else if session == "sprint" {
		var data struct {
			MRData struct {
				RaceTable struct {
					Races []struct {
						RaceName      string `json:"raceName"`
						SprintResults []struct {
							Position string `json:"position"`
							Driver   struct {
								GivenName  string `json:"givenName"`
								FamilyName string `json:"familyName"`
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
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return nil, err
		}
		if len(data.MRData.RaceTable.Races) == 0 {
			dialog.ShowInformation("No Data", fmt.Sprintf("No %s data is available for the current weekend.", session), mainWindow)
			return buildMainMenu(mainWindow), nil
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
	} else {
		var data struct {
			MRData struct {
				RaceTable struct {
					Races []struct {
						RaceName string `json:"raceName"`
						Results  []struct {
							Position string `json:"position"`
							Driver   struct {
								GivenName  string `json:"givenName"`
								FamilyName string `json:"familyName"`
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
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return nil, err
		}
		if len(data.MRData.RaceTable.Races) == 0 {
			dialog.ShowInformation("No Data", fmt.Sprintf("No %s data is available for the current weekend.", session), mainWindow)
			return buildMainMenu(mainWindow), nil
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
	backBtn := widget.NewButton("Back to Home", func() {
		mainWindow.SetContent(buildMainMenu(mainWindow))
	})

	return container.NewBorder(
		titleLabel,
		backBtn,
		nil,
		nil,
		container.NewVScroll(table),
	), nil
}

func main() {
	flag.BoolVar(&debugMode, "debug", false, "Enable debug mode")
	flag.Parse()

	f1App := app.NewWithID("com.f1tray.app")

	iconFile, err := os.Open("assets/tray_icon.png")
	if err != nil {
		log.Printf("Failed to load icon: %v", err)
	} else {
		iconData, err := io.ReadAll(iconFile)
		if err != nil {
			log.Printf("Failed to read icon: %v", err)
		} else {
			resource := fyne.NewStaticResource("tray_icon.png", iconData)
			f1App.SetIcon(resource)
		}
		iconFile.Close()
	}

	mainWindow := f1App.NewWindow("F1 Tray")
	mainWindow.Resize(fyne.NewSize(500, 400))
	mainWindow.SetContent(buildMainMenu(mainWindow))

	f1App.Settings().SetTheme(theme.DefaultTheme())

	mainWindow.SetCloseIntercept(func() {
		mainWindow.Hide()
	})

	f1App.Lifecycle().SetOnStopped(func() {
		log.Println("F1 Tray shutting down")
	})

	mainWindow.Show()
	mainWindow.SetMaster()
	f1App.Run()
}
