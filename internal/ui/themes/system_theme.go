package themes

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type SystemTheme struct{}

func (c SystemTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	// Return a transparent color for the separator to remove grid lines.
	if name == theme.ColorNameSeparator {
		return color.NRGBA{R: 0, G: 0, B: 0, A: 0}
	}
	// Always force System variant
	return theme.DefaultTheme().Color(name, fyne.ThemeVariant(3))
}

func (c SystemTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (c SystemTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (c SystemTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
