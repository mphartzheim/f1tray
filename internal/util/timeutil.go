package util

import (
	"fmt"
	"time"

	"github.com/mphartzheim/f1tray/internal/appstate"
)

// Converts a date and time string to the user's local time.
func FormatToLocal(dateStr, timeStr string) string {
	raw := fmt.Sprintf("%sT%s", dateStr, timeStr)
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return fmt.Sprintf("%s %s", dateStr, timeStr)
	}
	return t.Local().Format("2006-01-02 15:04:05 MST")
}

// Parses a full formatted datetime string.
func ParseDateTime(raw string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05 MST", raw)
}

// Returns the current time (mocked if debug mode is on).
func GetNow(state *appstate.AppState) time.Time {
	if state.Debug {
		loc, _ := time.LoadLocation("America/Chicago")
		return time.Date(2025, 4, 3, 21, 30, 0, 0, loc) // Hard coded for testing
	}
	return time.Now().UTC()
}
