package processes

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

// SendNotification sends a desktop notification using the Fyne app.
// Note: Desktop notifications will only work on supported platforms.
func SendNotification(a fyne.App, title, content string) {
	notif := fyne.NewNotification(title, content)
	a.SendNotification(notif)
}

// ShowInAppNotification displays an in-app dialog notification.
func ShowInAppNotification(w fyne.Window, title, content string) {
	dialog.ShowInformation(title, content, w)
}
