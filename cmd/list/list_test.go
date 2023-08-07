package list

import (
	"testing"

	"sqlpkg.org/cli/cmd"
	"sqlpkg.org/cli/lockfile"
	"sqlpkg.org/cli/logx"
)

func TestList(t *testing.T) {
	_, lockPath := cmd.SetupTestRepo(t)
	defer cmd.TeardownTestRepo(t)
	cmd.CopyTestRepo(t, "")
	mem := logx.Mock()

	args := []string{}
	err := List(args)
	if err != nil {
		t.Fatalf("list error: %v", err)
	}

	mem.Print()
	mem.MustHave(t, "gathered 2 packages")
	mem.MustHave(t, "added 2 packages to the lockfile")

	lck, err := lockfile.ReadLocal(lockPath)
	if err != nil {
		t.Fatal("failed to read lockfile")
	}
	if len(lck.Packages) != 2 {
		t.Fatalf("unexpected package count: %v", len(lck.Packages))
	}
	if !lck.Has("nalgeon/example") {
		t.Fatal("nalgeon/example not found in the lockfile")
	}
	if !lck.Has("sqlite/stmt") {
		t.Fatal("sqlite/stmt not found in the lockfile")
	}
}
