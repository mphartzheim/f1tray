package themes

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type DarkTheme struct{}

func (c DarkTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	if name == theme.ColorNamePrimary {
		return color.NRGBA{R: 0xFF, G: 0x18, B: 0x01, A: 0xFF} // F1 Red
	}
	// Return a transparent color for the separator to remove grid lines.
	if name == theme.ColorNameSeparator {
		return color.Transparent
	}
	// Always force dark variant
	return theme.DefaultTheme().Color(name, fyne.ThemeVariant(0))
}

func (c DarkTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (c DarkTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (c DarkTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
