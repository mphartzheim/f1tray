package models

const (
	// ScheduleURL is the API endpoint for the full F1 schedule for a given year.
	ScheduleURL = "https://api.jolpi.ca/ergast/f1/%s.json"

	// UpcomingURL is the API endpoint for the next upcoming F1 race.
	UpcomingURL = "https://api.jolpi.ca/ergast/f1/%s/next.json"

	// DriversStandingsURL is the API endpoint for driver standings for a given year.
	DriversStandingsURL = "https://api.jolpi.ca/ergast/f1/%s/driverstandings.json"

	// ConstrcutorsStandingsURL is the API endpoint for constructor standings for a given year.
	ConstructorsStandingsURL = "https://api.jolpi.ca/ergast/f1/%s/constructorstandings.json"

	// RaceURL is the API endpoint for race results by year and round.
	RaceURL = "https://api.jolpi.ca/ergast/f1/%s/%s/results.json"

	// QualifyingURL is the API endpoint for qualifying results by year and round.
	QualifyingURL = "https://api.jolpi.ca/ergast/f1/%s/%s/qualifying.json"

	// SprintURL is the API endpoint for sprint results by year and round.
	SprintURL = "https://api.jolpi.ca/ergast/f1/%s/%s/sprint.json"

	// F1tvURL is the direct link to the F1TV streaming platform.
	F1tvURL = "https://f1tv.formula1.com/"

	// MapBaseURL is the base OpenStreetMap URL used for race location links.
	MapBaseURL = "https://www.openstreetmap.org/"
)
