package themes

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// SystemTheme defines a custom system theme implementation for the app.
type SystemTheme struct{}

// Color returns the color for the given theme element, with custom primary and separator overrides.
func (c SystemTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	// Return a transparent color for the separator to remove grid lines.
	if name == theme.ColorNameSeparator {
		return color.NRGBA{R: 0, G: 0, B: 0, A: 0}
	}
	// Always force System variant
	return theme.DefaultTheme().Color(name, fyne.ThemeVariant(3))
}

// Font returns the font resource for the specified text style.
func (c SystemTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Icon returns the icon resource for the given icon name.
func (c SystemTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size returns the size value for the specified theme size name.
func (c SystemTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
