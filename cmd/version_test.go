package cmd

import (
	"testing"

	"sqlpkg.org/cli/httpx"
	"sqlpkg.org/cli/spec"
)

func TestResolveVersion(t *testing.T) {
	t.Run("specific", func(t *testing.T) {
		pkg := &spec.Package{
			Owner: "nalgeon", Name: "example", Version: "0.1.0",
			Repository: "https://github.com/nalgeon/example",
			Assets: spec.Assets{
				Path: &spec.AssetPath{
					Value:    "https://github.com/nalgeon/example/releases/download/0.1.0",
					IsRemote: true,
				},
				Files: map[string]string{
					"darwin-arm64": "example-0.1.0-darwin.zip",
					"linux-amd64":  "example-0.1.0-linux.zip",
				},
			},
		}

		err := ResolveVersion(pkg)
		if err != nil {
			t.Fatalf("ResolveVersion: unexpected error %v", err)
		}
		if pkg.Version != "0.1.0" {
			t.Errorf("ResolveVersion: unexpected Version %v", pkg.Version)
		}
		if pkg.Assets.Path.Value != "https://github.com/nalgeon/example/releases/download/0.1.0" {
			t.Errorf("ResolveVersion: unexpected Assets.Path %v", pkg.Assets.Path.Value)
		}
		if pkg.Assets.Files["darwin-arm64"] != "example-0.1.0-darwin.zip" ||
			pkg.Assets.Files["linux-amd64"] != "example-0.1.0-linux.zip" {
			t.Errorf("ResolveVersion: unexpected Assets.Files %v", pkg.Assets.Files)
		}
	})
	t.Run("latest", func(t *testing.T) {
		httpx.Mock("github")
		pkg := &spec.Package{
			Owner: "nalgeon", Name: "example", Version: "latest",
			Repository: "https://github.com/nalgeon/example",
			Assets: spec.Assets{
				Path: &spec.AssetPath{
					Value:    "https://github.com/nalgeon/example/releases/download/{latest}",
					IsRemote: true,
				},
				Files: map[string]string{
					"darwin-arm64": "example-{latest}-darwin.zip",
					"linux-amd64":  "example-{latest}-linux.zip",
				},
			},
		}

		err := ResolveVersion(pkg)
		if err != nil {
			t.Fatalf("ResolveVersion: unexpected error %v", err)
		}
		if pkg.Version != "0.2.0" {
			t.Errorf("ResolveVersion: unexpected Version %v", pkg.Version)
		}
		if pkg.Assets.Path.Value != "https://github.com/nalgeon/example/releases/download/0.2.0" {
			t.Errorf("ResolveVersion: unexpected Assets.Path %v", pkg.Assets.Path.Value)
		}
		if pkg.Assets.Files["darwin-arm64"] != "example-0.2.0-darwin.zip" ||
			pkg.Assets.Files["linux-amd64"] != "example-0.2.0-linux.zip" {
			t.Errorf("ResolveVersion: unexpected Assets.Files %v", pkg.Assets.Files)
		}
	})
}

func TestHasNewVersion(t *testing.T) {
	t.Run("yes", func(t *testing.T) {
		SetupTestRepo(t)
		defer TeardownTestRepo(t)
		CopyTestRepo(t)

		pkg, err := ReadSpec("./testdata/sqlpkg.json")
		if err != nil {
			t.Fatalf("ReadSpec: %v", err)
		}

		has := HasNewVersion(pkg)
		if !has {
			t.Errorf("HasNewVersion: expected true, got false")
		}
	})
	t.Run("no", func(t *testing.T) {
		SetupTestRepo(t)
		defer TeardownTestRepo(t)
		CopyTestRepo(t)

		pkg, err := ReadSpec("./testdata/sqlpkg.json")
		if err != nil {
			t.Fatalf("ReadSpec: %v", err)
		}
		pkg.Version = "0.1.0"

		has := HasNewVersion(pkg)
		if has {
			t.Errorf("HasNewVersion: expected false, got true")
		}
	})
	t.Run("not versioned", func(t *testing.T) {
		SetupTestRepo(t)
		defer TeardownTestRepo(t)
		CopyTestRepo(t)

		pkg, err := ReadSpec("./testdata/.sqlpkg/sqlite/stmt/sqlpkg.json")
		if err != nil {
			t.Fatalf("ReadSpec: %v", err)
		}

		has := HasNewVersion(pkg)
		if !has {
			t.Errorf("HasNewVersion: expected true, got false")
		}
	})
	t.Run("not installed", func(t *testing.T) {
		pkg, err := ReadSpec("./testdata/sqlpkg.json")
		if err != nil {
			t.Fatalf("ReadSpec: %v", err)
		}

		has := HasNewVersion(pkg)
		if !has {
			t.Errorf("HasNewVersion: expected true, got false")
		}
	})
}
