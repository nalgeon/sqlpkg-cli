package help

import (
	"testing"

	"sqlpkg.org/cli/cmd"
)

func TestHelp(t *testing.T) {
	cmd.WorkDir = "."
	repoDir, lockPath := cmd.SetupTestRepo(t)
	mem := cmd.SetupTestLogger()

	args := []string{}
	err := Help(args)
	if err != nil {
		t.Fatalf("help error: %v", err)
	}

	mem.Print()
	mem.MustHave(t, "sqlpkg is an SQLite package manager")
	mem.MustHave(t, "install")
	mem.MustHave(t, "uninstall")
	mem.MustHave(t, "update")
	mem.MustHave(t, "list")
	mem.MustHave(t, "init")
	mem.MustHave(t, "info")
	mem.MustHave(t, "which")
	mem.MustHave(t, "help")
	mem.MustHave(t, "version")

	cmd.TeardownTestRepo(t, repoDir, lockPath)
}
