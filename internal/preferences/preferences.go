package preferences

import (
	"encoding/json"
	"os"
	"time"
)

type Preferences struct {
	RaceReminderHours  int       `json:"raceReminderHours"`
	WeeklyReminderDay  string    `json:"weeklyReminderDay"`
	WeeklyReminderTime time.Time `json:"weeklyReminderTime"`
}

func LoadPrefs() (*Preferences, error) {
	// Load preferences from a file (example: preferences.json)
	file, err := os.Open("preferences.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var prefs Preferences
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&prefs)
	if err != nil {
		return nil, err
	}

	return &prefs, nil
}

func SavePrefs(prefs *Preferences) error {
	// Save preferences to a file (example: preferences.json)
	file, err := os.Create("preferences.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(prefs)
	if err != nil {
		return err
	}

	return nil
}
