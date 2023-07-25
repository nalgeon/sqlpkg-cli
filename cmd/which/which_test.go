package which

import (
	"strings"
	"testing"

	"sqlpkg.org/cli/cmd"
)

func TestExact(t *testing.T) {
	cmd.WorkDir = "."
	repoDir, lockPath := cmd.SetupTestRepo(t)
	cmd.CopyTestRepo(t, "")

	t.Run("exact", func(t *testing.T) {
		mem := cmd.SetupTestLogger()
		args := []string{"nalgeon/example"}
		err := Which(args)
		if err != nil {
			t.Fatalf("which error: %v", err)
		}

		mem.Print()
		mem.MustHave(t, ".sqlpkg/nalgeon/example/example")
		mem.MustNotHave(t, "exact match not found")
	})
	t.Run("version", func(t *testing.T) {
		mem := cmd.SetupTestLogger()
		args := []string{"nalgeon/version"}
		err := Which(args)
		if err != nil {
			t.Fatalf("which error: %v", err)
		}

		mem.Print()
		mem.MustHave(t, ".sqlpkg/nalgeon/version/version0")
		mem.MustNotHave(t, "exact match not found")
	})
	t.Run("prefix", func(t *testing.T) {
		mem := cmd.SetupTestLogger()
		args := []string{"nalgeon/prefix"}
		err := Which(args)
		if err != nil {
			t.Fatalf("which error: %v", err)
		}

		mem.Print()
		mem.MustHave(t, ".sqlpkg/nalgeon/prefix/libprefix")
		mem.MustNotHave(t, "exact match not found")
	})

	cmd.TeardownTestRepo(t, repoDir, lockPath)
}

func TestPossible(t *testing.T) {
	cmd.WorkDir = "."
	repoDir, lockPath := cmd.SetupTestRepo(t)
	cmd.CopyTestRepo(t, "")
	mem := cmd.SetupTestLogger()

	args := []string{"sqlite/stmt"}
	err := Which(args)
	if err != nil {
		t.Fatalf("which error: %v", err)
	}

	mem.Print()
	mem.MustHave(t, "exact match not found")
	mem.MustHave(t, ".sqlpkg/sqlite/stmt/stmtvtab")

	cmd.TeardownTestRepo(t, repoDir, lockPath)
}

func TestNotFound(t *testing.T) {
	cmd.WorkDir = "."
	repoDir, lockPath := cmd.SetupTestRepo(t)
	cmd.CopyTestRepo(t, "")
	cmd.SetupTestLogger()

	args := []string{"sqlite/broken"}
	err := Which(args)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "extension file is not found") {
		t.Fatalf("unexpected error: %v", err)
	}

	cmd.TeardownTestRepo(t, repoDir, lockPath)
}

func TestUnknown(t *testing.T) {
	cmd.WorkDir = "."
	repoDir, lockPath := cmd.SetupTestRepo(t)
	cmd.SetupTestLogger()

	args := []string{"sqlite/unknown"}
	err := Which(args)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "package is not installed") {
		t.Fatalf("unexpected error: %v", err)
	}

	cmd.TeardownTestRepo(t, repoDir, lockPath)
}
