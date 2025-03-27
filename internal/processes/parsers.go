package processes

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

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

	// Use the first race as the "next" race.
	race := races[0]
	location := fmt.Sprintf("%s, %s", race.Circuit.Location.Locality, race.Circuit.Location.Country)
	title := fmt.Sprintf("Next Race: %s (%s)", race.RaceName, location)

	// Helper function to convert UTC date/time to local.
	localize := func(dateStr, timeStr string) (string, string) {
		// Combine date and time into a RFC3339 string.
		// Example: dateStr "2023-05-01" and timeStr "14:00:00" become "2023-05-01T14:00:00"
		combined := fmt.Sprintf("%sT%s", dateStr, timeStr)
		// Append "Z" only if timeStr doesn't already include a timezone indicator.
		if len(timeStr) > 0 && timeStr[len(timeStr)-1] != 'Z' && !strings.Contains(timeStr, "+") && !strings.Contains(timeStr, "-") {
			combined += "Z"
		}

		t, err := time.Parse(time.RFC3339, combined)
		if err != nil {
			fmt.Println("Error parsing datetime:", err)
			// If parsing fails, return the original values.
			return dateStr, timeStr
		}
		local := t.Local()
		return local.Format("2006-01-02"), local.Format("15:04 MST")
	}

	var rows [][]string

	// Append session rows first, converting each session's time.
	if race.FirstPractice.Date != "" && race.FirstPractice.Time != "" {
		d, t := localize(race.FirstPractice.Date, race.FirstPractice.Time)
		rows = append(rows, []string{"Practice 1", d, t})
	}
	if race.SecondPractice.Date != "" && race.SecondPractice.Time != "" {
		d, t := localize(race.SecondPractice.Date, race.SecondPractice.Time)
		rows = append(rows, []string{"Practice 2", d, t})
	}
	if race.ThirdPractice.Date != "" && race.ThirdPractice.Time != "" {
		d, t := localize(race.ThirdPractice.Date, race.ThirdPractice.Time)
		rows = append(rows, []string{"Practice 3", d, t})
	}
	if race.Qualifying.Date != "" && race.Qualifying.Time != "" {
		d, t := localize(race.Qualifying.Date, race.Qualifying.Time)
		rows = append(rows, []string{"Qualifying", d, t})
	}
	if race.Sprint.Date != "" && race.Sprint.Time != "" {
		d, t := localize(race.Sprint.Date, race.Sprint.Time)
		rows = append(rows, []string{"Sprint", d, t})
	}

	// Append the Race row last.
	if race.Date != "" && race.Time != "" {
		d, t := localize(race.Date, race.Time)
		rows = append(rows, []string{"Race", d, t})
	}

	return title, rows, nil
}
