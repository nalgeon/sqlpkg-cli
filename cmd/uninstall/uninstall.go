package uninstall

import (
	"errors"
	"fmt"
	"os"

	"sqlpkg.org/cli/cmd"
	"sqlpkg.org/cli/fileio"
)

const uninstallHelp = "usage: sqlpkg uninstall package"

// Uninstall deletes the specified package.
func Uninstall(args []string) error {
	if len(args) != 1 {
		return errors.New(uninstallHelp)
	}

	cmd.PrintLocalRepo()

	fullName := args[0]
	cmd.Log("> uninstalling %s...", fullName)

	err := removePackageDir(fullName)
	if err != nil {
		return err
	}

	lck, err := cmd.ReadLockfile()
	if err != nil {
		return err
	}

	err = cmd.RemoveFromLockfile(lck, fullName)
	if err != nil {
		return err
	}

	cmd.Log("✓ uninstalled package %s", fullName)
	return nil
}

func removePackageDir(fullName string) error {
	dir, err := cmd.GetDirByFullName(fullName)
	if err != nil {
		return err
	}

	cmd.Debug("checking dir: %s", dir)
	if !fileio.Exists(dir) {
		cmd.Debug("package dir not found")
		return errors.New("package is not installed")
	}

	cmd.Debug("deleting dir: %s", dir)
	err = os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("failed to delete package dir: %w", err)
	}

	cmd.Debug("deleted package dir")
	return nil
}
