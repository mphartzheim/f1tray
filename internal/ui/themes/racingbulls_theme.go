package themes

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// RacingBullsTheme defines a custom theme implementation for Racing Bulls, with deep navy tones.
type RacingBullsTheme struct{}

// Color returns the color for the given theme element, with custom primary and separator overrides.
func (c RacingBullsTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	if name == theme.ColorNamePrimary {
		return color.NRGBA{R: 0x00, G: 0x20, B: 0x5B, A: 0xFF} // Racing Bulls Blue
	}
	// Return a transparent color for the separator to remove grid lines.
	if name == theme.ColorNameSeparator {
		return color.Transparent
	}
	// Always force dark variant
	return theme.DefaultTheme().Color(name, fyne.ThemeVariant(1))
}

// Font returns the font resource for the specified text style.
func (c RacingBullsTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Icon returns the icon resource for the given icon name.
func (c RacingBullsTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size returns the size value for the specified theme size name.
func (c RacingBullsTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
