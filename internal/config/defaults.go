package config

func defaultPreferences() Preferences {
	return Preferences{
		Themes: ThemesPreferences{
			Selected: "System",
		},
		Window: WindowPreferences{
			MinimizeToTray: true,
		},
		Clock: ClockPreferences{
			Use24Hour: true,
		},
		Notifications: NotificationPreferences{
			EnableSessionStart:  true,
			NotifyMinutesBefore: 10,
		},
		FavoriteDrivers:     []string{},
		FavoriteConstructor: "",
	}
}
