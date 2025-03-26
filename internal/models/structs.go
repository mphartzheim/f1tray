package models

import (
	"f1tray/internal/config"

	"fyne.io/fyne/v2"
)

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

type UpcomingResponse struct {
	MRData struct {
		RaceTable struct {
			Races []struct {
				RaceName string `json:"raceName"`
				Date     string `json:"date"`
				Time     string `json:"time"`
				Circuit  struct {
					CircuitName string `json:"circuitName"`
					Location    struct {
						Locality string `json:"locality"`
						Country  string `json:"country"`
					} `json:"Location"`
				} `json:"Circuit"`
				FirstPractice struct {
					Date string `json:"date"`
					Time string `json:"time"`
				} `json:"FirstPractice"`
				SecondPractice struct {
					Date string `json:"date"`
					Time string `json:"time"`
				} `json:"SecondPractice"`
				ThirdPractice struct {
					Date string `json:"date"`
					Time string `json:"time"`
				} `json:"ThirdPractice"`
				Qualifying struct {
					Date string `json:"date"`
					Time string `json:"time"`
				} `json:"Qualifying"`
				Sprint struct {
					Date string `json:"date"`
					Time string `json:"time"`
				} `json:"Sprint"`
			} `json:"Races"`
		} `json:"RaceTable"`
	} `json:"MRData"`
}

type TabData struct {
	Content fyne.CanvasObject
	Refresh func() bool
}

type AppState struct {
	DebugMode   bool
	Preferences config.Preferences
	// Add more things later like logger, user session, etc
}
