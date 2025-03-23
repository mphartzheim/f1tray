package schedule

import (
	"fmt"
	"strings"
	"time"

	"f1tray/internal/api"
	"f1tray/internal/notify"
)

// ScheduleNextRaceReminder sets a timer to notify N hours before the next race.
func ScheduleNextRaceReminder(testMode bool, hoursBefore int) {
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
		fmt.Println("[TEST MODE] Scheduling race reminder in 10 seconds...")
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

	reminderTime := raceTime.Add(-1 * time.Duration(hoursBefore) * time.Hour)
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
			notify.F1Reminder("F1 Reminder", message)
		})
	} else {
		fmt.Println("Too late to schedule reminder: race is too soon or already started.")
	}
}

// ScheduleWeeklyReminder sets a reminder every week at the configured day and hour.
func ScheduleWeeklyReminder(testMode bool, day string, hour int) {
	if testMode {
		fmt.Println("[TEST MODE] Scheduling weekly reminder in 10 seconds...")
		time.AfterFunc(10*time.Second, func() {
			notify.F1Reminder("F1 Weekly Reminder (Test)", "Check the schedule — F1 weekend might be coming!")
		})
		return
	}

	now := time.Now()
	next := nextWeeklyReminderTime(now, day, hour)
	delay := time.Until(next)

	fmt.Printf("Next weekly reminder scheduled for: %s\n", next.Format(time.RFC1123))

	time.AfterFunc(delay, func() {
		notify.F1Reminder("F1 Weekly Reminder", "Check the schedule — F1 weekend might be coming!")
		go ScheduleWeeklyReminder(false, day, hour) // reschedule for next week
	})
}

// nextWeeklyReminderTime returns the next configured weekday and hour.
func nextWeeklyReminderTime(t time.Time, day string, hour int) time.Time {
	weekdayMap := map[string]time.Weekday{
		"sunday":    time.Sunday,
		"monday":    time.Monday,
		"tuesday":   time.Tuesday,
		"wednesday": time.Wednesday,
		"thursday":  time.Thursday,
		"friday":    time.Friday,
		"saturday":  time.Saturday,
	}

	wd, ok := weekdayMap[strings.ToLower(day)]
	if !ok {
		wd = time.Wednesday
	}

	daysUntil := (int(wd) - int(t.Weekday()) + 7) % 7
	if daysUntil == 0 && t.Hour() >= hour {
		daysUntil = 7
	}
	next := t.AddDate(0, 0, daysUntil)
	return time.Date(next.Year(), next.Month(), next.Day(), hour, 0, 0, 0, next.Location())
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
