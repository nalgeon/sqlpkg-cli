package help

import (
	"testing"

	"sqlpkg.org/cli/cmd"
	"sqlpkg.org/cli/logx"
)

func TestHelp(t *testing.T) {
	cmd.SetupTestRepo(t)
	defer cmd.TeardownTestRepo(t)
	mem := logx.Mock()

	args := []string{}
	err := Help(args)
	if err != nil {
		t.Fatalf("help error: %v", err)
	}

	mem.Print()
	mem.MustHave(t, "sqlpkg is a package manager")
	mem.MustHave(t, "install")
	mem.MustHave(t, "uninstall")
	mem.MustHave(t, "update")
	mem.MustHave(t, "list")
	mem.MustHave(t, "init")
	mem.MustHave(t, "info")
	mem.MustHave(t, "which")
	mem.MustHave(t, "help")
	mem.MustHave(t, "version")
}
