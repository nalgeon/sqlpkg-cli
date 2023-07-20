package cmd

import (
	"path/filepath"
	"testing"

	"sqlpkg.org/cli/fileio"
	"sqlpkg.org/cli/lockfile"
)

func TestUninstall(t *testing.T) {
	workDir = "."
	repoDir, lockPath := setupRepo(t)
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

	lck, err := lockfile.ReadLocal(lockPath)
	if err != nil {
		t.Fatal("failed to read lockfile")
	}
	if lck.Has("asg017/hello") {
		t.Fatal("uninstalled package found in the lockfile")
	}

	teardownRepo(t, repoDir, lockPath)
}

func install(t *testing.T, repoDir string) {
	args := []string{filepath.Join(workDir, "testdata", "hello.json")}
	err := Install(args)
	if err != nil {
		t.Fatalf("installation error: %v", err)
	}
}
