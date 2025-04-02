package results

type RaceResultsResponse struct {
	MRData RaceResultsMRData `json:"MRData"`
}

type RaceResultsMRData struct {
	RaceTable RaceResultsTable `json:"RaceTable"`
}

type RaceResultsTable struct {
	Races []RaceResultsEvent `json:"Races"`
}

type RaceResultsEvent struct {
	RaceName string                   `json:"raceName"`
	Date     string                   `json:"date"`
	Circuit  RaceResultsCircuit       `json:"Circuit"`
	Results  []RaceResultsPositioning `json:"Results"`
}

type RaceResultsCircuit struct {
	CircuitName string                   `json:"circuitName"`
	Location    RaceResultsCircuitLocale `json:"Location"`
}

type RaceResultsCircuitLocale struct {
	Locality string `json:"locality"`
	Country  string `json:"country"`
}

type RaceResultsPositioning struct {
	Position    string                 `json:"position"`
	Driver      RaceResultsDriver      `json:"Driver"`
	Constructor RaceResultsConstructor `json:"Constructor"`
	Grid        string                 `json:"grid"`
	Laps        string                 `json:"laps"`
	Status      string                 `json:"status"`
	Time        *RaceResultsFinishTime `json:"Time"` // May be nil for DNFs
}

type RaceResultsDriver struct {
	DriverID    string `json:"driverId"`
	GivenName   string `json:"givenName"`
	FamilyName  string `json:"familyName"`
	Nationality string `json:"nationality"`
}

type RaceResultsConstructor struct {
	Name string `json:"name"`
}

type RaceResultsFinishTime struct {
	Time string `json:"time"`
}

type QualifyingResponse struct {
	MRData QualifyingMRData `json:"MRData"`
}

type QualifyingMRData struct {
	RaceTable QualifyingRaceTable `json:"RaceTable"`
}

type QualifyingRaceTable struct {
	Races []QualifyingEvent `json:"Races"`
}

type QualifyingEvent struct {
	RaceName string             `json:"raceName"`
	Date     string             `json:"date"`
	Circuit  RaceResultsCircuit `json:"Circuit"` // reuse
	Results  []QualifyingEntry  `json:"QualifyingResults"`
}

type QualifyingEntry struct {
	Position    string                 `json:"position"`
	Driver      RaceResultsDriver      `json:"Driver"`      // reuse
	Constructor RaceResultsConstructor `json:"Constructor"` // reuse
	Q1          string                 `json:"Q1"`
	Q2          string                 `json:"Q2"`
	Q3          string                 `json:"Q3"`
}

type SprintResponse struct {
	MRData SprintMRData `json:"MRData"`
}

type SprintMRData struct {
	RaceTable SprintRaceTable `json:"RaceTable"`
}

type SprintRaceTable struct {
	Races []SprintEvent `json:"Races"`
}

type SprintEvent struct {
	RaceName      string              `json:"raceName"`
	Date          string              `json:"date"`
	Circuit       RaceResultsCircuit  `json:"Circuit"` // Reuse from race models
	SprintResults []SprintResultEntry `json:"SprintResults"`
}

type SprintResultEntry struct {
	Position    string                 `json:"position"`
	Driver      RaceResultsDriver      `json:"Driver"`      // Reuse
	Constructor RaceResultsConstructor `json:"Constructor"` // Reuse
	Grid        string                 `json:"grid"`
	Laps        string                 `json:"laps"`
	Status      string                 `json:"status"`
	Time        *RaceResultsFinishTime `json:"Time"` // Optional
}
