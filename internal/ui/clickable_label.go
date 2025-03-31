package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ClickableLabel is a label widget that supports single and double tap interactions.
type ClickableLabel struct {
	widget.BaseWidget
	Text           string
	OnTapped       func()
	OnDoubleTapped func()
	Clickable      bool
	textColor      color.Color
	fontSize       float32
}

// NewClickableLabel creates and returns a new ClickableLabel with optional tap support.
func NewClickableLabel(text string, tapped func(), clickable bool) *ClickableLabel {
	// Use the current app's theme settings to get the foreground color and text size.
	currentApp := fyne.CurrentApp()
	initialColor := theme.Current().Color(theme.ColorNameForeground, currentApp.Settings().ThemeVariant())
	defaultFontSize := theme.TextSize()
	cl := &ClickableLabel{
		Text:      text,
		Clickable: clickable,
		OnTapped:  tapped,
		textColor: initialColor,
		fontSize:  defaultFontSize,
	}
	cl.ExtendBaseWidget(cl)
	return cl
}

// SetTextColor updates the text color of the label and refreshes it.
func (cl *ClickableLabel) SetTextColor(c color.Color) {
	cl.textColor = c
	cl.Refresh()
}

// SetFontSize updates the font size of the label and refreshes it.
func (cl *ClickableLabel) SetFontSize(size float32) {
	cl.fontSize = size
	cl.Refresh()
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

// CreateRenderer implements fyne.Widget and returns a custom renderer.
func (cl *ClickableLabel) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(cl.Text, cl.textColor)
	text.Alignment = fyne.TextAlignLeading
	text.TextSize = cl.fontSize
	return &clickableLabelRenderer{
		label:   cl,
		text:    text,
		objects: []fyne.CanvasObject{text},
	}
}

// clickableLabelRenderer is the custom renderer for ClickableLabel.
type clickableLabelRenderer struct {
	label   *ClickableLabel
	text    *canvas.Text
	objects []fyne.CanvasObject
}

func (r *clickableLabelRenderer) Layout(size fyne.Size) {
	r.text.Resize(size)
}

func (r *clickableLabelRenderer) MinSize() fyne.Size {
	return r.text.MinSize()
}

func (r *clickableLabelRenderer) Refresh() {
	// Update the text, color, and font size from the label.
	r.text.Text = r.label.Text
	r.text.Color = r.label.textColor
	r.text.TextSize = r.label.fontSize
	canvas.Refresh(r.label)
}

func (r *clickableLabelRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *clickableLabelRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *clickableLabelRenderer) Destroy() {}
