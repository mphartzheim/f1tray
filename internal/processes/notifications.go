package processes

import (
	"os/exec"
	"runtime"
	"time"

	"f1tray/internal/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// ShowInAppNotification sets the message and hides it after 5 seconds.
func ShowInAppNotification(label *widget.Label, wrapper fyne.CanvasObject, message string) {
	label.SetText(message)
	wrapper.Show()
	label.Show()

	time.AfterFunc(5*time.Second, func() {
		wrapper.Hide()
	})
}

// PlayNotificationSound plays a simple platform-specific notification sound.
func PlayNotificationSound(prefs config.Preferences) {
	if !prefs.EnableSound {
		return
	}

	switch runtime.GOOS {
	case "linux":
		_ = exec.Command("canberra-gtk-play", "--id", "message").Start()
	case "darwin":
		_ = exec.Command("afplay", "/System/Library/Sounds/Glass.aiff").Start()
	case "windows":
		_ = exec.Command("powershell", "-c", `[console]::beep(1000,300)`).Start()
	}
}
