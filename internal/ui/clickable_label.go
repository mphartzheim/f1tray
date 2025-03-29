package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// ClickableLabel is a custom widget that supports both single and double taps,
// and can indicate whether it is clickable.
type ClickableLabel struct {
	widget.Label
	// Callback for a single tap event.
	OnTapped func()
	// Callback for a double tap event.
	OnDoubleTapped func()
	// Indicates if the label should be interactive.
	Clickable bool
}

// NewClickableLabel creates a new ClickableLabel.
func NewClickableLabel(text string, tapped func(), clickable bool) *ClickableLabel {
	cl := &ClickableLabel{Clickable: clickable}
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

// Implement desktop.Hoverable (empty implementations if not needed).
func (cl *ClickableLabel) MouseIn(me *desktop.MouseEvent)    {}
func (cl *ClickableLabel) MouseOut()                         {}
func (cl *ClickableLabel) MouseMoved(me *desktop.MouseEvent) {}

// Implement the desktop.Cursorable interface.
// Returns a pointer cursor if Clickable is true, otherwise the default cursor.
func (cl *ClickableLabel) Cursor() desktop.Cursor {
	if cl.Clickable {
		return desktop.PointerCursor
	}
	return desktop.DefaultCursor
}
