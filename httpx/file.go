package httpx

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

var responders = map[string]func(data []byte) io.Reader{
	".json": respondJSON,
}

// FileClient serves responses from the file system instead of remote calls.
// Should be used for testing purposes only.
type FileClient struct {
	routes map[string]string
}

// NewFileClient creates a new FileClient.
func NewFileClient() *FileClient {
	routes := map[string]string{}
	return &FileClient{routes: routes}
}

// AddRoute defines a new route that maps an URL to a local file path.
func (c *FileClient) AddRoute(url, path string) {
	c.routes[url] = path
}

// Do serves the file according to the request URL.
func (c *FileClient) Do(req *http.Request) (*http.Response, error) {
	filename, ok := c.routes[req.URL.String()]
	if !ok {
		resp := http.Response{
			Status:     http.StatusText(http.StatusNotFound),
			StatusCode: http.StatusNotFound,
		}
		return &resp, nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	responder, ok := responders[path.Ext(filename)]
	if !ok {
		panic("unsupported file type: " + filename)
	}
	rdr := responder(data)
	resp, err := http.ReadResponse(bufio.NewReader(rdr), req)
	if err != nil {
		panic(err)
	}
	return resp, nil
}

func respondJSON(data []byte) io.Reader {
	text := fmt.Sprintf("%s\n\n%s",
		"HTTP/1.1 200 OK\nContent-Type: application/json",
		data,
	)
	return strings.NewReader(text)
}
