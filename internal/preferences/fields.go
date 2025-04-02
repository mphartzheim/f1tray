package preferences

import (
	"github.com/mphartzheim/f1tray/internal/config"
	"github.com/mphartzheim/f1tray/internal/themes"

	"fyne.io/fyne/v2/widget"
)

func ThemeSelector(onChange func(string)) *widget.Select {
	prefs := config.Get()
	selectWidget := widget.NewSelect(themes.SortedThemeList(), func(selected string) {
		onChange(selected)
		prefs := config.Get()
		prefs.Themes.Selected = selected
		config.Set(prefs)
	})
	selectWidget.SetSelected(prefs.Themes.Selected)
	return selectWidget
}

func ClockToggle() *widget.Check {
	prefs := config.Get()
	toggle := widget.NewCheck("Use 24h Clock", func(checked bool) {
		prefs := config.Get()
		prefs.Clock.Use24Hour = checked
		config.Set(prefs)
	})
	toggle.SetChecked(prefs.Clock.Use24Hour)
	return toggle
}

func TrayToggle() *widget.Check {
	prefs := config.Get()
	toggle := widget.NewCheck("Minimize to Tray", func(checked bool) {
		prefs := config.Get()
		prefs.Window.MinimizeToTray = checked
		config.Set(prefs)
	})
	toggle.SetChecked(prefs.Window.MinimizeToTray)
	return toggle
}

func NotificationSlider() *widget.Slider {
	prefs := config.Get()
	slider := widget.NewSlider(0, 60)
	slider.Step = 5
	slider.OnChanged = func(val float64) {
		prefs := config.Get()
		prefs.Notifications.NotifyMinutesBefore = int(val)
		config.Set(prefs)
	}
	slider.SetValue(float64(prefs.Notifications.NotifyMinutesBefore))
	return slider
}
