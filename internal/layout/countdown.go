package layout

import (
	"fmt"
	"strings"
	"time"

	"github.com/mphartzheim/f1tray/internal/upcoming"

	"fyne.io/fyne/v2/data/binding"
)

func StartCountdown(countdown binding.String, nextRace *upcoming.NextRace) {
	go func() {
		for {
			sessionTime, label := GetNextSessionTime(nextRace)
			if sessionTime.IsZero() {
				countdown.Set("Countdown: no upcoming sessions")
				time.Sleep(10 * time.Second)
				continue
			}

			diff := time.Until(sessionTime)
			weeks := int(diff.Hours()) / 168
			days := (int(diff.Hours()) % 168) / 24
			hours := int(diff.Hours()) % 24
			minutes := int(diff.Minutes()) % 60
			seconds := int(diff.Seconds()) % 60

			units := []struct {
				value  int
				suffix string
			}{
				{weeks, "w"},
				{days, "d"},
				{hours, "h"},
				{minutes, "m"},
				{seconds, "s"},
			}

			var parts []string
			seenNonZero := false
			for _, unit := range units {
				if unit.value > 0 || seenNonZero {
					seenNonZero = true
					parts = append(parts, fmt.Sprintf("%d%s", unit.value, unit.suffix))
				}

				if unit.suffix == "s" && !seenNonZero {
					parts = append(parts, fmt.Sprintf("%d%s", unit.value, unit.suffix))
				}
			}

			formatted := fmt.Sprintf("Countdown to %s: %s", label, strings.Join(parts, " "))
			countdown.Set(formatted)
			time.Sleep(1 * time.Second)
		}
	}()
}

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
