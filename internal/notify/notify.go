package notify

import (
	"path/filepath"

	"github.com/gen2brain/beeep"
)

func F1Reminder(title, message string) error {
	iconPath := filepath.Join("assets", "tray_icon.png")
	return beeep.Notify(title, message, iconPath)
}
