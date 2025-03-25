package models

type QualifyingResultResponse struct {
	MRData struct {
		RaceTable struct {
			Races []struct {
				RaceName          string `json:"raceName"`
				QualifyingResults []struct {
					Position string `json:"position"`
					Driver   struct {
						FamilyName string `json:"familyName"`
						GivenName  string `json:"givenName"`
					} `json:"Driver"`
					Constructor struct {
						Name string `json:"name"`
					} `json:"Constructor"`
					Q1 string `json:"Q1"`
					Q2 string `json:"Q2"`
					Q3 string `json:"Q3"`
				} `json:"QualifyingResults"`
			} `json:"Races"`
		} `json:"RaceTable"`
	} `json:"MRData"`
}

type SprintResultResponse struct {
	MRData struct {
		RaceTable struct {
			Races []struct {
				RaceName      string `json:"raceName"`
				SprintResults []struct {
					Position string `json:"position"`
					Driver   struct {
						FamilyName string `json:"familyName"`
						GivenName  string `json:"givenName"`
					} `json:"Driver"`
					Constructor struct {
						Name string `json:"name"`
					} `json:"Constructor"`
					Time struct {
						Time string `json:"time"`
					} `json:"Time"`
					Status string `json:"status"`
				} `json:"SprintResults"`
			} `json:"Races"`
		} `json:"RaceTable"`
	} `json:"MRData"`
}

type RaceResultResponse struct {
	MRData struct {
		RaceTable struct {
			Races []struct {
				RaceName string `json:"raceName"`
				Results  []struct {
					Position string `json:"position"`
					Driver   struct {
						FamilyName string `json:"familyName"`
						GivenName  string `json:"givenName"`
					} `json:"Driver"`
					Constructor struct {
						Name string `json:"name"`
					} `json:"Constructor"`
					Time struct {
						Time string `json:"time"`
					} `json:"Time"`
					Status string `json:"status"`
				} `json:"Results"`
			} `json:"Races"`
		} `json:"RaceTable"`
	} `json:"MRData"`
}

type ScheduleResponse struct {
	MRData struct {
		RaceTable struct {
			Races []struct {
				Round    string `json:"round"`
				RaceName string `json:"raceName"`
				Date     string `json:"date"`
				Circuit  struct {
					CircuitName string `json:"circuitName"`
					Location    struct {
						Locality string `json:"locality"`
						Country  string `json:"country"`
					} `json:"Location"`
				} `json:"Circuit"`
			} `json:"Races"`
		} `json:"RaceTable"`
	} `json:"MRData"`
}
