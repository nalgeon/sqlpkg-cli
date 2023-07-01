package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nalgeon/sqlpkg-cli/internal/fileio"
)

func TestInstall(t *testing.T) {
	workDir = "."
	repoDir := setupRepo(t)

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

	teardownRepo(t, repoDir)
}

func setupRepo(t *testing.T) string {
	repoDir := filepath.Join(workDir, ".sqlpkg")
	err := os.RemoveAll(repoDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return repoDir
}

func teardownRepo(t *testing.T, repoDir string) {
	err := os.RemoveAll(repoDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
