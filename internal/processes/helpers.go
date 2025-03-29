package processes

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"f1tray/internal/config"
	"f1tray/internal/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// RefreshAllData refreshes all provided tabs and respects debug mode for forced refresh.
func RefreshAllData(label *widget.Label, wrapper fyne.CanvasObject, tabs ...models.TabData) {
	for _, tab := range tabs {
		if config.Get().DebugMode || tab.Refresh() {
		}
	}
}

// StartAutoRefresh checks an endpoint's hash on intervals and notifies the user if it changes after the first run.
func StartAutoRefresh(state *models.AppState, selectedYear string) {
	// Download and store the initial aggregated hash from your selected endpoints.
	prevHash, err := DownloadDataHash(selectedYear) // Assumes this function fetches a combined hash
	if err != nil {
		log.Println("Error downloading initial data hash:", err)
		// Optionally, handle the error (retry or exit)
	}

	// Determine refresh interval: 1 hour normally, 1 minute in debug mode.
	interval := time.Hour
	if config.Get().DebugMode {
		interval = time.Minute
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		// Download the current aggregated hash.
		currHash, err := DownloadDataHash(selectedYear)
		if err != nil {
			log.Println("Error downloading current data hash:", err)
			continue // Skip this tick if there's an error
		}

		// Compare the current hash with the previously stored hash.
		if currHash != prevHash {
			// Data has changed.
			// Only send notifications if it's not the first run.
			if !state.FirstRun {
				// Send a system notification via the Fyne framework.
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title:   "F1Tray",
					Content: "New F1 data is available!",
				})

				// Optionally play a sound based on user preferences.
				PlayNotificationSound()

			}
			// Update the stored hash with the new one.
			prevHash = currHash
			// After the first update, mark the first run as complete.
			state.FirstRun = false
		}
	}
}

// DownloadDataHash fetches data from endpoints, combines it, and returns a SHA-256 hash as a hex string.
func DownloadDataHash(selectedYear string) (string, error) {
	// Define your endpoints here; update these URLs as needed for your application.
	endpoints := []string{
		fmt.Sprintf(models.RaceResultsURL, selectedYear, "last"),
		fmt.Sprintf(models.QualifyingURL, selectedYear, "last"),
		fmt.Sprintf(models.SprintURL, selectedYear, "last"),
	}

	// Initialize a SHA-256 hash.
	hasher := sha256.New()

	// Loop through each endpoint, download its data, and write it into the hash.
	for _, url := range endpoints {
		resp, err := http.Get(url)
		if err != nil {
			return "", fmt.Errorf("error fetching data from %s: %w", url, err)
		}

		// Read the response body.
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close() // Ensure we close the body immediately.
		if err != nil {
			return "", fmt.Errorf("error reading data from %s: %w", url, err)
		}

		// Write the body into the hasher.
		_, err = hasher.Write(body)
		if err != nil {
			return "", fmt.Errorf("error writing data to hash: %w", err)
		}
	}

	// Compute the final hash and convert it to a hex string.
	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	return hashString, nil
}

// SetTrayIcon initializes the system tray icon and menu, with retry logic on Windows.
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

// AppendSessionRow appends a formatted session row to the table if date and time are provided.
func AppendSessionRow(rows [][]string, label, date, time string, use24h bool) [][]string {
	if date != "" && time != "" {
		d, t := Localize(date, time, use24h)
		rows = append(rows, []string{label, d, t})
	}
	return rows
}

// IsSessionInProgress returns true if the given session time is currently active.
func IsSessionInProgress(dateStr, timeStr string) bool {
	lower := strings.ToLower(dateStr + " " + timeStr)
	if config.Get().DebugMode && strings.Contains(lower, "race") {
		return true
	}

	start, err := time.Parse("2006-01-02 15:04", dateStr+" "+timeStr)
	if err != nil {
		return false
	}
	now := time.Now().UTC()
	duration := 1 * time.Hour
	if strings.Contains(lower, "race") {
		duration = 2 * time.Hour
	}
	return now.After(start) && now.Before(start.Add(duration))
}
