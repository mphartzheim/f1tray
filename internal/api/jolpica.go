package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Location represents the location details of the circuit.
type Location struct {
	Locality string `json:"locality"`
	Country  string `json:"country"`
}

// Circuit represents the circuit details.
type Circuit struct {
	CircuitName string   `json:"circuitName"`
	Location    Location `json:"Location"`
}

// Race represents the race details.
type Race struct {
	Season   string  `json:"season"`
	Round    string  `json:"round"`
	RaceName string  `json:"raceName"`
	Circuit  Circuit `json:"Circuit"`
	Date     string  `json:"date"`
	Time     string  `json:"time"`
}

// RaceTable holds a list of races.
type RaceTable struct {
	Races []Race `json:"Races"`
}

// MRData is the top-level structure of the JSON response.
type MRData struct {
	RaceTable RaceTable `json:"RaceTable"`
}

// APIResponse wraps MRData
type APIResponse struct {
	MRData MRData `json:"MRData"`
}

// FetchNextRace retrieves details of the next scheduled F1 race.
func FetchNextRace() (*Race, error) {
	url := "https://api.jolpi.ca/ergast/f1/current/next.json"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResponse APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(apiResponse.MRData.RaceTable.Races) == 0 {
		return nil, fmt.Errorf("no upcoming races found")
	}

	return &apiResponse.MRData.RaceTable.Races[0], nil
}
