package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CreateNotification creates a notification overlay container that is centered on the given window's canvas.
func CreateNotification(win fyne.Window) (*widget.Label, fyne.CanvasObject) {
	notificationLabel := widget.NewLabel("")
	notificationLabel.Alignment = fyne.TextAlignCenter

	// Build a padded horizontal container with the label and close button.
	popup := container.NewPadded(container.NewHBox(
		notificationLabel,
	))

	// Create a background rectangle.
	popupBG := canvas.NewRectangle(theme.Color(theme.ColorNamePrimary))
	popupBG.SetMinSize(fyne.NewSize(320, 50))

	// Stack the background and popup.
	notificationContainer := container.NewStack(popupBG, popup)
	// Center the notification container.
	centered := container.NewCenter(notificationContainer)
	// Wrap in a Max container so that it fills the canvas.
	wrapper := container.NewStack(centered)
	wrapper.Hide() // start hidden

	// We need the overlay reference in the close button callback.
	var overlay fyne.CanvasObject = wrapper
	return notificationLabel, overlay
}

// ShowNotification creates and displays the notification overlay (centered) on the provided window.
// The overlay is sized to fill the window and is automatically removed after 3 seconds.
func ShowNotification(win fyne.Window, text string) fyne.CanvasObject {
	label, overlay := CreateNotification(win)
	label.SetText(text)
	// Ensure the overlay fills the canvas.
	overlay.Resize(win.Canvas().Size())
	overlay.Show()
	win.Canvas().Overlays().Add(overlay)

	// Auto-remove the overlay after 3 seconds.
	go func() {
		time.Sleep(3 * time.Second)
		win.Canvas().Overlays().Remove(overlay)
	}()
	return overlay
}
