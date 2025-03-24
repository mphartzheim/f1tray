package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// BackButton returns a button that sets the given content when clicked.
func BackButton(label string, win fyne.Window, content fyne.CanvasObject) *widget.Button {
	return widget.NewButton(label, func() {
		win.SetContent(content)
	})
}
