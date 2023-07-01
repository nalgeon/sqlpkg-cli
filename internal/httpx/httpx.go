// Package httpx provides high-level HTTP operations.
package httpx

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var client = http.Client{Timeout: 3 * time.Second}

// Exists checks if the specified url exists.
func Exists(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		return false
	}
	return resp.StatusCode == http.StatusOK
}

// GetBody issues a GET request with an Accept header and returns the response body.
func GetBody(url string, accept string) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", accept)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got http status %d", resp.StatusCode)
	}

	return resp.Body, nil
}

// GetJSON issues a GET request and decodes the response as JSON into the val.
func GetJSON(url string, val any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got http status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, val)
	if err != nil {
		return err
	}

	return nil
}
