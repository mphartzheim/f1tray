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
// If withLoading is true, an infinite progress bar is added.
func CreateNotification(win fyne.Window, withLoading bool) (*widget.Label, fyne.CanvasObject) {
	notificationLabel := widget.NewLabel("")
	notificationLabel.Alignment = fyne.TextAlignCenter

	// Build the content container.
	contentItems := []fyne.CanvasObject{notificationLabel}
	var progressBar *widget.ProgressBarInfinite
	if withLoading {
		progressBar = widget.NewProgressBarInfinite()
		contentItems = append(contentItems, progressBar)
	}

	// Use a VBox to stack the label (and progress bar, if any) vertically.
	popupContent := container.NewVBox(contentItems...)
	popup := container.NewPadded(popupContent)

	// Create a background rectangle.
	popupBG := canvas.NewRectangle(theme.Color(theme.ColorNamePrimary))
	popupBG.SetMinSize(fyne.NewSize(320, 50))

	// Stack the background and popup.
	notificationContainer := container.NewStack(popupBG, popup)
	// Center the notification container.
	centered := container.NewCenter(notificationContainer)
	// Wrap in a container that fills the canvas.
	wrapper := container.NewStack(centered)
	wrapper.Hide() // start hidden

	return notificationLabel, wrapper
}

// ShowNotification creates and displays the notification overlay (centered) on the provided window.
// If withLoading is true, an infinite loading bar is displayed and the overlay is not auto-removed.
// Otherwise, the overlay is automatically removed after 3 seconds.
func ShowNotification(win fyne.Window, text string, withLoading bool) fyne.CanvasObject {
	label, overlay := CreateNotification(win, withLoading)
	label.SetText(text)
	overlay.Resize(win.Canvas().Size())
	overlay.Show()
	win.Canvas().Overlays().Add(overlay)

	if !withLoading {
		// Auto-remove the overlay after 3 seconds when not in loading mode.
		go func() {
			time.Sleep(3 * time.Second)
			win.Canvas().Overlays().Remove(overlay)
		}()
	}

	return overlay
}
