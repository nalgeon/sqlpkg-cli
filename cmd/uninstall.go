package cmd

import (
	"errors"
)

const uninstallHelp = "usage: sqlpkg uninstall package"

// Uninstall deletes the specified package.
func Uninstall(args []string) error {
	if len(args) != 1 {
		return errors.New(uninstallHelp)
	}

	printLocalRepo()

	fullName := args[0]
	log("> uninstalling %s...", fullName)

	err := removePackageDir(fullName)
	if err != nil {
		return err
	}

	lck, err := readLockfile()
	if err != nil {
		return err
	}

	err = removeFromLockfile(lck, fullName)
	if err != nil {
		return err
	}

	log("✓ uninstalled package %s", fullName)
	return nil
}
