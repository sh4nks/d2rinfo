package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	UserAgentHeader = "D2RLoader (https://github.com/sh4nks/d2rloader)"
)

// FetchJSON is a utils function to fetch and unmarshal JSON from a URL with custom headers.
func FetchJSON(ctx context.Context, url string, username string, token string) (map[string]any, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set custom headers
	req.Header.Set("User-Agent", UserAgentHeader)
	req.Header.Set("x-emu-username", username)
	req.Header.Set("x-emu-token", token)

	client := &http.Client{Timeout: 10 * time.Second} // Add a timeout for external requests
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request to %s: %w", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK HTTP status from %s: %s", url, res.Status)
	}

	var data map[string]any
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON from %s: %w", url, err)
	}

	return data, nil
}
