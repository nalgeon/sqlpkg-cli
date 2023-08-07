package which

import (
	"strings"
	"testing"

	"sqlpkg.org/cli/cmd"
	"sqlpkg.org/cli/logx"
)

func TestExact(t *testing.T) {
	cmd.SetupTestRepo(t)
	defer cmd.TeardownTestRepo(t)
	cmd.CopyTestRepo(t, "")

	t.Run("exact", func(t *testing.T) {
		mem := logx.Mock()
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
		mem := logx.Mock()
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
		mem := logx.Mock()
		args := []string{"nalgeon/prefix"}
		err := Which(args)
		if err != nil {
			t.Fatalf("which error: %v", err)
		}

		mem.Print()
		mem.MustHave(t, ".sqlpkg/nalgeon/prefix/libprefix")
		mem.MustNotHave(t, "exact match not found")
	})
}

func TestPossible(t *testing.T) {
	cmd.SetupTestRepo(t)
	defer cmd.TeardownTestRepo(t)
	cmd.CopyTestRepo(t, "")
	mem := logx.Mock()

	args := []string{"sqlite/stmt"}
	err := Which(args)
	if err != nil {
		t.Fatalf("which error: %v", err)
	}

	mem.Print()
	mem.MustHave(t, "exact match not found")
	mem.MustHave(t, ".sqlpkg/sqlite/stmt/stmtvtab")
}

func TestNotFound(t *testing.T) {
	cmd.SetupTestRepo(t)
	defer cmd.TeardownTestRepo(t)
	cmd.CopyTestRepo(t, "")
	logx.Mock()

	args := []string{"sqlite/broken"}
	err := Which(args)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "extension file is not found") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnknown(t *testing.T) {
	cmd.SetupTestRepo(t)
	defer cmd.TeardownTestRepo(t)
	logx.Mock()

	args := []string{"sqlite/unknown"}
	err := Which(args)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "package is not installed") {
		t.Fatalf("unexpected error: %v", err)
	}
}
