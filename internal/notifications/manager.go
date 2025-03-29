package notifications

import (
	"fmt"
	"strings"
	"time"

	"f1tray/internal/config"
	"f1tray/internal/models"

	"fyne.io/fyne/v2"
)

// SessionType represents a type of race session.
type SessionType string

const (
	Practice   SessionType = "Practice"
	Qualifying SessionType = "Qualifying"
	Race       SessionType = "Race"
)

// notified tracks sessions that have already been notified (for production).
var notified = map[string]bool{}

// CheckAndSendNotifications examines a session and sends a notification if conditions are met.
// If the session is a test (its Label contains "Test"), it always sends a notification
// without marking the session as notified.
func CheckAndSendNotifications(session models.SessionInfo) {
	// Convert session.Type (a string) into our SessionType.
	sType := ParseSessionType(session.Type)
	prefs := config.Get()

	var settings *config.SessionNotificationSettings
	switch sType {
	case Practice:
		settings = prefs.Notifications.Practice
	case Qualifying:
		settings = prefs.Notifications.Qualifying
	case Race:
		settings = prefs.Notifications.Race
	default:
		return
	}

	// Generate keys for production notifications.
	keyStart := fmt.Sprintf("%s_start", sType)
	keyBefore := fmt.Sprintf("%s_before", sType)

	now := time.Now()

	// If this is a test session, bypass updating the notified map.
	isTest := strings.Contains(session.Label, "Test")

	// BEFORE SESSION notification.
	if settings.NotifyBefore {
		offset := time.Duration(settings.BeforeValue)
		if settings.BeforeUnit == "hours" {
			offset *= time.Hour
		} else {
			offset *= time.Minute
		}
		targetTime := session.StartTime.Add(-offset)
		if now.After(targetTime) && now.Before(session.StartTime) {
			send(session.Label, fmt.Sprintf("Starting in %d %s", settings.BeforeValue, settings.BeforeUnit), settings.PlaySoundBefore)
			if !isTest {
				notified[keyBefore] = true
			}
			// Return after sending notification to avoid sending both before and start notifications.
			return
		}
	}

	// AT SESSION START notification.
	if settings.NotifyOnStart {
		if now.After(session.StartTime) && now.Before(session.StartTime.Add(1*time.Minute)) {
			send(session.Label, "Session has started!", settings.PlaySoundOnStart)
			if !isTest {
				notified[keyStart] = true
			}
		}
	}
}

// ParseSessionType converts a string to a SessionType.
func ParseSessionType(t string) SessionType {
	switch strings.ToLower(t) {
	case "practice":
		return Practice
	case "qualifying":
		return Qualifying
	case "race":
		return Race
	default:
		return SessionType("unknown")
	}
}

// send triggers a desktop notification and plays a sound if requested.
func send(title, message string, withSound bool) {
	fyne.CurrentApp().SendNotification(&fyne.Notification{
		Title:   title,
		Content: message,
	})
	if withSound {
		PlayNotificationSound()
	}
}
