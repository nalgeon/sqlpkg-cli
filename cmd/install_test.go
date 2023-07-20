package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"sqlpkg.org/cli/fileio"
	"sqlpkg.org/cli/lockfile"
)

func TestInstall(t *testing.T) {
	workDir = "."
	repoDir, lockPath := setupRepo(t)

	pkgDir := filepath.Join(repoDir, "asg017", "hello")
	args := []string{filepath.Join(workDir, "testdata", "hello.json")}
	IsVerbose = true
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

	assets, _ := filepath.Glob(filepath.Join(pkgDir, "hello0.*"))
	if len(assets) == 0 {
		t.Fatal("asset files do not exist")
	}

	lck, err := lockfile.ReadLocal(lockPath)
	if err != nil {
		t.Fatal("failed to read lockfile")
	}
	if !lck.Has("asg017/hello") {
		t.Fatal("installed package not found in the lockfile")
	}

	pkg := lck.Packages["asg017/hello"]
	nAssets := len(pkg.Assets.Files)
	nChecksums := len(pkg.Assets.Checksums)
	if nChecksums != nAssets {
		t.Fatalf("got %d checksums, want %d", nChecksums, nAssets)
	}

	teardownRepo(t, repoDir, lockPath)
}

func setupRepo(t *testing.T) (string, string) {
	repoDir := filepath.Join(workDir, ".sqlpkg")
	err := os.RemoveAll(repoDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lockPath := filepath.Join(workDir, "sqlpkg.lock")
	err = os.RemoveAll(lockPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return repoDir, lockPath
}

func teardownRepo(t *testing.T, repoDir, lockPath string) {
	err := os.RemoveAll(repoDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = os.RemoveAll(lockPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
