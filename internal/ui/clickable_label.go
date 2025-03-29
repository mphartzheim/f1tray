package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// ClickableLabel is a label widget that supports single and double tap interactions.
type ClickableLabel struct {
	widget.Label
	// Callback for a single tap event.
	OnTapped func()
	// Callback for a double tap event.
	OnDoubleTapped func()
	// Indicates if the label should be interactive.
	Clickable bool
}

// NewClickableLabel creates and returns a new ClickableLabel with optional tap support.
func NewClickableLabel(text string, tapped func(), clickable bool) *ClickableLabel {
	cl := &ClickableLabel{Clickable: clickable}
	cl.Text = text
	cl.OnTapped = tapped
	cl.ExtendBaseWidget(cl)
	return cl
}

// Tapped handles single-tap events on the ClickableLabel.
func (cl *ClickableLabel) Tapped(_ *fyne.PointEvent) {
	if cl.OnTapped != nil {
		cl.OnTapped()
	}
}

// DoubleTapped handles double-tap events on the ClickableLabel.
func (cl *ClickableLabel) DoubleTapped(_ *fyne.PointEvent) {
	if cl.OnDoubleTapped != nil {
		cl.OnDoubleTapped()
	}
}

// MouseIn is called when the mouse pointer enters the ClickableLabel (unused).
func (cl *ClickableLabel) MouseIn(me *desktop.MouseEvent) {}

// MouseOut is called when the mouse pointer exits the ClickableLabel (unused).
func (cl *ClickableLabel) MouseOut() {}

// MouseMoved is called when the mouse moves over the ClickableLabel (unused).
func (cl *ClickableLabel) MouseMoved(me *desktop.MouseEvent) {}

// Cursor returns the appropriate cursor based on the Clickable state.
func (cl *ClickableLabel) Cursor() desktop.Cursor {
	if cl.Clickable {
		return desktop.PointerCursor
	}
	return desktop.DefaultCursor
}
