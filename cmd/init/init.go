package init

import (
	"errors"
	"fmt"
	"os"

	"sqlpkg.org/cli/cmd"
	"sqlpkg.org/cli/fileio"
	"sqlpkg.org/cli/spec"
)

const initHelp = "usage: sqlpkg init"

// Init creates an empty local package repository.
func Init(args []string) error {
	if len(args) != 0 {
		return errors.New(initHelp)
	}

	if fileio.Exists(spec.DirName) {
		return errors.New(".sqlpkg dir already exists")
	}

	err := os.Mkdir(spec.DirName, 0755)
	if err != nil {
		return fmt.Errorf("failed to create a local repository: %w", err)
	}

	cmd.Log("✓ created a local repository")
	return nil
}
