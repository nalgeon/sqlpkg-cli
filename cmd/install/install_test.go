package install

import (
	"path/filepath"
	"testing"

	"sqlpkg.org/cli/cmd"
	"sqlpkg.org/cli/fileio"
	"sqlpkg.org/cli/lockfile"
)

func TestInstall(t *testing.T) {
	cmd.WorkDir = "."
	repoDir, lockPath := cmd.SetupTestRepo(t)

	pkgDir := filepath.Join(repoDir, "nalgeon", "example")
	args := []string{filepath.Join(cmd.WorkDir, "testdata", "sqlpkg.json")}
	cmd.IsVerbose = true

	err := Install(args)
	if err != nil {
		t.Fatalf("installation error: %v", err)
	}

	if !fileio.Exists(pkgDir) {
		t.Fatalf("package dir does not exist: %v", pkgDir)
	}

	specPath := filepath.Join(pkgDir, "sqlpkg.json")
	if !fileio.Exists(specPath) {
		t.Fatalf("spec file does not exist: %v", specPath)
	}

	assets, _ := filepath.Glob(filepath.Join(pkgDir, "example.*"))
	if len(assets) == 0 {
		t.Fatal("asset files do not exist")
	}

	lck, err := lockfile.ReadLocal(lockPath)
	if err != nil {
		t.Fatal("failed to read lockfile")
	}
	if !lck.Has("nalgeon/example") {
		t.Fatal("installed package not found in the lockfile")
	}

	pkg := lck.Packages["nalgeon/example"]
	nAssets := len(pkg.Assets.Files)
	nChecksums := len(pkg.Assets.Checksums)
	if nChecksums != nAssets {
		t.Fatalf("got %d checksums, want %d", nChecksums, nAssets)
	}
	if pkg.Assets.Files["linux-amd64"] != "example-linux-0.1.0-x86.zip" {
		t.Fatalf("unexpected linux asset: %s", pkg.Assets.Files["linux-amd64"])
	}

	cmd.TeardownTestRepo(t, repoDir, lockPath)
}
