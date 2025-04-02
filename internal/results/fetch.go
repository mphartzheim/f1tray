package results

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mphartzheim/f1tray/internal/appstate"
)

func FetchRaceResults(state *appstate.AppState, url string) (*RaceResultsEvent, error) {
	if state.Debug {
		fmt.Println("🛠️ Fetching Race Results from:", url)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch results: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read results response: %w", err)
	}

	var parsed RaceResultsResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal results: %w", err)
	}

	if len(parsed.MRData.RaceTable.Races) == 0 {
		return nil, fmt.Errorf("no races in response")
	}

	return &parsed.MRData.RaceTable.Races[0], nil
}

func FetchQualifyingResults(state *appstate.AppState, url string) (*QualifyingEvent, error) {
	if state.Debug {
		fmt.Println("🛠️ Fetching Qualifying Results from:", url)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch qualifying: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read qualifying response: %w", err)
	}

	var parsed QualifyingResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal qualifying: %w", err)
	}

	if len(parsed.MRData.RaceTable.Races) == 0 {
		return nil, fmt.Errorf("no qualifying races found")
	}

	return &parsed.MRData.RaceTable.Races[0], nil
}

func FetchSprintResults(state *appstate.AppState, url string) (*SprintEvent, error) {
	if state.Debug {
		fmt.Println("🛠️ Fetching Sprint Results from:", url)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sprint results: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read sprint results: %w", err)
	}

	var parsed SprintResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal sprint results: %w", err)
	}

	if len(parsed.MRData.RaceTable.Races) == 0 || len(parsed.MRData.RaceTable.Races[0].SprintResults) == 0 {
		return nil, nil // No sprint data available
	}

	return &parsed.MRData.RaceTable.Races[0], nil
}
