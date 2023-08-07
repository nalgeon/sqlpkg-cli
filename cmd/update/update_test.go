package update

import (
	"path/filepath"
	"strings"
	"testing"

	"sqlpkg.org/cli/cmd"
	"sqlpkg.org/cli/fileio"
	"sqlpkg.org/cli/lockfile"
	"sqlpkg.org/cli/logx"
	"sqlpkg.org/cli/spec"
)

func TestUpdate(t *testing.T) {
	repoDir, lockPath := cmd.SetupTestRepo(t)
	defer cmd.TeardownTestRepo(t)
	cmd.CopyTestRepo(t, "success")

	mem := logx.Mock()

	args := []string{"nalgeon/example"}
	err := Update(args)
	if err != nil {
		t.Fatalf("update error: %v", err)
	}

	validateLog(t, mem)
	validatePackage(t, repoDir, lockPath, "nalgeon", "example")
}

func TestUpdateAll(t *testing.T) {
	repoDir, lockPath := cmd.SetupTestRepo(t)
	defer cmd.TeardownTestRepo(t)
	cmd.CopyTestRepo(t, "success")

	mem := logx.Mock()

	args := []string{}
	err := UpdateAll(args)
	if err != nil {
		t.Fatalf("update error: %v", err)
	}

	validateLog(t, mem)
	mem.MustHave(t, "updated 1 packages")
	validatePackage(t, repoDir, lockPath, "nalgeon", "example")
}

func TestLatest(t *testing.T) {
	cmd.SetupTestRepo(t)
	defer cmd.TeardownTestRepo(t)
	cmd.CopyTestRepo(t, "latest")

	mem := logx.Mock()

	args := []string{"nalgeon/example"}
	err := Update(args)
	if err != nil {
		t.Fatalf("update error: %v", err)
	}

	mem.Print()
	mem.MustHave(t, "already at the latest version")

	pkg, err := spec.ReadLocal(spec.Path(cmd.WorkDir, "nalgeon", "example"))
	if err != nil {
		t.Fatalf("read pkg error: %v", err)
	}
	if pkg.Version != "0.1.0" {
		t.Fatalf("unexpected version: %v", pkg.Version)
	}
}

func TestNoVersion(t *testing.T) {
	cmd.SetupTestRepo(t)
	defer cmd.TeardownTestRepo(t)
	cmd.CopyTestRepo(t, "version")

	mem := logx.Mock()

	args := []string{"nalgeon/example"}
	err := Update(args)
	if err != nil {
		t.Fatalf("update error: %v", err)
	}

	mem.Print()
	mem.MustHave(t, "added package to the lockfile")
	mem.MustHave(t, "updated package nalgeon/example")
}

func TestInvalidChecksum(t *testing.T) {
	cmd.SetupTestRepo(t)
	defer cmd.TeardownTestRepo(t)
	cmd.CopyTestRepo(t, "checksum")

	args := []string{"nalgeon/example"}
	err := Update(args)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "asset checksum is invalid") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func validateLog(t *testing.T, mem *logx.Memory) {
	mem.Print()
	mem.MustHave(t, "read existing lockfile")
	mem.MustHave(t, "found local spec")
	mem.MustHave(t, "read package nalgeon/example, version = 0.1.0")
	mem.MustHave(t, "updating nalgeon/example")
	mem.MustHave(t, "local package version = 0.1.0")
	mem.MustHave(t, "read 4 checksums")
	mem.MustHave(t, "checking remote asset")
	mem.MustHave(t, "downloaded example-")
	mem.MustHave(t, "asset checksum is valid")
	mem.MustHave(t, "unpacked 1 files from example-")
	mem.MustHave(t, "added package to the lockfile")
	mem.MustHave(t, "updated package nalgeon/example to 0.2.0")
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
