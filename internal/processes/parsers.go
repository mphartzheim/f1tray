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
		timeOrStatus := getTimeOrStatus(res.Status, res.Time.Time)
		driverName := buildDriverDisplayName(res.Driver.GivenName, res.Driver.FamilyName, res.Driver.URL)

		url := models.ConstructorURLMap[res.Constructor.Name]
		constructorName := buildConstructorDisplayName(res.Constructor.Name, url)

		rows[i] = []string{
			res.Position,
			"",
			driverName,
			constructorName,
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
		timeOrStatus := getTimeOrStatus(res.Status, res.Time.Time)
		driverName := buildDriverDisplayName(res.Driver.GivenName, res.Driver.FamilyName, res.Driver.URL)

		url := models.ConstructorURLMap[res.Constructor.Name]
		constructorName := buildConstructorDisplayName(res.Constructor.Name, url)

		rows[i] = []string{
			res.Position,
			"",
			driverName,
			constructorName,
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
		bestTime := bestQualifyingTime(res.Q1, res.Q2, res.Q3)
		driverName := buildDriverDisplayName(res.Driver.GivenName, res.Driver.FamilyName, res.Driver.URL)

		url := models.ConstructorURLMap[res.Constructor.Name]
		constructorName := buildConstructorDisplayName(res.Constructor.Name, url)

		rows[i] = []string{
			res.Position,
			"",
			driverName,
			constructorName,
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
		driverName := buildDriverDisplayName(s.Driver.GivenName, s.Driver.FamilyName, s.Driver.URL)

		constructor := s.Constructors[0]
		url := models.ConstructorURLMap[constructor.Name]
		constructorName := buildConstructorDisplayName(constructor.Name, url)

		rows[i] = []string{
			s.Position,
			"",
			driverName,
			constructorName,
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
		url := models.ConstructorURLMap[s.Constructor.Name]
		constructorName := buildConstructorDisplayName(s.Constructor.Name, url)

		rows[i] = []string{
			s.Position,
			"",
			constructorName,
			s.Points,
		}
	}

	title := fmt.Sprintf("Constructor Standings (%s)", standingsLists[0].Season)
	return title, rows, nil
}

// buildDriverDisplayName creates a driver display string with an embedded fallback URL and clickable indicator.
func buildDriverDisplayName(givenName, familyName, url string) string {
	fullName := fmt.Sprintf("%s %s", givenName, familyName)
	if url != "" {
		return fmt.Sprintf("%s|||%s%s", fullName, url, " 👤")
	}
	return fullName
}

// buildConstructorDisplayName creates a constructor display string with an embedded fallback URL and clickable indicator.
func buildConstructorDisplayName(name, _ string) string {
	// Step 1: Try the ConstructorURLMap first.
	if url, ok := models.ConstructorURLMap[name]; ok && url != "" {
		return fmt.Sprintf("%s|||%s 🌐", name, url)
	}

	// Step 2: Try fallback from AllConstructors (downloaded JSON).
	for _, c := range models.AllConstructors {
		if c.Name == name && c.URL != "" {
			return fmt.Sprintf("%s|||%s 🌐", name, c.URL)
		}
	}

	// Step 3: No URL available.
	return name
}

// getTimeOrStatus returns the time if available, otherwise the status.
func getTimeOrStatus(status, timeField string) string {
	if timeField != "" {
		return timeField
	}
	return status
}

// bestQualifyingTime returns the best qualifying time available, checking Q3 first, then Q2, and finally Q1.
func bestQualifyingTime(q1, q2, q3 string) string {
	if q3 != "" {
		return q3
	} else if q2 != "" {
		return q2
	}
	return q1
}
