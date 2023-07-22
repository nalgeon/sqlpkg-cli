package init

import (
	"strings"
	"testing"

	"sqlpkg.org/cli/cmd"
)

func TestInit(t *testing.T) {
	cmd.WorkDir = "."
	repoDir, lockPath := cmd.SetupTestRepo(t)
	mem := cmd.SetupTestLogger()

	args := []string{}
	err := Init(args)
	if err != nil {
		t.Fatalf("init error: %v", err)
	}

	mem.Print()
	mem.MustHave(t, "created a local repository")

	cmd.TeardownTestRepo(t, repoDir, lockPath)
}

func TestAlreadyExists(t *testing.T) {
	cmd.WorkDir = "."
	repoDir, lockPath := cmd.SetupTestRepo(t)
	cmd.SetupTestLogger()

	args := []string{}
	_ = Init(args)
	err := Init(args)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("unexpected error: %v", err)
	}

	cmd.TeardownTestRepo(t, repoDir, lockPath)
}
