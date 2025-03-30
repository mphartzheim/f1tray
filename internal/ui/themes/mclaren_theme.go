package themes

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// McLarenTheme defines a custom theme implementation with McLaren's signature papaya orange.
type McLarenTheme struct{}

// Color returns the color for the given theme element, with custom primary and separator overrides.
func (c McLarenTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	if name == theme.ColorNamePrimary {
		return color.NRGBA{R: 0xFF, G: 0x87, B: 0x00, A: 0xFF} // McLaren Papaya Orange
	}
	// Return a transparent color for the separator to remove grid lines.
	if name == theme.ColorNameSeparator {
		return color.Transparent
	}
	// Always force dark variant
	return theme.DefaultTheme().Color(name, fyne.ThemeVariant(0))
}

// Font returns the font resource for the specified text style.
func (c McLarenTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Icon returns the icon resource for the given icon name.
func (c McLarenTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size returns the size value for the specified theme size name.
func (c McLarenTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
