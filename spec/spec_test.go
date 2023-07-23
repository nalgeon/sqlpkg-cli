package spec

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"sqlpkg.org/cli/fileio"
)

func TestPackage_FullName(t *testing.T) {
	p := &Package{Owner: "nalgeon", Name: "example", Version: "0.1.0"}
	if p.FullName() != "nalgeon/example" {
		t.Errorf("FullName: unexpected value %v", p.FullName())
	}
}

func TestPackage_ExpandVars(t *testing.T) {
	t.Run("infer asset path", func(t *testing.T) {
		p := &Package{
			Owner: "nalgeon", Name: "example", Version: "0.1.0",
			Repository: "https://github.com/nalgeon/sqlite-example",
			Assets: Assets{
				Files: map[string]string{"linux-amd64": "example-linux-{version}-x86.zip"},
			},
		}

		p.ExpandVars()
		want := "https://github.com/nalgeon/sqlite-example/releases/download/0.1.0"
		if p.Assets.Path.Value != want {
			t.Errorf("ExpandVars: unexpected Assets.Path = %v", p.Assets.Path)
		}
	})
	t.Run("expand asset path", func(t *testing.T) {
		p := &Package{
			Owner: "nalgeon", Name: "example", Version: "0.1.0",
			Assets: Assets{
				Path: &AssetPath{
					Value:    "https://antonz.org/{owner}/{name}/{version}",
					IsRemote: true,
				},
				Files: map[string]string{"linux-amd64": "example-linux-{version}-x86.zip"},
			},
		}

		p.ExpandVars()
		want := "https://antonz.org/nalgeon/example/0.1.0"
		if p.Assets.Path.Value != want {
			t.Errorf("ExpandVars: unexpected Assets.Path = %v", p.Assets.Path)
		}
	})
	t.Run("expand asset files", func(t *testing.T) {
		p := &Package{
			Owner: "nalgeon", Name: "example", Version: "0.1.0",
			Assets: Assets{
				Files: map[string]string{
					"linux-amd64":   "example-linux-{version}-x86.zip",
					"windows-amd64": "example-win-{version}-x64.zip",
				},
			},
		}

		p.ExpandVars()
		if p.Assets.Files["linux-amd64"] != "example-linux-0.1.0-x86.zip" {
			t.Errorf("ExpandVars: unexpected Assets.Files = %v", p.Assets.Files["linux-amd64"])
		}
		if p.Assets.Files["windows-amd64"] != "example-win-0.1.0-x64.zip" {
			t.Errorf("ExpandVars: unexpected Assets.Files = %v", p.Assets.Files["windows-amd64"])
		}
	})
	t.Run("latest version", func(t *testing.T) {
		p := &Package{
			Owner: "nalgeon", Name: "example", Version: "latest",
			Assets: Assets{
				Path: &AssetPath{
					Value:    "https://antonz.org/{version}",
					IsRemote: true,
				},
				Files: map[string]string{"linux-amd64": "example-linux-{version}-x86.zip"},
			},
		}

		p.ExpandVars()
		want := "https://antonz.org/{latest}"
		if p.Assets.Path.Value != want {
			t.Errorf("ExpandVars: unexpected Assets.Path = %v", p.Assets.Path)
		}
		if p.Assets.Files["linux-amd64"] != "example-linux-{latest}-x86.zip" {
			t.Errorf("ExpandVars: unexpected Assets.Files = %v", p.Assets.Files["linux-amd64"])
		}
	})
}

func TestPackage_ReplaceLatest(t *testing.T) {
	t.Run("latest", func(t *testing.T) {
		p := &Package{
			Owner: "nalgeon", Name: "example", Version: "latest",
			Assets: Assets{
				Path: &AssetPath{
					Value:    "https://antonz.org/{latest}",
					IsRemote: true,
				},
				Files: map[string]string{
					"linux-amd64":   "example-linux-{latest}-x86.zip",
					"windows-amd64": "example-win-{latest}-x64.zip",
				},
			},
		}

		p.ReplaceLatest("0.1.0")
		if p.Version != "0.1.0" {
			t.Errorf("ReplaceLatest: unexpected Version = %v", p.Version)
		}
		if p.Assets.Path.Value != "https://antonz.org/0.1.0" {
			t.Errorf("ReplaceLatest: unexpected Assets.Path = %v", p.Assets.Path)
		}
		if p.Assets.Files["linux-amd64"] != "example-linux-0.1.0-x86.zip" {
			t.Errorf("ReplaceLatest: unexpected Assets.Files = %v", p.Assets.Files["linux-amd64"])
		}
		if p.Assets.Files["windows-amd64"] != "example-win-0.1.0-x64.zip" {
			t.Errorf("ReplaceLatest: unexpected Assets.Files = %v", p.Assets.Files["windows-amd64"])
		}
	})

	t.Run("specific", func(t *testing.T) {
		p := &Package{
			Owner: "nalgeon", Name: "example", Version: "0.1.0",
			Assets: Assets{
				Path: &AssetPath{
					Value:    "https://antonz.org/0.1.0",
					IsRemote: true,
				},
				Files: map[string]string{
					"linux-amd64":   "example-linux-0.1.0-x86.zip",
					"windows-amd64": "example-win-0.1.0-x64.zip",
				},
			},
		}

		p.ReplaceLatest("0.2.0")
		if p.Version != "0.1.0" {
			t.Errorf("ReplaceLatest: unexpected Version = %v", p.Version)
		}
		if p.Assets.Path.Value != "https://antonz.org/0.1.0" {
			t.Errorf("ReplaceLatest: unexpected Assets.Path = %v", p.Assets.Path)
		}
		if p.Assets.Files["linux-amd64"] != "example-linux-0.1.0-x86.zip" {
			t.Errorf("ReplaceLatest: unexpected Assets.Files = %v", p.Assets.Files["linux-amd64"])
		}
		if p.Assets.Files["windows-amd64"] != "example-win-0.1.0-x64.zip" {
			t.Errorf("ReplaceLatest: unexpected Assets.Files = %v", p.Assets.Files["windows-amd64"])
		}
	})
}

