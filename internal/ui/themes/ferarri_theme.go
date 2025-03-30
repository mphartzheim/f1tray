package themes

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// FerrariTheme defines a custom theme implementation with Ferrari's signature Rosso Corsa red.
type FerrariTheme struct{}

// Color returns the color for the given theme element, with custom primary and separator overrides.
func (c FerrariTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	if name == theme.ColorNamePrimary {
		return color.NRGBA{R: 0xDC, G: 0x00, B: 0x00, A: 0xFF} // Ferrari Rosso Corsa
	}
	// Return a transparent color for the separator to remove grid lines.
	if name == theme.ColorNameSeparator {
		return color.Transparent
	}
	// Always force dark variant
	return theme.DefaultTheme().Color(name, fyne.ThemeVariant(0))
}

// Font returns the font resource for the specified text style.
func (c FerrariTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Icon returns the icon resource for the given icon name.
func (c FerrariTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size returns the size value for the specified theme size name.
func (c FerrariTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
