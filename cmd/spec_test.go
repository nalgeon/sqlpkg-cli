package cmd

import (
	"testing"

	"sqlpkg.org/cli/httpx"
)

func TestReadSpec(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		pkg, err := ReadSpec("./testdata/sqlpkg.json")
		if err != nil {
			t.Fatalf("ReadSpec: unexpected error %v", err)
		}
		if pkg.FullName() != "nalgeon/example" || pkg.Version != "0.2.0" {
			t.Errorf("ReadSpec: unexpected package %s@%s", pkg.FullName(), pkg.Version)
		}
		if pkg.Assets.Files["darwin-arm64"] != "example-0.2.0-darwin.zip" {
			t.Errorf("ReadSpec: unexpected darwin asset %v", pkg.Assets.Files["darwin-arm64"])
		}
		if pkg.Assets.Files["linux-amd64"] != "example-0.2.0-linux.zip" {
			t.Errorf("ReadSpec: unexpected linux asset %v", pkg.Assets.Files["linux-amd64"])
		}
	})
	t.Run("not found", func(t *testing.T) {
		_, err := ReadSpec("./testdata/missing.json")
		if err == nil {
			t.Fatal("ReadSpec: expected error, got nil")
		}
	})
}

func TestFindSpec(t *testing.T) {
	t.Run("installed", func(t *testing.T) {
		SetupTestRepo(t)
		defer TeardownTestRepo(t)
		CopyTestRepo(t)

		pkg, err := FindSpec("nalgeon/example")
		if err != nil {
			t.Fatalf("FindSpec: unexpected error %v", err)
		}
		if pkg.FullName() != "nalgeon/example" || pkg.Version != "0.1.0" {
			t.Errorf("FindSpec: unexpected package %s@%s", pkg.FullName(), pkg.Version)
		}
	})
	t.Run("remote", func(t *testing.T) {
		httpx.Mock()

		pkg, err := FindSpec("nalgeon/example")
		if err != nil {
			t.Fatalf("FindSpec: unexpected error %v", err)
		}
		if pkg.FullName() != "nalgeon/example" || pkg.Version != "0.2.0" {
			t.Errorf("FindSpec: unexpected package %s@%s", pkg.FullName(), pkg.Version)
		}
	})
}

func TestReadInstalledSpec(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		SetupTestRepo(t)
		defer TeardownTestRepo(t)
		CopyTestRepo(t)

		pkg := ReadInstalledSpec("nalgeon/example")
		if pkg == nil {
			t.Fatal("ReadInstalledSpec: expected package, got nil")
		}
		if pkg.FullName() != "nalgeon/example" || pkg.Version != "0.1.0" {
			t.Errorf("ReadInstalledSpec: unexpected package %s@%s", pkg.FullName(), pkg.Version)
		}
	})
	t.Run("not found", func(t *testing.T) {
		pkg := ReadInstalledSpec("nalgeon/example")
		if pkg != nil {
			t.Fatalf("ReadInstalledSpec: expected nul, got %v", pkg)
		}
	})
}

func TestReadChecksums(t *testing.T) {
	t.Run("exist", func(t *testing.T) {
		pkg, err := ReadSpec("./testdata/checksums/sqlpkg.json")
		if err != nil {
			t.Fatalf("ReadSpec: %v", err)
		}

		err = ReadChecksums(pkg)
		if err != nil {
			t.Fatalf("ReadChecksums: unexpected error %v", err)
		}
		if len(pkg.Assets.Checksums) != 2 {
			t.Fatalf("ReadChecksums: unexpected checksum count %v", len(pkg.Assets.Checksums))
		}
		if pkg.Assets.Checksums["example-darwin.zip"][:19] != "sha256-e3de533fdc23" {
			t.Errorf(
				"ReadChecksums: unexpected darwin checksum %v",
				pkg.Assets.Checksums["example-darwin.zip"],
			)
		}
		if pkg.Assets.Checksums["example-linux.zip"][:19] != "sha256-6bc24897dde2" {
			t.Errorf(
				"ReadChecksums: unexpected linux checksum %v",
				pkg.Assets.Checksums["example-linux.zip"],
			)
		}
	})
	t.Run("not found", func(t *testing.T) {
		pkg, err := ReadSpec("./testdata/sqlpkg.json")
		if err != nil {
			t.Fatalf("ReadSpec: %v", err)
		}

		err = ReadChecksums(pkg)
		if err != nil {
			t.Fatalf("ReadChecksums: unexpected error %v", err)
		}
		if len(pkg.Assets.Checksums) != 0 {
			t.Fatalf("ReadChecksums: unexpected checksum count %v", len(pkg.Assets.Checksums))
		}
	})
}
