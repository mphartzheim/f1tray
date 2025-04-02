package schedule

type ScheduledRace struct {
	Season   string           `json:"season"`
	Round    string           `json:"round"`
	URL      string           `json:"url"`
	RaceName string           `json:"raceName"`
	Circuit  ScheduledCircuit `json:"Circuit"`
	Date     string           `json:"date"`
	Time     string           `json:"time"`
}

type ScheduledCircuit struct {
	CircuitID   string            `json:"circuitId"`
	URL         string            `json:"url"`
	CircuitName string            `json:"circuitName"`
	Location    ScheduledLocation `json:"Location"`
}

type ScheduledLocation struct {
	Lat      string `json:"lat"`
	Long     string `json:"long"`
	Locality string `json:"locality"`
	Country  string `json:"country"`
}

type ScheduleAPIResponse struct {
	MRData ScheduleMRData `json:"MRData"`
}

type ScheduleMRData struct {
	RaceTable ScheduleRaceTable `json:"RaceTable"`
}

type ScheduleRaceTable struct {
	Season string          `json:"season"`
	Races  []ScheduledRace `json:"Races"`
}
