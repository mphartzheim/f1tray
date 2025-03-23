package preferences

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type UserPrefs struct {
	RaceReminderHours  int    `json:"raceReminderHours"`
	WeeklyReminderDay  string `json:"weeklyReminderDay"`
	WeeklyReminderHour int    `json:"weeklyReminderHour"`
}

func getConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appDir := filepath.Join(configDir, "f1tray")
	err = os.MkdirAll(appDir, 0755)
	if err != nil {
		return "", err
	}

	return filepath.Join(appDir, "config.json"), nil
}

func LoadPrefs() (UserPrefs, error) {
	path, err := getConfigPath()
	if err != nil {
		return UserPrefs{}, err
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		// Return default preferences if no file exists
		return UserPrefs{
			RaceReminderHours:  1,
			WeeklyReminderDay:  "Wednesday",
			WeeklyReminderHour: 12,
		}, nil
	} else if err != nil {
		return UserPrefs{}, err
	}

	var prefs UserPrefs
	if err := json.Unmarshal(data, &prefs); err != nil {
		return UserPrefs{}, fmt.Errorf("invalid config file: %w", err)
	}

	return prefs, nil
}

func SavePrefs(prefs UserPrefs) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(prefs, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
