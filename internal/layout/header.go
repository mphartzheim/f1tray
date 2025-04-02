package layout

import (
	"fmt"
	"time"

	"github.com/mphartzheim/f1tray/internal/upcoming"
)

// GetNextSessionTime returns the next future session time and its label.
func GetNextSessionTime(race *upcoming.NextRace) (time.Time, string) {
	type sessionInfo struct {
		name string
		date string
		time string
	}

	// Helper to safely collect sessions
	collectSession := func(name string, s *upcoming.NextSession) sessionInfo {
		if s == nil {
			return sessionInfo{name, "", ""}
		}
		return sessionInfo{name, s.Date, s.Time}
	}

	sessions := []sessionInfo{
		collectSession("First Practice", race.FirstPractice),
		collectSession("Second Practice", race.SecondPractice),
		collectSession("Third Practice", race.ThirdPractice),
		collectSession("Sprint", race.Sprint),
		collectSession("Qualifying", race.Qualifying),
		{"Race", race.Date, race.Time},
	}

	now := time.Now().UTC()

	for _, s := range sessions {
		if s.date == "" {
			continue
		}
		raw := fmt.Sprintf("%sT%s", s.date, s.time)
		t, err := time.Parse(time.RFC3339, raw)
		if err == nil && t.After(now) {
			return t.Local(), s.name
		}
	}

	return time.Time{}, ""
}
