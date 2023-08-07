package init

import (
	"strings"
	"testing"

	"sqlpkg.org/cli/cmd"
	"sqlpkg.org/cli/logx"
)

func TestInit(t *testing.T) {
	cmd.SetupTestRepo(t)
	defer cmd.TeardownTestRepo(t)
	mem := logx.Mock()

	args := []string{}
	err := Init(args)
	if err != nil {
		t.Fatalf("init error: %v", err)
	}

	mem.Print()
	mem.MustHave(t, "created a local repository")
}

func TestAlreadyExists(t *testing.T) {
	cmd.SetupTestRepo(t)
	defer cmd.TeardownTestRepo(t)
	logx.Mock()

	args := []string{}
	_ = Init(args)
	err := Init(args)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("unexpected error: %v", err)
	}
}
