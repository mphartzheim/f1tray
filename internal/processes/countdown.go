package processes

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"github.com/mphartzheim/f1tray/internal/models"
)

func StartCountdown(binding binding.String, state *models.AppState) {
	fmt.Println("Countdown started — finding next session...")
	for {
		fmt.Println("📦 Loaded sessions:", len(state.UpcomingSessions))
		for i, s := range state.UpcomingSessions {
			fmt.Printf("  [%d] %s — %s\n", i, s.Label, s.StartTime.Format(time.RFC3339))
			if s.StartTime.IsZero() {
				fmt.Println("     ⚠️ StartTime is zero (not parsed)")
			}
		}

		next := findNextSession(state)
		if next == nil {
			binding.Set("No upcoming sessions")
			time.Sleep(15 * time.Second)
			continue
		}

		fmt.Printf("Next session found: %s at %s\n", next.Label, next.StartTime.Format(time.RFC3339))
		for {
			now := time.Now()
			if now.After(next.StartTime) {
				break // Find the next session
			}
			remaining := next.StartTime.Sub(now)
			binding.Set(fmt.Sprintf("Next: %s in %02dh %02dm %02ds",
				next.Label,
				int(remaining.Hours()),
				int(remaining.Minutes())%60,
				int(remaining.Seconds())%60,
			))
			time.Sleep(1 * time.Second)
		}
	}
}

func findNextSession(state *models.AppState) *models.SessionInfo {
	now := time.Now()
	fmt.Println("🕒 Current time:", now.Format(time.RFC3339))

	for _, session := range state.UpcomingSessions {
		fmt.Printf("🔎 Checking session: %s\n", session.Label)
		fmt.Println("    StartTime:", session.StartTime.Format(time.RFC3339))
		fmt.Println("    After now? ", session.StartTime.After(now))

		if session.StartTime.After(now) {
			fmt.Println("✅ Found next session:", session.Label)
			return &models.SessionInfo{
				Label:     session.Label,
				StartTime: session.StartTime,
			}
		}
	}

	fmt.Println("❌ No valid upcoming session found.")
	return nil
}
