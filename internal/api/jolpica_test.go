package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetNextRace(t *testing.T) {
	// Create a mock HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"MRData": map[string]interface{}{
				"RaceTable": map[string]interface{}{
					"Races": []map[string]interface{}{
						{"round": "3"},
					},
				},
			},
		})
	}))
	defer ts.Close()

	// Temporarily override baseURL for testing
	originalBaseURL := baseURL
	baseURL = ts.URL
	defer func() { baseURL = originalBaseURL }()

	round, err := GetNextRace()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if round != "3" {
		t.Errorf("Expected round '3', got '%s'", round)
	}
}
