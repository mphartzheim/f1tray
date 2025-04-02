package upcoming

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchNextRace() (*NextRace, error) {
	resp, err := http.Get("https://api.jolpi.ca/ergast/f1/current/next.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res NextRaceResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	if len(res.MRData.RaceTable.Races) == 0 {
		return nil, fmt.Errorf("no races found")
	}

	return &res.MRData.RaceTable.Races[0], nil
}
