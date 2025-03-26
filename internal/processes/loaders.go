package processes

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

// hashCache stores computed SHA256 hashes for endpoints to detect changes.
var hashCache = make(map[string]string)

// FetchData retrieves the raw response body from the given URL, returning an error if the request fails.
// It does not perform any parsing or hashing checks.
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
