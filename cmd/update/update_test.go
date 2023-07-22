package update

import (
	"path/filepath"
	"strings"
	"testing"

	"sqlpkg.org/cli/cmd"
	"sqlpkg.org/cli/fileio"
	"sqlpkg.org/cli/lockfile"
)

func TestUpdate(t *testing.T) {
	cmd.WorkDir = "."
	repoDir, lockPath := cmd.SetupTestRepo(t)
	cmd.CopyTestRepo(t)

	cmd.IsVerbose = true
	args := []string{"nalgeon/example"}

	err := Update(args)
	if err != nil {
		t.Fatalf("update error: %v", err)
	}

	validateUpdatedPackage(t, repoDir, lockPath, "nalgeon", "example")

	cmd.TeardownTestRepo(t, repoDir, lockPath)
}

func TestUpdateAll(t *testing.T) {
	cmd.WorkDir = "."
	repoDir, lockPath := cmd.SetupTestRepo(t)
	cmd.CopyTestRepo(t)

	cmd.IsVerbose = true
	args := []string{}

	err := UpdateAll(args)
	if err != nil {
		t.Fatalf("update error: %v", err)
	}

	validateUpdatedPackage(t, repoDir, lockPath, "nalgeon", "example")

	cmd.TeardownTestRepo(t, repoDir, lockPath)
}

func validateUpdatedPackage(t *testing.T, repoDir, lockPath, owner, name string) {
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
		t.Fatal("updated package not found in the lockfile")
	}

	pkg := lck.Packages["nalgeon/example"]
	assetName := "example-linux-0.2.0-x86.zip"
	if pkg.Assets.Files["linux-amd64"] != assetName {
		t.Fatalf("unexpected linux asset: %s", pkg.Assets.Files["linux-amd64"])
	}
	if !strings.HasPrefix(pkg.Assets.Checksums[assetName], "sha256-00cfa32") {
		t.Fatalf("unexpected linux checksum: %s", pkg.Assets.Checksums[assetName])
	}
}
