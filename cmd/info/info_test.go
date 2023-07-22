package info

import (
	"testing"

	"sqlpkg.org/cli/cmd"
)

func TestInfo(t *testing.T) {
	cmd.WorkDir = "."
	repoDir, lockPath := cmd.SetupTestRepo(t)
	cmd.CopyTestRepo(t, "")
	mem := cmd.SetupTestLogger()

	args := []string{"nalgeon/example"}
	err := Info(args)
	if err != nil {
		t.Fatalf("info error: %v", err)
	}

	mem.Print()
	mem.MustHave(t, "nalgeon/example@0.1.0 by Anton Zhiyanov")
	mem.MustHave(t, "Example extension")
	mem.MustHave(t, "https://github.com/nalgeon/sqlite-example")
	mem.MustHave(t, "license: MIT")
	mem.MustHave(t, "✓ installed")

	cmd.TeardownTestRepo(t, repoDir, lockPath)
}
