package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// ClickableCell is a custom widget that displays a text label and supports clicking and hovering.
type ClickableCell struct {
	widget.BaseWidget
	Label      *widget.Label
	IsPointer  bool
	OnTapped   func()
	OnMouseIn  func()
	OnMouseOut func()
}

// NewClickableCell returns a new instance of ClickableCell.
func NewClickableCell() *ClickableCell {
	cc := &ClickableCell{
		Label: widget.NewLabel(""),
	}
	cc.ExtendBaseWidget(cc)
	return cc
}

// CreateRenderer implements fyne.Widget.
func (cc *ClickableCell) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(cc.Label)
}

// Tapped is called when the cell is clicked.
func (cc *ClickableCell) Tapped(ev *fyne.PointEvent) {
	if cc.OnTapped != nil {
		cc.OnTapped()
	}
}

// TappedSecondary can be implemented if needed.
func (cc *ClickableCell) TappedSecondary(ev *fyne.PointEvent) {}

// MouseIn implements desktop.Hoverable.
func (cc *ClickableCell) MouseIn(ev *desktop.MouseEvent) {
	if cc.OnMouseIn != nil {
		cc.OnMouseIn()
	}
}

// MouseOut implements desktop.Hoverable.
func (cc *ClickableCell) MouseOut() {
	if cc.OnMouseOut != nil {
		cc.OnMouseOut()
	}
}

// MouseMoved implements desktop.Hoverable (required for interface).
func (cc *ClickableCell) MouseMoved(ev *desktop.MouseEvent) {}

// Cursor returns the cursor to be displayed when hovering over the cell.
func (cc *ClickableCell) Cursor() desktop.Cursor {
	if cc.IsPointer {
		return desktop.PointerCursor
	}
	return desktop.DefaultCursor
}
