package uiupdaters

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/mphartzheim/f1tray/internal/appstate"
	"github.com/mphartzheim/f1tray/internal/models"
	"github.com/mphartzheim/f1tray/internal/results"
	"github.com/mphartzheim/f1tray/internal/schedule"
	"github.com/mphartzheim/f1tray/internal/standings"
	"github.com/mphartzheim/f1tray/internal/upcoming"
)

func FetchScheduleUI(state *appstate.AppState) (fyne.CanvasObject, error) {
	races, err := schedule.FetchSchedule(state)
	if err != nil {
		return nil, err
	}
	return container.NewStack(schedule.CreateScheduleTable(state, races)), nil
}

func FetchUpcomingUI(state *appstate.AppState) (fyne.CanvasObject, error) {
	nextRace, err := upcoming.FetchNextRace()
	if err != nil {
		return widget.NewLabel("Failed to load upcoming race."), nil
	}
	return container.NewStack(upcoming.CreateUpcomingTable(state, nextRace)), nil
}

func FetchRaceResultsUI(state *appstate.AppState) (fyne.CanvasObject, error) {
	raceURL := fmt.Sprintf(models.RaceURL, state.SelectedYear, "last")
	data, err := results.FetchRaceResults(state, raceURL)
	if err != nil {
		return nil, err
	}
	return results.CreateRaceResultsTable(state, data), nil
}

func FetchQualifyingResultsUI(state *appstate.AppState) (fyne.CanvasObject, error) {
	qualURL := fmt.Sprintf(models.QualifyingURL, state.SelectedYear, "last")
	data, err := results.FetchQualifyingResults(state, qualURL)
	if err != nil {
		return nil, err
	}
	return results.CreateQualifyingResultsTable(state, data), nil
}

func FetchSprintResultsUI(state *appstate.AppState) (fyne.CanvasObject, error) {
	sprintURL := fmt.Sprintf(models.SprintURL, state.SelectedYear, "last")
	data, err := results.FetchSprintResults(state, sprintURL)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return widget.NewLabel("Not a Sprint event."), nil
	}
	return results.CreateSprintResultsTable(state, data), nil
}

func FetchDriverStandingsUI(state *appstate.AppState) (fyne.CanvasObject, error) {
	driverURL := fmt.Sprintf(models.DriversStandingsURL, state.SelectedYear)
	data, err := standings.FetchDriverStandings(state, driverURL)
	if err != nil {
		return nil, err
	}
	return standings.CreateDriverStandingsTable(state, data), nil
}

func FetchConstructorStandingsUI(state *appstate.AppState) (fyne.CanvasObject, error) {
	constructorURL := fmt.Sprintf(models.ConstructorsStandingsURL, state.SelectedYear)
	data, err := standings.FetchConstructorStandings(state, constructorURL)
	if err != nil {
		return nil, err
	}
	return standings.CreateConstructorStandingsTable(state, data), nil
}
