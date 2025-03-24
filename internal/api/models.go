package api

type MRDataContainer[T any] struct {
	MRData struct {
		RaceTable struct {
			Races []T `json:"Races"`
		} `json:"RaceTable"`
	} `json:"MRData"`
}

type QualifyingRace struct {
	RaceName          string             `json:"raceName"`
	QualifyingResults []QualifyingResult `json:"QualifyingResults"`
}

type SprintRace struct {
	RaceName      string         `json:"raceName"`
	SprintResults []SprintResult `json:"SprintResults"`
}

type ResultsRace struct {
	RaceName string        `json:"raceName"`
	Results  []FinalResult `json:"Results"`
}

type QualifyingResult struct {
	Position    string      `json:"position"`
	Driver      Driver      `json:"Driver"`
	Constructor Constructor `json:"Constructor"`
	Q1          string      `json:"Q1"`
	Q2          string      `json:"Q2"`
	Q3          string      `json:"Q3"`
}

type SprintResult struct {
	Position    string      `json:"position"`
	Driver      Driver      `json:"Driver"`
	Constructor Constructor `json:"Constructor"`
	Time        struct {
		Time string `json:"time"`
	} `json:"Time"`
	Status string `json:"status"`
}

type FinalResult struct {
	Position    string      `json:"position"`
	Driver      Driver      `json:"Driver"`
	Constructor Constructor `json:"Constructor"`
	Time        struct {
		Time string `json:"time"`
	} `json:"Time"`
	Status string `json:"status"`
}

type Driver struct {
	GivenName  string `json:"givenName"`
	FamilyName string `json:"familyName"`
}

type Constructor struct {
	Name string `json:"name"`
}

type RaceSchedule struct {
	RaceName string `json:"raceName"`
	Date     string `json:"date"`
	Round    string `json:"round"`
}
