package upcoming

type NextRaceResponse struct {
	MRData NextMRData `json:"MRData"`
}

type NextMRData struct {
	RaceTable NextRaceTable `json:"RaceTable"`
}

type NextRaceTable struct {
	Season string     `json:"season"`
	Races  []NextRace `json:"Races"`
}

type NextRace struct {
	Season         string       `json:"season"`
	Round          string       `json:"round"`
	URL            string       `json:"url"`
	RaceName       string       `json:"raceName"`
	Circuit        NextCircuit  `json:"Circuit"`
	Date           string       `json:"date"`
	Time           string       `json:"time"`
	FirstPractice  *NextSession `json:"FirstPractice,omitempty"`
	SecondPractice *NextSession `json:"SecondPractice,omitempty"`
	ThirdPractice  *NextSession `json:"ThirdPractice,omitempty"`
	Sprint         *NextSession `json:"Sprint,omitempty"`
	Qualifying     *NextSession `json:"Qualifying,omitempty"`
}

type NextCircuit struct {
	CircuitID   string       `json:"circuitId"`
	URL         string       `json:"url"`
	CircuitName string       `json:"circuitName"`
	Location    NextLocation `json:"Location"`
}

type NextLocation struct {
	Lat      string `json:"lat"`
	Long     string `json:"long"`
	Locality string `json:"locality"`
	Country  string `json:"country"`
}

type NextSession struct {
	Date string `json:"date"`
	Time string `json:"time"`
}
