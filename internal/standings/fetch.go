package standings

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mphartzheim/f1tray/internal/appstate"
)

func FetchDriverStandings(state *appstate.AppState, url string) ([]DriverStandingItem, error) {
	if state.Debug {
		fmt.Println("🛠️ Fetching Driver Standings from:", url)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch driver standings: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read standings response: %w", err)
	}

	var parsed DriverStandingsResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal driver standings: %w", err)
	}

	if len(parsed.MRData.StandingsTable.StandingsLists) == 0 {
		return nil, fmt.Errorf("no standings data found")
	}

	return parsed.MRData.StandingsTable.StandingsLists[0].DriverStandings, nil
}

func FetchConstructorStandings(state *appstate.AppState, url string) ([]ConstructorStandingPosition, error) {
	if state.Debug {
		fmt.Println("🛠️ Fetching Constructor Standings from:", url)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch constructor standings: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read constructor standings response: %w", err)
	}

	var parsed ConstructorStandingsResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal constructor standings: %w", err)
	}

	if len(parsed.MRData.StandingsTable.StandingsLists) == 0 {
		return nil, fmt.Errorf("no constructor standings data found")
	}

	return parsed.MRData.StandingsTable.StandingsLists[0].ConstructorStandings, nil
}
