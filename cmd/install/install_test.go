package install

import (
	"path/filepath"
	"strings"
	"testing"

	"sqlpkg.org/cli/cmd"
	"sqlpkg.org/cli/fileio"
	"sqlpkg.org/cli/lockfile"
)

func TestFull(t *testing.T) {
	cmd.WorkDir = "."
	repoDir, lockPath := cmd.SetupTestRepo(t)
	mem := cmd.SetupTestLogger()

	args := []string{filepath.Join(cmd.WorkDir, "testdata", "full", "sqlpkg.json")}
	err := Install(args)
	if err != nil {
		t.Fatalf("installation error: %v", err)
	}

	mem.Print()
	mem.MustHave(t, "installing testdata/full/sqlpkg.json")
	mem.MustHave(t, "read package nalgeon/example, version = 0.1.0")
	mem.MustHave(t, "read 4 checksums")
	mem.MustHave(t, "downloaded example-")
	mem.MustHave(t, "asset checksum is valid")
	mem.MustHave(t, "unpacked 1 files")
	mem.MustHave(t, "created new lockfile")
	mem.MustHave(t, "added package to the lockfile")
	mem.MustHave(t, "installed package nalgeon/example")

	validatePackage(t, repoDir, lockPath, "nalgeon", "example")

	cmd.TeardownTestRepo(t, repoDir, lockPath)
}

func TestLockfile(t *testing.T) {
	cmd.WorkDir = "."
	repoDir, lockPath := cmd.SetupTestRepo(t)
	cmd.CopyTestRepo(t, "lockfile")
	mem := cmd.SetupTestLogger()

	args := []string{}
	err := InstallAll(args)
	if err != nil {
		t.Fatalf("installation error: %v", err)
	}

	mem.Print()
	mem.MustHave(t, "loaded the lockfile with 1 packages")
	mem.MustHave(t, "read package nalgeon/example, version = 0.2.0")
	mem.MustHave(t, "locked version = 0.1.0")
	mem.MustHave(t, "installed package nalgeon/example")

	validatePackage(t, repoDir, lockPath, "nalgeon", "example")

	cmd.TeardownTestRepo(t, repoDir, lockPath)
}

func TestMinimal(t *testing.T) {
	cmd.WorkDir = "."
	repoDir, lockPath := cmd.SetupTestRepo(t)
	mem := cmd.SetupTestLogger()

	args := []string{filepath.Join(cmd.WorkDir, "testdata", "minimal", "sqlpkg.json")}
	err := Install(args)
	if err != nil {
		t.Fatalf("installation error: %v", err)
	}

	mem.Print()
	mem.MustHave(t, "read package nalgeon/example")
	mem.MustHave(t, "missing spec checksum file")
	mem.MustHave(t, "downloaded example")
	mem.MustHave(t, "spec is missing asset checksum")
	mem.MustHave(t, "not an archive, skipping unpack")
	mem.MustHave(t, "added package to the lockfile")
	mem.MustHave(t, "installed package nalgeon/example")

	cmd.TeardownTestRepo(t, repoDir, lockPath)
}

func TestAlreadyInstalled(t *testing.T) {
	cmd.WorkDir = "."
	repoDir, lockPath := cmd.SetupTestRepo(t)
	cmd.CopyTestRepo(t, "installed")
	mem := cmd.SetupTestLogger()

	args := []string{filepath.Join(cmd.WorkDir, "testdata", "installed", "sqlpkg.json")}
	err := Install(args)
	if err != nil {
		t.Fatalf("installation error: %v", err)
	}

	mem.Print()
	mem.MustHave(t, "already at the latest version")

	cmd.TeardownTestRepo(t, repoDir, lockPath)
}

func TestInvalidChecksum(t *testing.T) {
	cmd.WorkDir = "."
	repoDir, lockPath := cmd.SetupTestRepo(t)
	cmd.SetupTestLogger()

	args := []string{filepath.Join(cmd.WorkDir, "testdata", "checksum", "sqlpkg.json")}
	err := Install(args)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "asset checksum is invalid") {
		t.Fatalf("unexpected error: %v", err)
	}

	cmd.TeardownTestRepo(t, repoDir, lockPath)
}

func TestUnsupportedPlatform(t *testing.T) {
	cmd.WorkDir = "."
	repoDir, lockPath := cmd.SetupTestRepo(t)
	cmd.SetupTestLogger()

	args := []string{filepath.Join(cmd.WorkDir, "testdata", "unsupported", "sqlpkg.json")}
	err := Install(args)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported platform") {
		t.Fatalf("unexpected error: %v", err)
	}

	cmd.TeardownTestRepo(t, repoDir, lockPath)
}

func TestUnknown(t *testing.T) {
	cmd.WorkDir = "."
	repoDir, lockPath := cmd.SetupTestRepo(t)
	cmd.SetupTestLogger()

	args := []string{"sqlite/unknown"}
	err := Install(args)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to read package spec") {
		t.Fatalf("unexpected error: %v", err)
	}

	cmd.TeardownTestRepo(t, repoDir, lockPath)
}

func validatePackage(t *testing.T, repoDir, lockPath, owner, name string) {
	pkgDir := filepath.Join(repoDir, owner, name)

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
	if !strings.HasPrefix(pkg.Specfile, "./testdata/") || !strings.HasSuffix(pkg.Specfile, "sqlpkg.json") {
		t.Fatalf("unexpected specfile: %v", pkg.Specfile)
	}

	nAssets := len(pkg.Assets.Files)
	nChecksums := len(pkg.Assets.Checksums)
	if nChecksums != nAssets {
		t.Fatalf("got %d checksums, want %d", nChecksums, nAssets)
	}
	if pkg.Assets.Files["linux-amd64"] != "example-linux-0.1.0-x86.zip" {
		t.Fatalf("unexpected linux asset: %s", pkg.Assets.Files["linux-amd64"])
	}
}
