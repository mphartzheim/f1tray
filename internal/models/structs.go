package models

import (
	"time"

	"fyne.io/fyne/v2"
)

// QualifyingResultResponse represents the parsed JSON response for qualifying results.
type QualifyingResultResponse struct {
	MRData struct {
		RaceTable struct {
			Season string `json:"season"`
			Races  []struct {
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

// SprintResultResponse represents the parsed JSON response for sprint results.
type SprintResultResponse struct {
	MRData struct {
		RaceTable struct {
			Season string `json:"season"`
			Races  []struct {
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

// RaceResultResponse represents the parsed JSON response for race results.
type RaceResultResponse struct {
	MRData struct {
		RaceTable struct {
			Season string `json:"season"`
			Races  []struct {
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

// ScheduleResponse represents the parsed JSON response for the full race schedule.
type ScheduleResponse struct {
	MRData struct {
		RaceTable struct {
			Season string `json:"season"`
			Races  []struct {
				Round    string `json:"round"`
				RaceName string `json:"raceName"`
				Date     string `json:"date"`
				Circuit  struct {
					CircuitName string `json:"circuitName"`
					Location    struct {
						Locality string `json:"locality"`
						Country  string `json:"country"`
						Lat      string `json:"lat"`
						Long     string `json:"long"`
					} `json:"Location"`
				} `json:"Circuit"`
			} `json:"Races"`
		} `json:"RaceTable"`
	} `json:"MRData"`
}

// UpcomingResponse represents the parsed JSON response for the next upcoming race and sessions.
type UpcomingResponse struct {
	MRData struct {
		RaceTable struct {
			Season string `json:"season"`
			Races  []struct {
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

// TabData contains the content and refresh logic for a UI tab.
type TabData struct {
	Content fyne.CanvasObject
	Refresh func() bool
}

// AppState tracks application state like whether it's the first run.
type AppState struct {
	FirstRun         bool
	UpcomingSessions []SessionInfo
}

// SessionInfo holds individual session data for notification purposes.
type SessionInfo struct {
	Type      string    // "Practice", "Qualifying", "Race"
	StartTime time.Time // UTC timestamp of session start
	Label     string    // Human-readable name (e.g., "Free Practice 1")
}
