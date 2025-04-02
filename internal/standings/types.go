package standings

type DriverStandingsResponse struct {
	MRData DriverStandingsMRData `json:"MRData"`
}

type DriverStandingsMRData struct {
	StandingsTable DriverStandingsTable `json:"StandingsTable"`
}

type DriverStandingsTable struct {
	StandingsLists []DriverStandingsList `json:"StandingsLists"`
}

type DriverStandingsList struct {
	Season          string               `json:"season"`
	Round           string               `json:"round"`
	DriverStandings []DriverStandingItem `json:"DriverStandings"`
}

type DriverStandingItem struct {
	Position     string                      `json:"position"`
	Points       string                      `json:"points"`
	Wins         string                      `json:"wins"`
	Driver       DriverStandingDriver        `json:"Driver"`
	Constructors []DriverStandingConstructor `json:"Constructors"`
}

type DriverStandingDriver struct {
	DriverID    string `json:"driverId"`
	GivenName   string `json:"givenName"`
	FamilyName  string `json:"familyName"`
	DateOfBirth string `json:"dateOfBirth"`
	Nationality string `json:"nationality"`
}

type DriverStandingConstructor struct {
	Name string `json:"name"`
}

type ConstructorStandingsResponse struct {
	MRData ConstructorStandingsMRData `json:"MRData"`
}

type ConstructorStandingsMRData struct {
	StandingsTable ConstructorStandingsTable `json:"StandingsTable"`
}

type ConstructorStandingsTable struct {
	StandingsLists []ConstructorStandingsList `json:"StandingsLists"`
}

type ConstructorStandingsList struct {
	Season               string                        `json:"season"`
	Round                string                        `json:"round"`
	ConstructorStandings []ConstructorStandingPosition `json:"ConstructorStandings"`
}

type ConstructorStandingPosition struct {
	Position    string                         `json:"position"`
	Points      string                         `json:"points"`
	Wins        string                         `json:"wins"`
	Constructor ConstructorStandingConstructor `json:"Constructor"`
}

type ConstructorStandingConstructor struct {
	ConstructorID string `json:"constructorId"`
	Name          string `json:"name"`
	Nationality   string `json:"nationality"`
}
