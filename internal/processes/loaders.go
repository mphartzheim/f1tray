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

// UpdateTabs is a callback function that should be set by the main package.
// When ReloadOtherTabs is called, UpdateTabs will update the tab container's
// "Race Results", "Qualifying", and "Sprint" tabs with the new content.
var UpdateTabs func(resultsContent, qualifyingContent, sprintContent fyne.CanvasObject)

// ReloadOtherTabs updates the Race Results, Qualifying, and Sprint tabs with new content.
// It calls the UpdateTabs callback if it has been set.
func ReloadOtherTabs(resultsContent, qualifyingContent, sprintContent fyne.CanvasObject) {
	if UpdateTabs != nil {
		UpdateTabs(resultsContent, qualifyingContent, sprintContent)
	}
}
