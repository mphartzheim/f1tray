package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CreateNotification creates and returns a notification label and its wrapper container.
// It encapsulates the notification overlay creation.
func CreateNotification() (*widget.Label, fyne.CanvasObject) {
	notificationLabel := widget.NewLabel("")
	notificationLabel.Alignment = fyne.TextAlignCenter

	// Temporary container to build the notification overlay
	notificationWrapper := container.NewWithoutLayout()

	// Close button for hiding the notification overlay
	closeButton := widget.NewButton("✕", func() {
		notificationWrapper.Hide()
	})
	closeButton.Importance = widget.LowImportance

	// Build the popup container with the label and close button
	popup := container.NewPadded(container.NewHBox(
		notificationLabel, layout.NewSpacer(), closeButton,
	))

	// Background for the popup
	popupBG := canvas.NewRectangle(theme.Color(theme.ColorNamePrimary))
	popupBG.SetMinSize(fyne.NewSize(320, 50))
	notificationContainer := container.NewStack(popupBG, popup)
	notificationWrapper = container.NewCenter(notificationContainer)
	notificationWrapper.Hide()

	return notificationLabel, notificationWrapper
}
