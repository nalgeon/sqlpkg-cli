package info

import (
	"testing"

	"sqlpkg.org/cli/cmd"
	"sqlpkg.org/cli/logx"
)

func TestInfo(t *testing.T) {
	cmd.SetupTestRepo(t)
	defer cmd.TeardownTestRepo(t)
	cmd.CopyTestRepo(t, "")
	mem := logx.Mock()

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
}
