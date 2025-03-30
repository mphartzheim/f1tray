package models

import (
	"encoding/xml"
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
						DriverID        string `json:"driverId"`
						PermanentNumber string `json:"permanentNumber"`
						Code            string `json:"code"`
						URL             string `json:"url"`
						GivenName       string `json:"givenName"`
						FamilyName      string `json:"familyName"`
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
						DriverID        string `json:"driverId"`
						PermanentNumber string `json:"permanentNumber"`
						Code            string `json:"code"`
						URL             string `json:"url"`
						GivenName       string `json:"givenName"`
						FamilyName      string `json:"familyName"`
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
						DriverID        string `json:"driverId"`
						PermanentNumber string `json:"permanentNumber"`
						Code            string `json:"code"`
						URL             string `json:"url"`
						GivenName       string `json:"givenName"`
						FamilyName      string `json:"familyName"`
						DateOfBirth     string `json:"dateOfBirth"`
						Nationality     string `json:"nationality"`
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

// DriverStandingsResponse represents the JSON structure for driver standings from the Ergast API.
type DriverStandingsResponse struct {
	MRData struct {
		StandingsTable struct {
			Season         string `json:"season"`
			StandingsLists []struct {
				Season          string `json:"season"`
				Round           string `json:"round"`
				DriverStandings []struct {
					Position     string `json:"position"`
					PositionText string `json:"positionText"`
					Points       string `json:"points"`
					Wins         string `json:"wins"`
					Driver       struct {
						DriverID        string `json:"driverId"`
						PermanentNumber string `json:"permanentNumber"`
						Code            string `json:"code"`
						URL             string `json:"url"`
						GivenName       string `json:"givenName"`
						FamilyName      string `json:"familyName"`
						DateOfBirth     string `json:"dateOfBirth"`
						Nationality     string `json:"nationality"`
					} `json:"Driver"`
					Constructors []struct {
						ConstructorID string `json:"constructorId"`
						URL           string `json:"url"`
						Name          string `json:"name"`
						Nationality   string `json:"nationality"`
					} `json:"Constructors"`
				} `json:"DriverStandings"`
			} `json:"StandingsLists"`
		} `json:"StandingsTable"`
	} `json:"MRData"`
}

// ConstructorStandingsResponse represents the JSON structure for constructor standings from the Ergast API.
type ConstructorStandingsResponse struct {
	MRData struct {
		StandingsTable struct {
			Season         string `json:"season"`
			StandingsLists []struct {
				Season               string `json:"season"`
				Round                string `json:"round"`
				ConstructorStandings []struct {
					Position     string `json:"position"`
					PositionText string `json:"positionText"`
					Points       string `json:"points"`
					Wins         string `json:"wins"`
					Constructor  struct {
						ConstructorID string `json:"constructorId"`
						URL           string `json:"url"`
						Name          string `json:"name"`
						Nationality   string `json:"nationality"`
					} `json:"Constructor"`
				} `json:"ConstructorStandings"`
			} `json:"StandingsLists"`
		} `json:"StandingsTable"`
	} `json:"MRData"`
}

// MRData represents the top-level structure of the API response.
type MRData struct {
	XMLName        xml.Name       `json:"MRData"`
	Series         string         `json:"series"`
	Limit          string         `json:"limit"`
	Offset         string         `json:"offset"`
	Total          string         `json:"total"`
	StandingsTable StandingsTable `json:"StandingsTable"`
}

// StandingsTable contains the season and a list of standings.
type StandingsTable struct {
	Season         string          `json:"season"`
	StandingsLists []StandingsList `json:"StandingsLists"`
}

// StandingsList represents the standings after a particular round.
type StandingsList struct {
	Season          string           `json:"season"`
	Round           string           `json:"round"`
	DriverStandings []DriverStanding `json:"DriverStandings"`
}

// DriverStanding holds the position, points, and other details of a driver.
type DriverStanding struct {
	Position     string        `json:"position"`
	PositionText string        `json:"positionText"`
	Points       string        `json:"points"`
	Wins         string        `json:"wins"`
	Driver       Driver        `json:"Driver"`
	Constructors []Constructor `json:"Constructors"`
}

// Driver contains personal and professional information about a driver.
type Driver struct {
	DriverID        string `json:"driverId"`
	PermanentNumber string `json:"permanentNumber"`
	Code            string `json:"code"`
	URL             string `json:"url"`
	GivenName       string `json:"givenName"`
	FamilyName      string `json:"familyName"`
	DateOfBirth     string `json:"dateOfBirth"`
	Nationality     string `json:"nationality"`
}

// Constructor represents a constructor (team) in the standings.
type Constructor struct {
	ConstructorID string `json:"constructorId"`
	URL           string `json:"url"`
	Name          string `json:"name"`
	Nationality   string `json:"nationality"`
}

// TabData contains the content and refresh logic for a UI tab.
type TabData struct {
	Content fyne.CanvasObject
	Refresh func() bool
}

// AppState tracks application state like whether it's the first run.
type AppState struct {
	FirstRun         bool          `json:"first_run"`
	UpcomingSessions []SessionInfo `json:"upcoming_sessions"`
	FavoriteDrivers  []string      `json:"favorite_drivers"`
}

// SessionInfo holds individual session data for notification purposes.
type SessionInfo struct {
	Type      string    // "Practice", "Qualifying", "Race"
	StartTime time.Time // UTC timestamp of session start
	Label     string    // Human-readable name (e.g., "Free Practice 1")
}
