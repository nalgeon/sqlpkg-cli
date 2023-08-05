// Package httpx provides high-level HTTP operations.
package httpx

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

var client = Client(&http.Client{Timeout: 3 * time.Second})

// Client is something that can send HTTP requests.
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

// IsURL checks if the path is an url.
func IsURL(path string) bool {
	u, err := url.Parse(path)
	if err != nil {
		return false
	}
	if u.Scheme == "" {
		return false
	}
	return true
}

// Hostname returns the domain part of the url
// or an empty string if the url is invalid.
func Hostname(rawUrl string) string {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return ""
	}
	return u.Hostname()
}

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

// GetJSON issues a GET request and decodes the response as JSON.
func GetJSON[T any](url string) (*T, error) {
	body, err := GetBody(url, "application/json")
	if err != nil {
		return nil, err
	}
	defer body.Close()

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var val T
	err = json.Unmarshal(data, &val)
	if err != nil {
		return nil, err
	}

	return &val, nil
}

// GetBytes issues a GET request and decodes the response as bytes.
func GetBytes(url string) ([]byte, error) {
	body, err := GetBody(url, "*/*")
	if err != nil {
		return nil, err
	}
	defer body.Close()
	return io.ReadAll(body)
}
