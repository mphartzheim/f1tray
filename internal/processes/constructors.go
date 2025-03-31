package processes

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mphartzheim/f1tray/internal/config"
)

const baseConstructorsURL = "https://api.jolpi.ca/ergast/f1/constructors.json"

// ConstructorsResponse represents the expected JSON structure from the API.
type ConstructorsResponse struct {
	MRData struct {
		Limit            int `json:"limit,string"`
		Offset           int `json:"offset,string"`
		Total            int `json:"total,string"`
		ConstructorTable struct {
			Constructors []json.RawMessage `json:"Constructors"`
		} `json:"ConstructorTable"`
	} `json:"MRData"`
}

// LoadOrUpdateConstructorsJSON ensures a local copy of constructors.json is available.
// It downloads the file only if it is not already present, otherwise it returns the existing file.
func LoadOrUpdateConstructorsJSON() (string, error) {
	// Determine config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config dir: %v", err)
	}

	cachePath := filepath.Join(configDir, "f1tray", "constructors.json")

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %v", err)
	}

	// Check if the file already exists; if so, skip the download.
	if _, err := os.Stat(cachePath); err == nil {
		return cachePath, nil
	}

	// File does not exist, so download and merge the JSON data.
	limit := 100
	offset := 0
	var allConstructors []json.RawMessage
	total := 0

	for {
		// Construct the URL with limit and offset
		url := fmt.Sprintf("%s?limit=%d&offset=%d", baseConstructorsURL, limit, offset)
		resp, err := http.Get(url)
		if err != nil {
			return "", fmt.Errorf("failed to fetch constructor data: %v", err)
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return "", fmt.Errorf("failed to read constructor response: %v", err)
		}

		// Unmarshal the current page
		var cr ConstructorsResponse
		if err := json.Unmarshal(body, &cr); err != nil {
			return "", fmt.Errorf("failed to unmarshal constructor response: %v", err)
		}

		// Append this page's constructors to our accumulator
		allConstructors = append(allConstructors, cr.MRData.ConstructorTable.Constructors...)
		total = cr.MRData.Total

		// Increment offset; if we've fetched all records, break out of the loop
		offset += limit
		if offset >= total {
			break
		}
	}

	// Build a merged response using the accumulated data.
	mergedResponse := ConstructorsResponse{}
	mergedResponse.MRData.Limit = limit
	mergedResponse.MRData.Offset = 0
	mergedResponse.MRData.Total = total
	mergedResponse.MRData.ConstructorTable.Constructors = allConstructors

	mergedJSON, err := json.Marshal(mergedResponse)
	if err != nil {
		return "", fmt.Errorf("failed to marshal merged constructors JSON: %v", err)
	}

	// Save merged JSON to file.
	if err := os.WriteFile(cachePath, mergedJSON, 0644); err != nil {
		return "", fmt.Errorf("failed to write constructors.json: %v", err)
	}

	// Update config with new hash.
	newHash := hashBytes(mergedJSON)
	prefs := config.Get()
	prefs.LastConstructorHash = newHash
	if err := config.SaveConfig(*prefs); err != nil {
		return "", fmt.Errorf("failed to save updated config: %v", err)
	}

	return cachePath, nil
}

func hashBytes(data []byte) string {
	sum := sha256.Sum256(data)
	return fmt.Sprintf("%x", sum)
}
