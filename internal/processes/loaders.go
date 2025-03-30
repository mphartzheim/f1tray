package processes

import (
	"fmt"
	"io"
	"net/http"

	"fyne.io/fyne/v2"
)

// FetchData retrieves the response body from a URL.
func FetchData(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("request creation failed: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read error: %w", err)
	}

	return body, nil
}

// UpdateTabs updates the Race Results, Qualifying, and Sprint tabs via callback from main.
var UpdateTabs func(resultsContent, qualifyingContent, sprintContent fyne.CanvasObject)

// UpdateStandingsTabs allows UI components to update the inner Drivers and Constructors tab content.
var UpdateStandingsTabs func(driversContent, constructorsContent fyne.CanvasObject)

// ReloadOtherTabs refreshes the Race Results, Qualifying, and Sprint tabs using UpdateTabs.
func ReloadOtherTabs(resultsContent, qualifyingContent, sprintContent fyne.CanvasObject) {
	if UpdateTabs != nil {
		UpdateTabs(resultsContent, qualifyingContent, sprintContent)
	}
}