func TestPackage_AssetPath(t *testing.T) {
	t.Run("supported platform", func(t *testing.T) {
		p := &Package{
			Owner: "nalgeon", Name: "example", Version: "0.1.0",
			Assets: Assets{
				Path: &AssetPath{Value: "https://antonz.org", IsRemote: true},
				Files: map[string]string{
					"linux-amd64":   "example-lin-x86.zip",
					"windows-amd64": "example-win-x64.zip",
				},
			},
		}

		path, err := p.AssetPath("linux", "amd64")
		if err != nil {
			t.Fatalf("AssetPath: unexpected error %v", err)
		}
		if path.Value != "https://antonz.org/example-lin-x86.zip" {
			t.Errorf("AssetPath: unexpected value %s", path)
		}
	})
	t.Run("unsupported platform", func(t *testing.T) {
		p := &Package{
			Owner: "nalgeon", Name: "example", Version: "0.1.0",
			Assets: Assets{
				Path: &AssetPath{Value: "https://antonz.org", IsRemote: true},
				Files: map[string]string{
					"linux-amd64":   "example-lin-x86.zip",
					"windows-amd64": "example-win-x64.zip",
				},
			},
		}

		_, err := p.AssetPath("darwin", "amd64")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "platform is not supported") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("missing base path", func(t *testing.T) {
		p := &Package{
			Owner: "nalgeon", Name: "example", Version: "0.1.0",
			Assets: Assets{
				Files: map[string]string{
					"linux-amd64":   "example-lin-x86.zip",
					"windows-amd64": "example-win-x64.zip",
				},
			},
		}

		_, err := p.AssetPath("linux", "amd64")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "asset path is not set") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestPackage_Save(t *testing.T) {
	p := &Package{
		Owner: "nalgeon", Name: "example", Version: "0.1.0",
		Homepage:    "https://antonz.org",
		Repository:  "https://github.com/nalgeon/sqlite-example",
		Specfile:    "https://github.com/nalgeon/sqlite-example/raw/main/sqlpkg.json",
		Authors:     []string{"Anton Zhiyanov", "Mystic Stranger"},
		License:     "MIT",
		Description: "Just an example.",
		Keywords:    []string{"sqlite-example", "extenstion"},
		Symbols:     []string{"get", "set", "del"},
		Assets: Assets{
			Path: &AssetPath{Value: "https://antonz.org", IsRemote: true},
			Files: map[string]string{
				"linux-amd64":   "example-lin-x86.zip",
				"windows-amd64": "example-win-x64.zip",
			},
		},
	}
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatalf("MkdirTemp: unexpected error %v", err)
	}
	defer os.RemoveAll(dir)

	err = p.Save(dir)
	if err != nil {
		t.Fatalf("Save: unexpected error %v", err)
	}

	path := filepath.Join(dir, FileName)
	got, err := fileio.ReadJSON[Package](path)
	if err != nil {
		t.Fatalf("ReadJSON: unexpected error %v", err)
	}
	if !reflect.DeepEqual(got, p) {
		t.Errorf("Save: unexpected decoded spec: %+v", got)
	}
}

func TestPackage_Dir(t *testing.T) {
	got := Dir("testdata", "nalgeon", "example")
	if got != filepath.Join("testdata", DirName, "nalgeon", "example") {
		t.Errorf("Dir: unexpected value %v", got)
	}
}

func TestPackage_Path(t *testing.T) {
	got := Path("testdata", "nalgeon", "example")
	if got != filepath.Join("testdata", DirName, "nalgeon", "example", FileName) {
		t.Errorf("Path: unexpected value %v", got)
	}
}

func TestPackage_inferAssetUrl(t *testing.T) {
	tests := []struct{ name, url, want string }{
		{
			"github",
			"https://github.com/nalgeon/sqlite-example",
			"{repository}/releases/download/{version}",
		},
		{
			"custom",
			"https://antonz.org/sqlite-example",
			"",
		},
		{
			"local",
			"sqlite-example",
			"",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := inferAssetUrl(test.url)
			if got != test.want {
				t.Errorf("inferAssetUrl: unexpected value %v", got)
			}
		})
	}
}
