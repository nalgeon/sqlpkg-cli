package uninstall

import (
	"errors"
	"fmt"
	"os"

	"sqlpkg.org/cli/cmd"
	"sqlpkg.org/cli/fileio"
	"sqlpkg.org/cli/logx"
)

const uninstallHelp = "usage: sqlpkg uninstall package"

// Uninstall deletes the specified package.
func Uninstall(args []string) error {
	if len(args) != 1 {
		return errors.New(uninstallHelp)
	}

	cmd.PrintScope()

	fullName := args[0]
	logx.Log("> uninstalling %s...", fullName)

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

	logx.Log("✓ uninstalled package %s", fullName)
	return nil
}

func removePackageDir(fullName string) error {
	dir, err := cmd.GetDirByFullName(fullName)
	if err != nil {
		return err
	}

	logx.Debug("checking dir: %s", dir)
	if !fileio.Exists(dir) {
		logx.Debug("package dir not found")
		return errors.New("package is not installed")
	}

	logx.Debug("deleting dir: %s", dir)
	err = os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("failed to delete package dir: %w", err)
	}

	logx.Debug("deleted package dir")
	return nil
}
