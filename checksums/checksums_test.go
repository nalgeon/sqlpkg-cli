package checksums

import (
	"errors"
	"path/filepath"
	"testing"

	"sqlpkg.org/cli/httpx"
)

func TestExists(t *testing.T) {
	t.Run("file", func(t *testing.T) {
		path := filepath.Join("testdata", "checksums.txt")
		ok := Exists(path, false)
		if !ok {
			t.Errorf("Exists: unexpected %v", ok)
		}
	})
	t.Run("http", func(t *testing.T) {
		httpx.Mock()
		path := filepath.Join("https://antonz.org/checksums.txt")
		ok := Exists(path, true)
		if !ok {
			t.Errorf("Exists: unexpected %v", ok)
		}
	})
}

func TestRead(t *testing.T) {
	t.Run("file", func(t *testing.T) {
		path := filepath.Join("testdata", "checksums.txt")
		sums, err := Read(path, false)
		if err != nil {
			t.Fatalf("Read: unexpected error %v", err)
		}
		if len(sums) != 3 {
			t.Fatalf("Read: unexpected length %v", len(sums))
		}
		if sums["example-linux.zip"] != "sha256-6bc24897dde2c7f00cf435055a6853358cb06fcb5a2a789877903ebec0b9298d" {
			t.Fatalf("Read: unexpected example-linux.zip = %v", sums["example-linux.zip"])
		}
		if sums["example-macos.zip"] != "sha256-e3de533fdc23e0d953572c2b544ecc2951b890758af0a00b5a42695ae59ee7ac" {
			t.Fatalf("Read: unexpected example-macos.zip = %v", sums["example-macos.zip"])
		}
		if sums["example-win.zip"] != "sha256-f0d2d705bbe641bf2950a51253820e85de04373b7f428f109f69df1d85fa0654" {
			t.Fatalf("Read: unexpected example-win.zip = %v", sums["example-win.zip"])
		}
	})
	t.Run("http", func(t *testing.T) {
		httpx.Mock()
		path := filepath.Join("https://antonz.org/checksums.txt")
		sums, err := Read(path, true)
		if err != nil {
			t.Fatalf("Read: unexpected error %v", err)
		}
		if len(sums) != 3 {
			t.Fatalf("Read: unexpected length %v", len(sums))
		}
		if sums["example-linux.zip"] != "sha256-6bc24897dde2c7f00cf435055a6853358cb06fcb5a2a789877903ebec0b9298d" {
			t.Fatalf("Read: unexpected example-linux.zip = %v", sums["example-linux.zip"])
		}
		if sums["example-macos.zip"] != "sha256-e3de533fdc23e0d953572c2b544ecc2951b890758af0a00b5a42695ae59ee7ac" {
			t.Fatalf("Read: unexpected example-macos.zip = %v", sums["example-macos.zip"])
		}
		if sums["example-win.zip"] != "sha256-f0d2d705bbe641bf2950a51253820e85de04373b7f428f109f69df1d85fa0654" {
			t.Fatalf("Read: unexpected example-win.zip = %v", sums["example-win.zip"])
		}
	})
	t.Run("invalid file", func(t *testing.T) {
		path := filepath.Join("testdata", "checksums.json")
		_, err := Read(path, false)
		if !errors.Is(err, ErrInvalidFile) {
			t.Fatalf("Read: expected ErrInvalidFile, got %v", err)
		}
	})
	t.Run("invalid sum", func(t *testing.T) {
		path := filepath.Join("testdata", "checksums.sha1")
		_, err := Read(path, false)
		if !errors.Is(err, ErrInvalidSum) {
			t.Fatalf("Read: expected ErrInvalidSum, got %v", err)
		}
	})
}
