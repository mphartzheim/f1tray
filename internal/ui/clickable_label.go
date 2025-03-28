package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// ClickableLabel is a custom widget that supports both single and double taps.
type ClickableLabel struct {
	widget.Label
	// Callback for a single tap event.
	OnTapped func()
	// Callback for a double tap event.
	OnDoubleTapped func()
}

// NewClickableLabel creates a new ClickableLabel.
func NewClickableLabel(text string, tapped func()) *ClickableLabel {
	cl := &ClickableLabel{}
	cl.Text = text
	cl.OnTapped = tapped
	cl.ExtendBaseWidget(cl)
	return cl
}

// Tapped is called when the label is tapped.
func (cl *ClickableLabel) Tapped(_ *fyne.PointEvent) {
	if cl.OnTapped != nil {
		cl.OnTapped()
	}
}

// DoubleTapped is called when the label is double tapped.
func (cl *ClickableLabel) DoubleTapped(_ *fyne.PointEvent) {
	if cl.OnDoubleTapped != nil {
		cl.OnDoubleTapped()
	}
}
