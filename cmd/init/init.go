package init

import (
	"errors"
	"fmt"
	"os"

	"sqlpkg.org/cli/fileio"
	"sqlpkg.org/cli/logx"
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
		return fmt.Errorf("failed to create a project scope: %w", err)
	}

	logx.Log("✓ created a project scope")
	return nil
}
