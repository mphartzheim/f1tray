package processes

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	"fyne.io/fyne/v2"
)

// hashCache stores computed SHA256 hashes for endpoints to detect changes.
var hashCache = make(map[string]string)

// FetchData retrieves the response body from a URL and returns whether the content has changed.
func FetchData(url string) ([]byte, bool, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, false, fmt.Errorf("request creation failed: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, false, fmt.Errorf("fetch error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false, fmt.Errorf("read error: %w", err)
	}

	newHashBytes := sha256.Sum256(body)
	newHash := hex.EncodeToString(newHashBytes[:])

	if oldHash, ok := hashCache[url]; ok && oldHash == newHash {
		return body, false, nil // no new data
	}

	hashCache[url] = newHash
	return body, true, nil
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
