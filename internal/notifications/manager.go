package notifications

import (
	"fmt"
	"strings"
	"time"

	"f1tray/internal/config"

	"fyne.io/fyne/v2"
)

type SessionType string

const (
	Practice   SessionType = "Practice"
	Qualifying SessionType = "Qualifying"
	Race       SessionType = "Race"
)

type SessionInfo struct {
	Type      SessionType
	StartTime time.Time
	Title     string
}

// Track notifications already sent.
var notified = map[string]bool{}

func CheckAndSendNotifications(session SessionInfo) {
	prefs := config.Get()

	var settings *config.SessionNotificationSettings
	switch session.Type {
	case Practice:
		settings = prefs.Notifications.Practice
	case Qualifying:
		settings = prefs.Notifications.Qualifying
	case Race:
		settings = prefs.Notifications.Race
	default:
		return
	}

	keyStart := fmt.Sprintf("%s_start", session.Type)
	keyBefore := fmt.Sprintf("%s_before", session.Type)

	now := time.Now()

	// BEFORE SESSION
	if settings.NotifyBefore && !notified[keyBefore] {
		offset := time.Duration(settings.BeforeValue)
		if settings.BeforeUnit == "hours" {
			offset *= time.Hour
		} else {
			offset *= time.Minute
		}
		targetTime := session.StartTime.Add(-offset)

		if now.After(targetTime) && now.Before(session.StartTime) {
			send(session.Title, fmt.Sprintf("Starting in %d %s", settings.BeforeValue, settings.BeforeUnit), settings.PlaySoundBefore)
			notified[keyBefore] = true
		}
	}

	// AT SESSION START
	if settings.NotifyOnStart && !notified[keyStart] {
		if now.After(session.StartTime) && now.Before(session.StartTime.Add(1*time.Minute)) {
			send(session.Title, "Session has started!", settings.PlaySoundOnStart)
			notified[keyStart] = true
		}
	}
}

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

func send(title, message string, withSound bool) {
	fyne.CurrentApp().SendNotification(&fyne.Notification{
		Title:   title,
		Content: message,
	})
	if withSound {
		PlayNotificationSound()
	}
}
