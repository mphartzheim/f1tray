package main

import (
	"flag"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	fyneLayout "fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/mphartzheim/f1tray/internal/appstate"
	"github.com/mphartzheim/f1tray/internal/layout"
	"github.com/mphartzheim/f1tray/internal/themes"
	"github.com/mphartzheim/f1tray/internal/uiupdaters"
	"github.com/mphartzheim/f1tray/internal/upcoming"
	"github.com/mphartzheim/f1tray/internal/userconfig"
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

	cfg, _ := userconfig.Load()
	themeSelect := widget.NewSelect(themes.SortedThemeList(), func(selected string) {
		fyne.CurrentApp().Settings().SetTheme(themes.AvailableThemes()[selected])
		cfg.SelectedTheme = selected
		_ = userconfig.Save(cfg)
	})
	themeSelect.SetSelected(cfg.SelectedTheme)
	fyne.CurrentApp().Settings().SetTheme(themes.AvailableThemes()[cfg.SelectedTheme])

	prefsTab := container.NewVBox(
		container.NewHBox(
			widget.NewLabel("Select Theme:"),
			container.NewStack(themeSelect),
		),
	)

	nextRace, err := upcoming.FetchNextRace()
	if err != nil {
		fmt.Println("Error fetching next race:", err)
		nextRace = &upcoming.NextRace{}
	}

	scheduleContainer := container.NewStack()
	scheduleContainer.Resize(fyne.NewSize(900, 780))
	updateSchedule := func() {
		obj, err := uiupdaters.FetchScheduleUI(state)
		if err != nil {
			scheduleContainer.Objects = []fyne.CanvasObject{widget.NewLabel("Failed to load.")}
		} else {
			scheduleContainer.Objects = []fyne.CanvasObject{obj}
		}
		scheduleContainer.Refresh()
	}

	upcomingContainer := container.NewStack()
	upcomingContainer.Resize(fyne.NewSize(900, 780))
	updateUpcoming := func() {
		obj, err := uiupdaters.FetchUpcomingUI(state)
		if err != nil {
			upcomingContainer.Objects = []fyne.CanvasObject{widget.NewLabel("Failed to load.")}
		} else {
			upcomingContainer.Objects = []fyne.CanvasObject{obj}
		}
		upcomingContainer.Refresh()
	}

	raceResultsContainer := container.NewStack()
	raceResultsContainer.Resize(fyne.NewSize(900, 780))
	updateRaceResults := func() {
		obj, err := uiupdaters.FetchRaceResultsUI(state)
		if err != nil {
			raceResultsContainer.Objects = []fyne.CanvasObject{widget.NewLabel("Failed to load.")}
		} else {
			raceResultsContainer.Objects = []fyne.CanvasObject{obj}
		}
		raceResultsContainer.Refresh()
	}

	qualifyingResultsContainer := container.NewStack()
	qualifyingResultsContainer.Resize(fyne.NewSize(900, 780))
	updateQualifyingResults := func() {
		obj, err := uiupdaters.FetchQualifyingResultsUI(state)
		if err != nil {
			qualifyingResultsContainer.Objects = []fyne.CanvasObject{widget.NewLabel("Failed to load.")}
		} else {
			qualifyingResultsContainer.Objects = []fyne.CanvasObject{obj}
		}
		qualifyingResultsContainer.Refresh()
	}

	sprintResultsContainer := container.NewStack()
	sprintResultsContainer.Resize(fyne.NewSize(900, 780))
	updateSprintResults := func() {
		obj, err := uiupdaters.FetchSprintResultsUI(state)
		if err != nil {
			sprintResultsContainer.Objects = []fyne.CanvasObject{widget.NewLabel("Failed to load.")}
		} else {
			sprintResultsContainer.Objects = []fyne.CanvasObject{obj}
		}
		sprintResultsContainer.Refresh()
	}

	driverStandingsContainer := container.NewStack()
	driverStandingsContainer.Resize(fyne.NewSize(900, 780))
	updateDriverStandings := func() {
		obj, err := uiupdaters.FetchDriverStandingsUI(state)
		if err != nil {
			driverStandingsContainer.Objects = []fyne.CanvasObject{widget.NewLabel("Failed to load.")}
		} else {
			driverStandingsContainer.Objects = []fyne.CanvasObject{obj}
		}
		driverStandingsContainer.Refresh()
	}

	constructorStandingsContainer := container.NewStack()
	constructorStandingsContainer.Resize(fyne.NewSize(900, 780))
	updateConstructorStandings := func() {
		obj, err := uiupdaters.FetchConstructorStandingsUI(state)
		if err != nil {
			constructorStandingsContainer.Objects = []fyne.CanvasObject{widget.NewLabel("Failed to load.")}
		} else {
			constructorStandingsContainer.Objects = []fyne.CanvasObject{obj}
		}
		constructorStandingsContainer.Refresh()
	}

	countdown := binding.NewString()
	countdownLabel := widget.NewLabelWithData(countdown)
	seasonLabel := widget.NewLabel("Season: ")

	currentYearInt := time.Now().Year()
	years := make([]string, 0)
	for y := currentYearInt; y >= 1950; y-- {
		years = append(years, fmt.Sprintf("%d", y))
	}

	onYearSelected := func(newYear string) {
		state.SelectedYear = newYear
		results := make(chan tabLoadResult)

		updateSchedule() // Doesn't like being in a channel
		go func() { updateUpcoming(); results <- tabLoadResult{"Upcoming", nil} }()
		go func() { updateRaceResults(); results <- tabLoadResult{"Race Results", nil} }()
		go func() { updateQualifyingResults(); results <- tabLoadResult{"Qualifying Results", nil} }()
		go func() { updateSprintResults(); results <- tabLoadResult{"Sprint Results", nil} }()
		go func() { updateDriverStandings(); results <- tabLoadResult{"Driver Standings", nil} }()
		go func() { updateConstructorStandings(); results <- tabLoadResult{"Constructor Standings", nil} }()

		go func() {
			for i := 0; i < 6; i++ {
				res := <-results
				if state.Debug {
					fmt.Println("✅ Loaded", res.name)
				}
			}
		}()
	}

	seasonSelect := widget.NewSelect(years, func(selected string) {
		if state.Debug {
			fmt.Println("Selected season:", selected)
		}
		onYearSelected(selected)
	})
	if len(years) > 0 {
		seasonSelect.SetSelected(years[0])
		state.SelectedYear = years[0]
	}

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

	state.ResultsTabs = resultsTabs
	state.OuterTabs = outerTabs

	content := container.NewBorder(topRow, nil, nil, nil, outerTabs)
	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(900, 800))

	layout.StartCountdown(countdown, nextRace)
	myWindow.ShowAndRun()
}
