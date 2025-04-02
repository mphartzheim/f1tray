package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/mphartzheim/f1tray/internal/appstate"
	"github.com/mphartzheim/f1tray/internal/layout"
	"github.com/mphartzheim/f1tray/internal/models"
	"github.com/mphartzheim/f1tray/internal/results"
	"github.com/mphartzheim/f1tray/internal/schedule"
	"github.com/mphartzheim/f1tray/internal/standings"
	"github.com/mphartzheim/f1tray/internal/themes"
	"github.com/mphartzheim/f1tray/internal/upcoming"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	fyneLayout "fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var debugFlag = flag.Bool("debug", false, "Enable debug mode")

type tabLoadResult struct {
	name string
	err  error
}

func main() {
	flag.Parse()
	myApp := app.New()
	myWindow := myApp.NewWindow("F1 Tray Application")
	currentYear := fmt.Sprintf("%d", time.Now().Year())
	state := &appstate.AppState{Window: myWindow, SelectedYear: currentYear, Debug: *debugFlag}

	themeSelect := widget.NewSelect(themes.SortedThemeList(), func(selected string) {
		fyne.CurrentApp().Settings().SetTheme(themes.AvailableThemes()[selected])
	})
	themeSelect.SetSelected("System") // Or load from config
	prefsTab := container.NewVBox(
		container.NewHBox(
			widget.NewLabel("Select Theme:"),
			container.NewStack(themeSelect),
		),
	)

	nextRace, err := upcoming.FetchNextRace()
	if err != nil {
		fmt.Println("Error fetching upcoming race:", err)
		nextRace = &upcoming.NextRace{}
	}

	scheduleContainer := container.NewStack()
	scheduleContainer.Resize(fyne.NewSize(900, 780))
	updateSchedule := func() {
		races, err := schedule.FetchSchedule(state)
		if err != nil {
			fmt.Println("Error fetching schedule:", err)
			races = []schedule.ScheduledRace{}
		}
		scheduleContainer.Objects = []fyne.CanvasObject{
			container.NewStack(schedule.CreateScheduleTable(state, races)),
		}
		scheduleContainer.Refresh()
	}

	// Define countdown binding and label
	countdown := binding.NewString()
	countdownLabel := widget.NewLabelWithData(countdown)
	seasonLabel := widget.NewLabel("Season: ")

	// Build the list of seasons (from current year to 1950)
	currentYearInt := time.Now().Year()
	years := make([]string, 0)
	for y := currentYearInt; y >= 1950; y-- {
		years = append(years, fmt.Sprintf("%d", y))
	}

	upcomingContainer := container.NewStack()
	upcomingContainer.Resize(fyne.NewSize(900, 780))
	updateUpcoming := func() {
		nextRace, err := upcoming.FetchNextRace()
		if err != nil {
			fmt.Println("Error fetching upcoming race:", err)
			upcomingContainer.Objects = []fyne.CanvasObject{
				widget.NewLabel("Failed to load upcoming race."),
			}
		} else {
			upcomingContainer.Objects = []fyne.CanvasObject{
				container.NewStack(upcoming.CreateUpcomingTable(state, nextRace)),
			}
		}
		upcomingContainer.Refresh()
	}

	raceResultsContainer := container.NewStack(widget.NewLabel("Loading race results..."))
	raceResultsContainer.Resize(fyne.NewSize(900, 780))
	updateRaceResults := func() {
		raceURL := fmt.Sprintf(models.RaceURL, state.SelectedYear, "last")
		raceData, err := results.FetchRaceResults(state, raceURL)
		if err != nil {
			fmt.Println("Error fetching race results:", err)
			raceResultsContainer.Objects = []fyne.CanvasObject{
				widget.NewLabel("Failed to load race results."),
			}
		} else {
			raceResultsContainer.Objects = []fyne.CanvasObject{
				results.CreateRaceResultsTable(state, raceData),
			}
		}
		raceResultsContainer.Refresh()
	}

	qualifyingResultsContainer := container.NewStack(widget.NewLabel("Loading qualifying results..."))
	qualifyingResultsContainer.Resize(fyne.NewSize(900, 780))
	updateQualifyingResults := func() {
		qualURL := fmt.Sprintf(models.QualifyingURL, state.SelectedYear, "last")
		qualData, err := results.FetchQualifyingResults(state, qualURL)
		if err != nil {
			fmt.Println("Error fetching qualifying results:", err)
			qualifyingResultsContainer.Objects = []fyne.CanvasObject{
				widget.NewLabel("Failed to load qualifying results."),
			}
		} else {
			qualifyingResultsContainer.Objects = []fyne.CanvasObject{
				results.CreateQualifyingResultsTable(state, qualData),
			}
		}
		qualifyingResultsContainer.Refresh()
	}

	sprintResultsContainer := container.NewStack(widget.NewLabel("Loading sprint results..."))
	sprintResultsContainer.Resize(fyne.NewSize(900, 780))
	updateSprintResults := func() {
		sprintURL := fmt.Sprintf(models.SprintURL, state.SelectedYear, "last")
		sprintData, err := results.FetchSprintResults(state, sprintURL)

		if err != nil {
			fmt.Println("Error fetching sprint results:", err)
			sprintResultsContainer.Objects = []fyne.CanvasObject{
				widget.NewLabel("Failed to load sprint results."),
			}
		} else if sprintData == nil {
			sprintResultsContainer.Objects = []fyne.CanvasObject{
				widget.NewLabel("Not a Sprint event."),
			}
		} else {
			sprintResultsContainer.Objects = []fyne.CanvasObject{
				results.CreateSprintResultsTable(state, sprintData),
			}
		}
		sprintResultsContainer.Refresh()
	}

	driverStandingsContainer := container.NewStack(widget.NewLabel("Loading driver standings..."))
	driverStandingsContainer.Resize(fyne.NewSize(900, 780))
	updateDriverStandings := func() {
		driverURL := fmt.Sprintf(models.DriversStandingsURL, state.SelectedYear)
		data, err := standings.FetchDriverStandings(state, driverURL)

		if err != nil {
			fmt.Println("Error fetching driver standings:", err)
			driverStandingsContainer.Objects = []fyne.CanvasObject{
				widget.NewLabel("Failed to load driver standings."),
			}
		} else {
			driverStandingsContainer.Objects = []fyne.CanvasObject{
				standings.CreateDriverStandingsTable(state, data),
			}
		}
		driverStandingsContainer.Refresh()
	}

	constructorStandingsContainer := container.NewStack(widget.NewLabel("Loading constructor standings..."))
	constructorStandingsContainer.Resize(fyne.NewSize(900, 780))
	updateConstructorStandings := func() {
		constructorURL := fmt.Sprintf(models.ConstructorsStandingsURL, state.SelectedYear)
		data, err := standings.FetchConstructorStandings(state, constructorURL)

		if err != nil {
			fmt.Println("Error fetching constructor standings:", err)
			constructorStandingsContainer.Objects = []fyne.CanvasObject{
				widget.NewLabel("Failed to load constructor standings."),
			}
		} else {
			constructorStandingsContainer.Objects = []fyne.CanvasObject{
				standings.CreateConstructorStandingsTable(state, data),
			}
		}
		constructorStandingsContainer.Refresh()
	}

	onYearSelected := func(newYear string) {
		state.SelectedYear = newYear

		results := make(chan tabLoadResult)

		/* go func() {
			updateSchedule()
			results <- tabLoadResult{name: "Schedule", err: nil}
		}() */
		updateSchedule() // Doesn't like to be in a channel, fyne.Do[AndWait] errors

		go func() {
			updateUpcoming()
			results <- tabLoadResult{name: "Upcoming", err: nil}
		}()

		go func() {
			updateRaceResults()
			results <- tabLoadResult{name: "Race Results", err: nil}
		}()

		go func() {
			updateQualifyingResults()
			results <- tabLoadResult{name: "Qualifying Results", err: nil}
		}()

		go func() {
			updateSprintResults()
			results <- tabLoadResult{name: "Sprint Results", err: nil}
		}()

		go func() {
			updateDriverStandings()
			results <- tabLoadResult{name: "Driver Standings", err: nil}
		}()

		go func() {
			updateConstructorStandings()
			results <- tabLoadResult{name: "Constructor Standings", err: nil}
		}()

		go func() {
			for i := 0; i < 7; i++ {
				res := <-results
				if state.Debug {
					if res.err != nil {
						fmt.Println("❌ Error loading", res.name, ":", res.err)
					} else {
						fmt.Println("✅ Loaded", res.name)
					}
				}
			}
		}()
	}

	// Season dropdown and callback
	seasonSelect := widget.NewSelect(years, func(selected string) {
		if state.Debug {
			fmt.Println("Selected season:", selected)
		}
		state.SelectedYear = selected
		onYearSelected(selected)
	})

	// Set the default selection to current year
	if len(years) > 0 {
		seasonSelect.SetSelected(years[0])
		state.SelectedYear = years[0]
	}

	// Build the top row container
	topRow := container.NewHBox(seasonLabel, seasonSelect, fyneLayout.NewSpacer(), countdownLabel)

	resultsTabs := container.NewAppTabs(
		container.NewTabItem("Race", raceResultsContainer),
		container.NewTabItem("Qualifying", qualifyingResultsContainer),
		container.NewTabItem("Sprint", sprintResultsContainer),
	)

	standingsTabs := container.NewAppTabs(
		container.NewTabItem("Driver", driverStandingsContainer),
		container.NewTabItem("Constructor", constructorStandingsContainer),
	)

	preferencesTabs := container.NewAppTabs(
		container.NewTabItem("Main", prefsTab),
		container.NewTabItem("Notifications", widget.NewLabel("Notification preferences content")),
	)

	outerTabs := container.NewAppTabs(
		container.NewTabItem("Schedule", scheduleContainer),
		container.NewTabItem("Upcoming", upcomingContainer),
		container.NewTabItem("Results", resultsTabs),
		container.NewTabItem("Standings", standingsTabs),
		container.NewTabItem("Preferences", preferencesTabs),
	)

	content := container.NewBorder(topRow, nil, nil, nil, outerTabs)
	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(900, 800))

	// Start live countdown goroutine
	go func() {
		for {
			sessionTime, label := layout.GetNextSessionTime(nextRace)
			if sessionTime.IsZero() {
				countdown.Set("Countdown: no upcoming sessions")
				time.Sleep(10 * time.Second)
				continue
			}

			diff := time.Until(sessionTime)
			weeks := int(diff.Hours()) / 168
			days := (int(diff.Hours()) % 168) / 24
			hours := int(diff.Hours()) % 24
			minutes := int(diff.Minutes()) % 60
			seconds := int(diff.Seconds()) % 60

			units := []struct {
				value  int
				suffix string
			}{
				{weeks, "w"},
				{days, "d"},
				{hours, "h"},
				{minutes, "m"},
				{seconds, "s"},
			}

			var parts []string
			seenNonZero := false
			for _, unit := range units {
				if unit.value > 0 || seenNonZero {
					seenNonZero = true
					parts = append(parts, fmt.Sprintf("%d%s", unit.value, unit.suffix))
				}

				// Always include seconds if everything is zero
				if unit.suffix == "s" && !seenNonZero {
					parts = append(parts, fmt.Sprintf("%d%s", unit.value, unit.suffix))
				}
			}

			formatted := fmt.Sprintf("Countdown to %s: %s", label, strings.Join(parts, " "))
			countdown.Set(formatted)
			time.Sleep(1 * time.Second)
		}
	}()

	myWindow.ShowAndRun()
}
