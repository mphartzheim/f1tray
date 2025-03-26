package main

import (
	_ "embed"
	"runtime"
	"time"

	"f1tray/internal/config"
	"f1tray/internal/models"
	"f1tray/internal/processes"
	"f1tray/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

//go:embed assets/tray_icon.png
var trayIconBytes []byte

func main() {
	prefs := config.LoadConfig()

	myApp := app.NewWithID("f1tray")
	myWindow := myApp.NewWindow("F1 Viewer")

	// Create tab content
	scheduleTabData := ui.CreateScheduleTableTab(models.ScheduleURL, processes.ParseSchedule)
	upcomingTabData := ui.CreateUpcomingTab(models.UpcomingURL, processes.ParseUpcoming)
	resultsTabData := ui.CreateResultsTableTab(models.RaceResultsURL, processes.ParseRaceResults)
	qualifyingTabData := ui.CreateResultsTableTab(models.QualifyingURL, processes.ParseQualifyingResults)
	sprintTabData := ui.CreateResultsTableTab(models.SprintURL, processes.ParseSprintResults)
	preferencesContent := ui.CreatePreferencesTab(prefs, func(updated config.Preferences) {
		_ = config.SaveConfig(updated)
		prefs = updated
	})

	// Create floating notification with close button
	notificationLabel := widget.NewLabel("")
	notificationLabel.Alignment = fyne.TextAlignCenter

	notificationWrapper := container.NewWithoutLayout() // Defined early so we can reference it in the close button

	closeButton := widget.NewButton("✕", func() {
		notificationWrapper.Hide()
	})
	closeButton.Importance = widget.LowImportance

	notificationRow := container.NewHBox(
		notificationLabel,
		layout.NewSpacer(),
		closeButton,
	)

	popup := container.NewPadded(notificationRow)

	// Background layer
	popupBG := canvas.NewRectangle(theme.Color(theme.ColorNamePrimary))
	popupBG.SetMinSize(fyne.NewSize(320, 50))

	// Stack them: border > background > content
	notificationContainer := container.NewStack(popupBG, popup)
	notificationWrapper = container.NewCenter(notificationContainer)
	notificationWrapper.Hide()

	// Create tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("Schedule", scheduleTabData.Content),
		container.NewTabItem("Upcoming", upcomingTabData.Content),
		container.NewTabItem("Race Results", resultsTabData.Content),
		container.NewTabItem("Qualifying", qualifyingTabData.Content),
		container.NewTabItem("Sprint", sprintTabData.Content),
		container.NewTabItem("Preferences", preferencesContent),
	)

	// Stack tabs and floating notification
	stack := container.NewStack(
		tabs,
		notificationWrapper,
	)

	// Manual refresh button
	refreshButton := widget.NewButton("Refresh All Data", func() {
		refreshAllData(notificationLabel, notificationWrapper,
			scheduleTabData, upcomingTabData, resultsTabData, qualifyingTabData, sprintTabData)
	})

	// Auto-refresh every hour
	go func() {
		var refreshInterval time.Duration
		if prefs.DebugMode {
			refreshInterval = 1 * time.Minute
		} else {
			refreshInterval = 1 * time.Hour
		}

		ticker := time.NewTicker(refreshInterval)
		defer ticker.Stop()

		for range ticker.C {
			refreshAllData(notificationLabel, notificationWrapper,
				scheduleTabData, upcomingTabData, resultsTabData, qualifyingTabData, sprintTabData)
		}
	}()

	// Set window content
	myWindow.SetContent(container.NewBorder(
		refreshButton, nil, nil, nil, stack,
	))
	myWindow.Resize(fyne.NewSize(900, 600))

	// Setup tray
	iconResource := fyne.NewStaticResource("tray_icon.png", trayIconBytes)
	if desk, ok := myApp.(desktop.App); ok {
		go func() {
			maxAttempts := 5
			success := false

			if runtime.GOOS == "windows" {
				for i := 0; i < maxAttempts; i++ {
					func() {
						defer func() {
							if r := recover(); r != nil {
								// Optionally log panic info
							}
						}()

						desk.SetSystemTrayIcon(iconResource)
						success = true
					}()

					if success {
						break
					}

					println("[F1Tray] Attempt", i+1, "to set system tray icon failed. Retrying...")
					time.Sleep(2 * time.Second)
				}

				if !success {
					println("[F1Tray] Failed to set system tray icon after 5 attempts. Exiting.")
					myApp.Quit()
					return
				}
			} else {
				desk.SetSystemTrayIcon(iconResource)
			}

			// Tray icon was set successfully; now set the menu
			desk.SetSystemTrayMenu(fyne.NewMenu("F1 Tray",
				fyne.NewMenuItem("Schedule", func() { tabs.SelectIndex(0); myWindow.Show(); myWindow.RequestFocus() }),
				fyne.NewMenuItem("Upcoming", func() { tabs.SelectIndex(1); myWindow.Show(); myWindow.RequestFocus() }),
				fyne.NewMenuItem("Race Results", func() { tabs.SelectIndex(2); myWindow.Show(); myWindow.RequestFocus() }),
				fyne.NewMenuItem("Qualifying", func() { tabs.SelectIndex(3); myWindow.Show(); myWindow.RequestFocus() }),
				fyne.NewMenuItem("Sprint", func() { tabs.SelectIndex(4); myWindow.Show(); myWindow.RequestFocus() }),
				fyne.NewMenuItem("Preferences", func() { tabs.SelectIndex(5); myWindow.Show(); myWindow.RequestFocus() }),
				fyne.NewMenuItemSeparator(),
				fyne.NewMenuItem("Show", func() { myWindow.Show(); myWindow.RequestFocus() }),
				fyne.NewMenuItem("Quit", myApp.Quit),
			))
		}()
	}

	// Window visibility
	if prefs.HideOnOpen {
		myWindow.Hide()
	} else {
		myWindow.Show()
	}

	// Close intercept
	myWindow.SetCloseIntercept(func() {
		if prefs.CloseBehavior == "exit" {
			myApp.Quit()
		} else {
			myWindow.Hide()
		}
	})

	myApp.Run()
}

func refreshAllData(label *widget.Label, wrapper fyne.CanvasObject, tabs ...models.TabData) {
	updated := false
	for _, tab := range tabs {
		if tab.Refresh() {
			updated = true
		}
	}

	if updated {
		processes.PlayNotificationSound()
		processes.ShowInAppNotification(label, wrapper, "Data has been refreshed.")
	} else {
		processes.ShowInAppNotification(label, wrapper, "No new data to load.")
	}

}
