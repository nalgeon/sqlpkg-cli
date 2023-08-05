package spec

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"sqlpkg.org/cli/httpx"
)

func TestRead(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		path := filepath.Join("testdata", "sqlpkg.json")
		got, err := Read(path)
		if err != nil {
			t.Fatalf("Read: unexpected error %v", err)
		}

		want := &Package{
			Owner: "nalgeon", Name: "example", Version: "0.1.0",
			Homepage:    "https://github.com/nalgeon/sqlite-example/blob/main/README.md",
			Repository:  "https://github.com/nalgeon/sqlite-example",
			Specfile:    "testdata/sqlpkg.json",
			Authors:     []string{"Anton Zhiyanov"},
			License:     "MIT",
			Description: "Example extension.",
			Keywords:    []string{"sqlite-example"},
			Assets: Assets{
				Path: &AssetPath{Value: "{repository}/releases/download/{version}", IsRemote: true},
				Files: map[string]string{
					"darwin-amd64":  "example-macos-{version}-x86.zip",
					"darwin-arm64":  "example-macos-{version}-arm64.zip",
					"linux-amd64":   "example-linux-{version}-x86.zip",
					"windows-amd64": "example-win-{version}-x64.zip",
				},
			},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Read: unexpacted package %+v", got)
		}
	})
	t.Run("missing", func(t *testing.T) {
		path := filepath.Join("testdata", "missing", "sqlpkg.json")
		_, err := Read(path)
		if err == nil {
			t.Fatal("Read: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "no such file or directory") {
			t.Errorf("Read: unexpected error %v", err)
		}
	})
}

func TestReadRemote(t *testing.T) {
	httpx.Mock()
	t.Run("valid", func(t *testing.T) {
		url := "https://antonz.org/sqlpkg.json"
		got, err := ReadRemote(url)
		if err != nil {
			t.Fatalf("ReadRemote: unexpected error %v", err)
		}

		want := &Package{
			Owner: "nalgeon", Name: "example", Version: "0.1.0",
			Homepage:    "https://github.com/nalgeon/sqlite-example/blob/main/README.md",
			Repository:  "https://github.com/nalgeon/sqlite-example",
			Authors:     []string{"Anton Zhiyanov"},
			License:     "MIT",
			Description: "Example extension.",
			Keywords:    []string{"sqlite-example"},
			Assets: Assets{
				Path: &AssetPath{Value: "{repository}/releases/download/{version}", IsRemote: true},
				Files: map[string]string{
					"darwin-amd64":  "example-macos-{version}-x86.zip",
					"darwin-arm64":  "example-macos-{version}-arm64.zip",
					"linux-amd64":   "example-linux-{version}-x86.zip",
					"windows-amd64": "example-win-{version}-x64.zip",
				},
			},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Read: unexpacted package %+v", got)
		}
	})
	t.Run("missing", func(t *testing.T) {
		url := "https://github.com/nalgeon/sqlite-example/blob/main/missing.json"
		_, err := ReadRemote(url)
		if err == nil {
			t.Fatal("ReadRemote: expected error, got nil")
		}
		if !strings.Contains(err.Error(), "http status 404") {
			t.Errorf("ReadRemote: unexpected error %v", err)
		}
	})
}

func Test_expandPath(t *testing.T) {
	tests := []struct {
		name, path string
		want       []string
	}{
		{"local", "./testdata/sqlpkg.json", []string{"./testdata/sqlpkg.json"}},
		{
			"github",
			"github.com/nalgeon/example",
			[]string{"https://github.com/nalgeon/example/raw/main/sqlpkg.json"},
		},
		{
			"owner-name",
			"nalgeon/example",
			[]string{
				"nalgeon/example",
				"https://github.com/nalgeon/example/raw/main/sqlpkg.json",
				"https://github.com/nalgeon/sqlpkg/raw/main/pkg/nalgeon/example.json",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := expandPath(test.path)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("expandPath: unexpected value %v", got)
			}
		})
	}
}

func Test_inferReader(t *testing.T) {
	var fnLocal = ReadFunc(ReadLocal)
	var fnRemote = ReadFunc(ReadRemote)
	t.Run("local", func(t *testing.T) {
		got := inferReader("./testdata/sqlpkg.json")
		if &got == &fnLocal {
			t.Errorf("inferReader: unexpected value %#v", got)
		}
	})
	t.Run("remote", func(t *testing.T) {
		got := inferReader("https://antonz.org/sqlpkg.json")
		if &got == &fnRemote {
			t.Errorf("inferReader: unexpected value %#v", got)
		}
	})
}
