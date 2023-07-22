package install

import (
	"path/filepath"
	"testing"

	"sqlpkg.org/cli/cmd"
	"sqlpkg.org/cli/fileio"
	"sqlpkg.org/cli/lockfile"
	"sqlpkg.org/cli/logx"
)

func TestInstall(t *testing.T) {
	cmd.WorkDir = "."
	repoDir, lockPath := cmd.SetupTestRepo(t)
	memory := cmd.SetupTestLogger()

	args := []string{filepath.Join(cmd.WorkDir, "testdata", "sqlpkg.json")}
	err := Install(args)
	if err != nil {
		t.Fatalf("installation error: %v", err)
	}

	validateLog(t, memory)
	validatePackage(t, repoDir, lockPath, "nalgeon", "example")

	cmd.TeardownTestRepo(t, repoDir, lockPath)
}

func validateLog(t *testing.T, mem *logx.Memory) {
	mem.Print()
	mem.MustHave(t, "installing testdata/sqlpkg.json")
	mem.MustHave(t, "read package nalgeon/example, version = 0.1.0")
	mem.MustHave(t, "read 4 checksums")
	mem.MustHave(t, "downloaded example-")
	mem.MustHave(t, "asset checksum is valid")
	mem.MustHave(t, "unpacked 1 files")
	mem.MustHave(t, "created new lockfile")
	mem.MustHave(t, "added package to the lockfile")
	mem.MustHave(t, "installed package nalgeon/example")
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
	nAssets := len(pkg.Assets.Files)
	nChecksums := len(pkg.Assets.Checksums)
	if nChecksums != nAssets {
		t.Fatalf("got %d checksums, want %d", nChecksums, nAssets)
	}
	if pkg.Assets.Files["linux-amd64"] != "example-linux-0.1.0-x86.zip" {
		t.Fatalf("unexpected linux asset: %s", pkg.Assets.Files["linux-amd64"])
	}
}
