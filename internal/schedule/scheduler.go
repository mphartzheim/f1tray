package schedule

import (
	"fmt"
	"time"

	"f1tray/internal/api"
	"f1tray/internal/notify"
)

// ScheduleNextRaceReminder sets a timer to notify 1 hour before the next race.
func ScheduleNextRaceReminder(testMode bool) {
	race, err := api.FetchNextRace()
	if err != nil {
		fmt.Println("Could not schedule next race reminder:", err)
		return
	}

	layout := "2006-01-02 15:04:05Z"
	dateTimeStr := race.Date + " " + race.Time
	raceTime, err := time.Parse(layout, dateTimeStr)
	if err != nil {
		fmt.Println("Failed to parse race time:", err)
		return
	}

	if testMode {
		fmt.Println("[TEST MODE] Scheduling reminder in 10 seconds...")
		time.AfterFunc(10*time.Second, func() {
			message := fmt.Sprintf(
				"[TEST] Upcoming race: %s at %s, %s",
				race.RaceName,
				race.Circuit.Location.Locality,
				race.Circuit.Location.Country,
			)
			notify.F1Reminder("F1 Reminder (Test Mode)", message)
		})
		return
	}

	reminderTime := raceTime.Add(-1 * time.Hour)
	delay := time.Until(reminderTime)

	if delay > 0 {
		fmt.Println("Scheduling race reminder in:", delay)
		time.AfterFunc(delay, func() {
			message := fmt.Sprintf(
				"Upcoming race: %s at %s, %s",
				race.RaceName,
				race.Circuit.Location.Locality,
				race.Circuit.Location.Country,
			)
			notify.F1Reminder("F1 Reminder (1 hour to go!)", message)
		})
	} else {
		fmt.Println("Too late to schedule reminder: race is too soon or already started.")
	}
}

// TestRaceNotification sends an immediate API call and alert for testing from the tray menu.
func TestRaceNotification() {
	race, err := api.FetchNextRace()
	if err != nil {
		notify.F1Reminder("F1 API Error", err.Error())
		fmt.Println("API error:", err)
		return
	}

	message := fmt.Sprintf(
		"Next race: %s\nCircuit: %s\nLocation: %s, %s\nDate: %s",
		race.RaceName,
		race.Circuit.CircuitName,
		race.Circuit.Location.Locality,
		race.Circuit.Location.Country,
		race.Date,
	)

	notify.F1Reminder("F1 Next Race", message)
}
