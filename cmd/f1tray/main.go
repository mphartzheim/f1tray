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
	"fyne.io/fyne/v2/widget"
)

var debugMode bool

type RaceWeekend struct {
	MRData struct {
		RaceTable struct {
			Races []struct {
				RaceName string `json:"raceName"`
				Circuit  struct {
					CircuitName string `json:"circuitName"`
					Location    struct {
						Locality string `json:"locality"`
						Country  string `json:"country"`
					} `json:"Location"`
				} `json:"Circuit"`
				Date string `json:"date"`
			} `json:"Races"`
		} `json:"RaceTable"`
	} `json:"MRData"`
}

func getRaceWeekendInfo() (string, error) {
	resp, err := http.Get("https://api.jolpi.ca/ergast/f1/current.json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data RaceWeekend
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	if len(data.MRData.RaceTable.Races) == 0 {
		return "No race data available.", nil
	}

	race := data.MRData.RaceTable.Races[len(data.MRData.RaceTable.Races)-1]
	info := fmt.Sprintf("Next Race: %s\nCircuit: %s\nLocation: %s, %s\nDate: %s",
		race.RaceName,
		race.Circuit.CircuitName,
		race.Circuit.Location.Locality,
		race.Circuit.Location.Country,
		race.Date)

	return info, nil
}

func main() {
	// Parse debug flag
	flag.BoolVar(&debugMode, "debug", false, "Enable debug mode")
	flag.Parse()

	// Create app
	f1App := app.NewWithID("com.f1tray.app")
	// Follow system theme automatically (no need to set manually)

	// Load window icon
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

	// Create window
	mainWindow := f1App.NewWindow("F1 Tray")
	mainWindow.Resize(fyne.NewSize(400, 300))

	// Main content
	content := container.NewVBox(
		widget.NewLabel("Welcome to F1 Tray!"),
		widget.NewButton("Fetch Current Race Info", func() {
			info, err := getRaceWeekendInfo()
			if err != nil {
				dialog.ShowError(err, mainWindow)
				return
			}
			dialog.ShowInformation("Race Info", info, mainWindow)
		}),
	)

	if debugMode {
		content.Add(widget.NewLabel("[DEBUG MODE ENABLED]"))
		content.Add(widget.NewButton("Run Debug Test", func() {
			fmt.Println("Debug button clicked")
		}))
	}

	mainWindow.SetContent(content)
	mainWindow.ShowAndRun()
}
