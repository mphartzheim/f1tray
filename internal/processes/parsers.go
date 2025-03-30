package processes

import (
	"encoding/json"
	"fmt"

	"github.com/mphartzheim/f1tray/internal/config"
	"github.com/mphartzheim/f1tray/internal/models"
)

// ParseRaceResults extracts race result data into a table-friendly format from raw JSON.
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

// ParseSprintResults extracts sprint result data into a table-friendly format from raw JSON.
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

// ParseQualifyingResults extracts qualifying result data into a table-friendly format from raw JSON.
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

// ParseSchedule extracts the full race schedule into rows for display from raw JSON.
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
	return fmt.Sprintf("%s season schedule loaded", result.MRData.RaceTable.Season), rows, nil
}

// ParseUpcoming extracts session times for the next race into labeled rows using user time preferences.
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

	// Load user preferences to determine the time format.
	use24h := config.Get().Clock.Use24Hour

	race := races[0]
	location := fmt.Sprintf("%s, %s", race.Circuit.Location.Locality, race.Circuit.Location.Country)
	title := fmt.Sprintf("Next Race: %s (%s)", race.RaceName, location)

	var rows [][]string

	rows = AppendSessionRow(rows, "Practice 1", race.FirstPractice.Date, race.FirstPractice.Time, use24h)
	rows = AppendSessionRow(rows, "Practice 2", race.SecondPractice.Date, race.SecondPractice.Time, use24h)
	rows = AppendSessionRow(rows, "Practice 3", race.ThirdPractice.Date, race.ThirdPractice.Time, use24h)
	rows = AppendSessionRow(rows, "Qualifying", race.Qualifying.Date, race.Qualifying.Time, use24h)
	rows = AppendSessionRow(rows, "Sprint", race.Sprint.Date, race.Sprint.Time, use24h)
	rows = AppendSessionRow(rows, "Race", race.Date, race.Time, use24h)

	return title, rows, nil
}

// ParseDriverStandings extracts driver standings into a table-friendly format from raw JSON.
func ParseDriverStandings(body []byte) (string, [][]string, error) {
	var result models.DriverStandingsResponse
	err := json.Unmarshal(body, &result)
	if err != nil {
		return "", nil, fmt.Errorf("JSON error: %v", err)
	}

	standingsLists := result.MRData.StandingsTable.StandingsLists
	if len(standingsLists) == 0 {
		return "", nil, fmt.Errorf("no driver standings data found")
	}

	standings := standingsLists[0].DriverStandings
	rows := make([][]string, len(standings))
	for i, s := range standings {
		// Build the driver name.
		// If an API URL is provided, embed it using "|||" as a delimiter,
		// then append the clickable emoji.
		driverName := fmt.Sprintf("%s %s", s.Driver.GivenName, s.Driver.FamilyName)
		if s.Driver.URL != "" {
			driverName = fmt.Sprintf("%s|||%s%s", driverName, s.Driver.URL, " 👤")
		}
		rows[i] = []string{
			s.Position,
			driverName,
			s.Constructors[0].Name, // usually only one constructor per driver
			s.Points,
		}
	}

	title := fmt.Sprintf("Driver Standings (%s)", standingsLists[0].Season)
	return title, rows, nil
}

// ParseConstructorStandings extracts constructor standings into a table-friendly format from raw JSON.
func ParseConstructorStandings(body []byte) (string, [][]string, error) {
	var result models.ConstructorStandingsResponse
	err := json.Unmarshal(body, &result)
	if err != nil {
		return "", nil, fmt.Errorf("JSON error: %v", err)
	}

	standingsLists := result.MRData.StandingsTable.StandingsLists
	if len(standingsLists) == 0 {
		return "", nil, fmt.Errorf("no constructor standings data found")
	}

	standings := standingsLists[0].ConstructorStandings
	rows := make([][]string, len(standings))
	for i, s := range standings {
		rows[i] = []string{
			s.Position,
			s.Constructor.Name,
			s.Constructor.Nationality,
			s.Points,
		}
	}

	title := fmt.Sprintf("Constructor Standings (%s)", standingsLists[0].Season)
	return title, rows, nil
}
