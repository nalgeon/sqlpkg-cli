package uninstall

import (
	"path/filepath"
	"testing"

	"sqlpkg.org/cli/cmd"
	"sqlpkg.org/cli/fileio"
	"sqlpkg.org/cli/lockfile"
)

func TestUninstall(t *testing.T) {
	cmd.WorkDir = "."
	repoDir, lockPath := cmd.SetupTestRepo(t)
	cmd.CopyTestRepo(t)

	cmd.IsVerbose = true
	args := []string{"nalgeon/example"}
	err := Uninstall(args)
	if err != nil {
		t.Fatalf("uninstallation error: %v", err)
	}

	pkgDir := filepath.Join(repoDir, "nalgeon", "example")
	if fileio.Exists(pkgDir) {
		t.Fatalf("package dir still exists: %v", pkgDir)
	}

	lck, err := lockfile.ReadLocal(lockPath)
	if err != nil {
		t.Fatal("failed to read lockfile")
	}
	if lck.Has("nalgeon/example") {
		t.Fatal("uninstalled package found in the lockfile")
	}

	cmd.TeardownTestRepo(t, repoDir, lockPath)
}
