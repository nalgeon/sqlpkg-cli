package httpx

import (
	"io"
	"testing"
)

func TestIsURL(t *testing.T) {
	tests := []struct {
		url string
		ok  bool
	}{
		{"https://antonz.org/sqlpkg.json", true},
		{"https://github.com/nalgeon/sqlpkg/raw/main/pkg/sqlite/stmt.json", true},
		{"https://raw.githubusercontent.com/nalgeon/sqlpkg/main/pkg/sqlite/stmt.json", true},
		{"./testdata/sqlpkg.json", false},
		{"/Users/anton/sqlpkg.json", false},
		{"file:///Users/anton/sqlpkg.json", false},
	}
	for _, test := range tests {
		ok := IsURL(test.url)
		if ok != test.ok {
			t.Errorf("IsURL(%s): expected %v, got %v", test.url, test.ok, ok)
		}
	}
}

func TestHostname(t *testing.T) {
	tests := []struct {
		url  string
		host string
	}{
		{"https://antonz.org/sqlpkg.json", "antonz.org"},
		{"https://github.com/nalgeon/sqlpkg/raw/main/pkg/sqlite/stmt.json", "github.com"},
		{"https://raw.githubusercontent.com/nalgeon/sqlpkg/main/pkg/sqlite/stmt.json", "raw.githubusercontent.com"},
		{"./testdata/sqlpkg.json", ""},
		{"/Users/anton/sqlpkg.json", ""},
		{"file:///Users/anton/sqlpkg.json", ""},
	}
	for _, test := range tests {
		host := Hostname(test.url)
		if host != test.host {
			t.Errorf("Hostname(%s): expected %v, got %v", test.url, test.host, host)
		}
	}
}

func TestExists(t *testing.T) {
	srv := MockServer()
	defer srv.Close()

	t.Run("exists", func(t *testing.T) {
		ok := Exists(srv.URL + "/sqlpkg.json")
		if !ok {
			t.Errorf("Exists: unexpected %v", ok)
		}
	})
	t.Run("does not exist", func(t *testing.T) {
		ok := Exists(srv.URL + "/missing.json")
		if ok {
			t.Errorf("Exists: unexpected %v", ok)
		}
	})
}

func TestGetBody(t *testing.T) {
	srv := MockServer()
	defer srv.Close()

	t.Run("success", func(t *testing.T) {
		body, err := GetBody(srv.URL+"/example.txt", "text/plain")
		if err != nil {
			t.Errorf("GetBody: unexpected error %v", err)
		}
		defer body.Close()

		data, err := io.ReadAll(body)
		if err != nil {
			t.Errorf("io.ReadAll: unexpected error %v", err)
		}

		if string(data) != "example.txt" {
			t.Errorf("GetBody: unexpected value %q", string(data))
		}
	})
	t.Run("failure", func(t *testing.T) {
		_, err := GetBody(srv.URL+"/missing.txt", "text/plain")
		if err == nil {
			t.Error("GetBody: expected error, got nil")
		}
	})
}

func TestGetBytes(t *testing.T) {
	srv := MockServer()
	defer srv.Close()

	t.Run("success", func(t *testing.T) {
		data, err := GetBytes(srv.URL + "/example.txt")
		if err != nil {
			t.Errorf("GetBytes: unexpected error %v", err)
		}
		if string(data) != "example.txt" {
			t.Errorf("GetBytes: unexpected value %q", string(data))
		}
	})
	t.Run("failure", func(t *testing.T) {
		_, err := GetBytes(srv.URL + "/missing.txt")
		if err == nil {
			t.Error("GetBytes: expected error, got nil")
		}
	})
}

func TestGetJSON(t *testing.T) {
	srv := MockServer()
	defer srv.Close()

	t.Run("success", func(t *testing.T) {
		type Example struct{ Body string }
		ex, err := GetJSON[Example](srv.URL + "/example.json")
		if err != nil {
			t.Errorf("GetJSON: unexpected error %v", err)
		}
		if ex.Body != "example.txt" {
			t.Errorf("GetJSON: unexpected value %q", ex.Body)
		}
	})
	t.Run("failure", func(t *testing.T) {
		type Example struct{ Body string }
		_, err := GetJSON[Example](srv.URL + "/example.txt")
		if err == nil {
			t.Error("GetJSON: expected error, got nil")
		}
	})
}
