package themes

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type LightTheme struct{}

func (c LightTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	if name == theme.ColorNamePrimary {
		return color.NRGBA{R: 0xFF, G: 0x18, B: 0x01, A: 0xFF} // F1 Red
	}
	// Return a transparent color for the separator to remove grid lines.
	if name == theme.ColorNameSeparator {
		return color.Transparent
	}
	// Always force light variant
	return theme.DefaultTheme().Color(name, fyne.ThemeVariant(1))
}

func (c LightTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (c LightTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (c LightTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
