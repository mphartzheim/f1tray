package processes

import (
	"runtime"
	"time"

	"f1tray/internal/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func RefreshAllData(state models.AppState, label *widget.Label, wrapper fyne.CanvasObject, silent bool, tabs ...models.TabData) {
	updated := false
	for _, tab := range tabs {
		if state.DebugMode || tab.Refresh() {
			updated = true
		}
	}
	if updated {
		PlayNotificationSound()
		if !silent && label != nil && wrapper != nil {
			ShowInAppNotification(label, wrapper, "Data has been refreshed.")
		}
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "F1Tray",
			Content: "New F1 data is available!",
		})
	} else if !silent && label != nil && wrapper != nil {
		ShowInAppNotification(label, wrapper, "No new data to load.")
	}
}

func StartAutoRefresh(state models.AppState, label *widget.Label, wrapper fyne.CanvasObject, tabs ...models.TabData) {
	interval := time.Hour
	if state.DebugMode {
		interval = time.Minute
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		RefreshAllData(state, label, wrapper, false, tabs...)
	}
}

func SetTrayIcon(desk desktop.App, icon fyne.Resource, tabs *container.AppTabs, win fyne.Window) {
	maxAttempts := 5
	success := false

	if runtime.GOOS == "windows" {
		for i := 0; i < maxAttempts; i++ {
			func() {
				defer func() { recover() }()
				desk.SetSystemTrayIcon(icon)
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
			fyne.CurrentApp().Quit()
			return
		}
	} else {
		desk.SetSystemTrayIcon(icon)
	}

	desk.SetSystemTrayMenu(fyne.NewMenu("F1 Tray",
		fyne.NewMenuItem("Schedule", func() { tabs.SelectIndex(0); win.Show(); win.RequestFocus() }),
		fyne.NewMenuItem("Upcoming", func() { tabs.SelectIndex(1); win.Show(); win.RequestFocus() }),
		fyne.NewMenuItem("Race Results", func() { tabs.SelectIndex(2); win.Show(); win.RequestFocus() }),
		fyne.NewMenuItem("Qualifying", func() { tabs.SelectIndex(3); win.Show(); win.RequestFocus() }),
		fyne.NewMenuItem("Sprint", func() { tabs.SelectIndex(4); win.Show(); win.RequestFocus() }),
		fyne.NewMenuItem("Preferences", func() { tabs.SelectIndex(5); win.Show(); win.RequestFocus() }),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Show", func() { tabs.SelectIndex(0); win.Show(); win.RequestFocus() }),
		fyne.NewMenuItem("Quit", fyne.CurrentApp().Quit),
	))
}

func AppendSessionRow(rows [][]string, label, date, time string, use24h bool) [][]string {
	if date != "" && time != "" {
		d, t := Localize(date, time, use24h)
		rows = append(rows, []string{label, d, t})
	}
	return rows
}
