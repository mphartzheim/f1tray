package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// ClickableLabel is a label with a double-click handler.
type ClickableLabel struct {
	widget.Label
	OnDoubleTapped func()
	lastTap        time.Time
}

// NewClickableLabel creates a new ClickableLabel.
func NewClickableLabel(text string, onDoubleTapped func()) *ClickableLabel {
	cl := &ClickableLabel{
		Label:          *widget.NewLabel(text),
		OnDoubleTapped: onDoubleTapped,
	}
	return cl
}

// Tapped simulates a double-click.
func (c *ClickableLabel) Tapped(_ *fyne.PointEvent) {
	now := time.Now()
	if now.Sub(c.lastTap) <= 500*time.Millisecond {
		// Detected double-tap
		if c.OnDoubleTapped != nil {
			c.OnDoubleTapped()
		}
	}
	c.lastTap = now
}

// TappedSecondary satisfies the Tappable interface.
func (c *ClickableLabel) TappedSecondary(_ *fyne.PointEvent) {}
