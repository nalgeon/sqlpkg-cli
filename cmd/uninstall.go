package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/nalgeon/sqlpkg-cli/internal/fileio"
)

const uninstallHelp = "usage: sqlpkg uninstall package"

// Uninstall deletes the specified package.
func Uninstall(args []string) error {
	if len(args) != 1 {
		return errors.New(uninstallHelp)
	}

	fullName := args[0]
	dir, err := getDirByFullName(fullName)
	if err != nil {
		return err
	}

	log("> uninstalling %s...", fullName)
	debug("checking dir: %s", dir)
	if !fileio.Exists(dir) {
		return errors.New("package is not installed")
	}

	debug("deleting dir: %s", dir)
	err = os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("uninstall failed: %w", err)
	}

	log("✓ uninstalled package %s", fullName)
	return nil
}
