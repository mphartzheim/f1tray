package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"f1tray/internal/config"
)

var baseURL = config.BaseAPIURL + "/current"

// GetSessionResults fetches JSON data for a specific session
func GetSessionResults(session string) (*http.Response, error) {
	useNextPrefix := session == "qualifying" || session == "results"
	path := baseURL
	if useNextPrefix {
		path += "/next/"
	}
	url := fmt.Sprintf("%s%s.json", path, session)
	return http.Get(url)
}

// GetRaceSchedule fetches the full race schedule
func GetRaceSchedule() (*http.Response, error) {
	url := baseURL + ".json"
	return http.Get(url)
}

// GetNextRace fetches the next race's metadata
func GetNextRace() (string, error) {
	resp, err := http.Get(baseURL + "/next.json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data struct {
		MRData struct {
			RaceTable struct {
				Races []struct {
					Round string `json:"round"`
				} `json:"Races"`
			} `json:"RaceTable"`
		} `json:"MRData"`
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	if len(data.MRData.RaceTable.Races) == 0 {
		return "", fmt.Errorf("no upcoming races found")
	}

	return data.MRData.RaceTable.Races[0].Round, nil
}
