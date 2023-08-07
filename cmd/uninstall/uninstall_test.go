package uninstall

import (
	"path/filepath"
	"strings"
	"testing"

	"sqlpkg.org/cli/cmd"
	"sqlpkg.org/cli/fileio"
	"sqlpkg.org/cli/lockfile"
	"sqlpkg.org/cli/logx"
)

func TestUninstall(t *testing.T) {
	repoDir, lockPath := cmd.SetupTestRepo(t)
	defer cmd.TeardownTestRepo(t)
	cmd.CopyTestRepo(t, "")

	memory := logx.Mock()

	args := []string{"nalgeon/example"}
	err := Uninstall(args)
	if err != nil {
		t.Fatalf("uninstallation error: %v", err)
	}

	validateLog(t, memory)
	validatePackage(t, repoDir, lockPath, "nalgeon", "example")
}

func TestUnknown(t *testing.T) {
	cmd.SetupTestRepo(t)
	defer cmd.TeardownTestRepo(t)
	cmd.CopyTestRepo(t, "")
	logx.Mock()

	args := []string{"sqlite/unknown"}
	err := Uninstall(args)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "package is not installed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func validateLog(t *testing.T, mem *logx.Memory) {
	mem.Print()
	mem.MustHave(t, "uninstalling nalgeon/example")
	mem.MustHave(t, "deleting dir: .sqlpkg/nalgeon/example")
	mem.MustHave(t, "deleted package dir")
	mem.MustHave(t, "read existing lockfile")
	mem.MustHave(t, "removed package from the lockfile")
	mem.MustHave(t, "uninstalled package nalgeon/example")
}

func validatePackage(t *testing.T, repoDir, lockPath, owner, name string) {
	pkgDir := filepath.Join(repoDir, owner, name)

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
}
