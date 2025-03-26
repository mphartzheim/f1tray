package processes

import (
	"os/exec"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// SendNotification sends a desktop notification using the Fyne app.
// Note: Desktop notifications will only work on supported platforms.
func SendNotification(a fyne.App, title, content string) {
	notif := fyne.NewNotification(title, content)
	a.SendNotification(notif)
}

// ShowInAppNotification sets the message and hides it after 5 seconds.
func ShowInAppNotification(label *widget.Label, wrapper fyne.CanvasObject, message string) {
	label.SetText(message)
	wrapper.Show()
	label.Show()

	time.AfterFunc(5*time.Second, func() {
		wrapper.Hide()
	})
}

func PlayNotificationSound() {
	switch runtime.GOOS {
	case "linux":
		_ = exec.Command("canberra-gtk-play", "--id", "message").Start()
	case "darwin":
		_ = exec.Command("afplay", "/System/Library/Sounds/Glass.aiff").Start()
	case "windows":
		_ = exec.Command("powershell", "-c", `[console]::beep(1000,300)`).Start()
	}
}
