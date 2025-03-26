package processes

import (
	"encoding/json"
	"fmt"

	"f1tray/internal/models"
)

func ParseRaceResults(body []byte) (string, [][]string, error) {
	var result models.RaceResultResponse
	err := json.Unmarshal(body, &result)
	if err != nil {
		return "", nil, fmt.Errorf("JSON error: %v", err)
	}

	if len(result.MRData.RaceTable.Races) == 0 {
		return "", nil, fmt.Errorf("no race data found")
	}

	race := result.MRData.RaceTable.Races[0]
	rows := make([][]string, len(race.Results))
	for i, res := range race.Results {
		timeOrStatus := res.Status
		if res.Time.Time != "" {
			timeOrStatus = res.Time.Time
		}
		rows[i] = []string{
			res.Position,
			fmt.Sprintf("%s %s", res.Driver.GivenName, res.Driver.FamilyName),
			res.Constructor.Name,
			timeOrStatus,
		}
	}
	return race.RaceName, rows, nil
}

func ParseSprintResults(body []byte) (string, [][]string, error) {
	var result models.SprintResultResponse
	err := json.Unmarshal(body, &result)
	if err != nil {
		return "", nil, fmt.Errorf("JSON error: %v", err)
	}

	if len(result.MRData.RaceTable.Races) == 0 {
		return "", nil, fmt.Errorf("no sprint data found")
	}

	race := result.MRData.RaceTable.Races[0]
	rows := make([][]string, len(race.SprintResults))
	for i, res := range race.SprintResults {
		timeOrStatus := res.Status
		if res.Time.Time != "" {
			timeOrStatus = res.Time.Time
		}
		rows[i] = []string{
			res.Position,
			fmt.Sprintf("%s %s", res.Driver.GivenName, res.Driver.FamilyName),
			res.Constructor.Name,
			timeOrStatus,
		}
	}
	return race.RaceName, rows, nil
}

func ParseQualifyingResults(body []byte) (string, [][]string, error) {
	var result models.QualifyingResultResponse
	err := json.Unmarshal(body, &result)
	if err != nil {
		return "", nil, fmt.Errorf("JSON error: %v", err)
	}

	if len(result.MRData.RaceTable.Races) == 0 {
		return "", nil, fmt.Errorf("no qualifying data found")
	}

	race := result.MRData.RaceTable.Races[0]
	rows := make([][]string, len(race.QualifyingResults))
	for i, res := range race.QualifyingResults {
		bestTime := res.Q3
		if bestTime == "" {
			bestTime = res.Q2
		}
		if bestTime == "" {
			bestTime = res.Q1
		}
		rows[i] = []string{
			res.Position,
			fmt.Sprintf("%s %s", res.Driver.GivenName, res.Driver.FamilyName),
			res.Constructor.Name,
			bestTime,
		}
	}
	return race.RaceName, rows, nil
}

func ParseSchedule(body []byte) (string, [][]string, error) {
	var result models.ScheduleResponse
	err := json.Unmarshal(body, &result)
	if err != nil {
		return "", nil, fmt.Errorf("JSON error: %v", err)
	}

	races := result.MRData.RaceTable.Races
	rows := make([][]string, len(races))
	for i, race := range races {
		rows[i] = []string{
			race.Round,
			race.RaceName,
			race.Circuit.CircuitName,
			fmt.Sprintf("%s, %s (%s)", race.Circuit.Location.Locality, race.Circuit.Location.Country, race.Date),
		}
	}
	return "Current Season Schedule", rows, nil
}

func ParseUpcoming(body []byte) (string, [][]string, error) {
	var result models.UpcomingResponse
	err := json.Unmarshal(body, &result)
	if err != nil {
		return "", nil, fmt.Errorf("JSON error: %v", err)
	}

	races := result.MRData.RaceTable.Races
	if len(races) == 0 {
		return "", nil, fmt.Errorf("no upcoming race data found")
	}

	race := races[0]
	location := fmt.Sprintf("%s, %s", race.Circuit.Location.Locality, race.Circuit.Location.Country)
	title := fmt.Sprintf("Next Race: %s (%s)", race.RaceName, location)

	var rows [][]string

	// Append session rows first
	if race.FirstPractice.Date != "" {
		rows = append(rows, []string{"Practice 1", race.FirstPractice.Date, race.FirstPractice.Time})
	}
	if race.SecondPractice.Date != "" {
		rows = append(rows, []string{"Practice 2", race.SecondPractice.Date, race.SecondPractice.Time})
	}
	if race.ThirdPractice.Date != "" {
		rows = append(rows, []string{"Practice 3", race.ThirdPractice.Date, race.ThirdPractice.Time})
	}
	if race.Qualifying.Date != "" {
		rows = append(rows, []string{"Qualifying", race.Qualifying.Date, race.Qualifying.Time})
	}
	if race.Sprint.Date != "" {
		rows = append(rows, []string{"Sprint", race.Sprint.Date, race.Sprint.Time})
	}

	// Append the Race row last
	rows = append(rows, []string{"Race", race.Date, race.Time})

	return title, rows, nil
}
