package schedule

import (
	"fmt"
	"time"

	"f1tray/internal/notify"
)

// ScheduleWeeklyReminder sets a reminder every Wednesday at 12:00 PM local time.
// If testMode is true, it triggers in 10 seconds.
func ScheduleWeeklyReminder(testMode bool) {
	if testMode {
		fmt.Println("[TEST MODE] Scheduling weekly reminder in 10 seconds...")
		time.AfterFunc(10*time.Second, func() {
			notify.F1Reminder("F1 Weekly Reminder (Test)", "Check the schedule — F1 weekend might be coming!")
		})
		return
	}

	now := time.Now()
	next := nextWeeklyReminderTime(now)
	delay := time.Until(next)

	fmt.Println("Next weekly reminder scheduled for:", next)

	time.AfterFunc(delay, func() {
		notify.F1Reminder("F1 Weekly Reminder", "Check the schedule — F1 weekend might be coming!")
		go ScheduleWeeklyReminder(false) // reschedule for next week
	})
}

// nextWeeklyReminderTime returns the next Wednesday at 12:00 PM local time
func nextWeeklyReminderTime(t time.Time) time.Time {
	daysUntil := (int(time.Wednesday) - int(t.Weekday()) + 7) % 7
	if daysUntil == 0 && t.Hour() >= 12 {
		daysUntil = 7
	}
	next := t.AddDate(0, 0, daysUntil)
	return time.Date(next.Year(), next.Month(), next.Day(), 12, 0, 0, 0, next.Location())
}
