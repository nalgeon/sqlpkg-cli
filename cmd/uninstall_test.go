package cmd

import (
	"path/filepath"
	"testing"

	"github.com/nalgeon/sqlpkg-cli/internal/fileio"
)

func TestUninstall(t *testing.T) {
	workDir = "."
	repoDir := setupRepo(t)
	install(t, repoDir)

	IsVerbose = true
	args := []string{"asg017/hello"}
	err := Uninstall(args)
	if err != nil {
		t.Fatalf("uninstallation error: %v", err)
	}

	pkgDir := filepath.Join(repoDir, "asg017", "hello")
	if fileio.Exists(pkgDir) {
		t.Fatalf("package dir still exists: %v", pkgDir)
	}

	teardownRepo(t, repoDir)
}

func install(t *testing.T, repoDir string) {
	args := []string{filepath.Join(workDir, "testdata", "hello.json")}
	err := Install(args)
	if err != nil {
		t.Fatalf("installation error: %v", err)
	}
}
