package cmd

import (
	"errors"
	"fmt"

	"github.com/nalgeon/sqlpkg-cli/internal/spec"
)

const installHelp = "usage: sqlpkg install [package]"

// InstallAll installs all packages from the lockfile.
func InstallAll(args []string) error {
	printLocalRepo()

	lck, err := readLockfile()
	if err != nil {
		return err
	}
	debug("loaded the lockfile with %d packages", len(lck.Packages))

	if len(lck.Packages) == 0 {
		log("no packages found in the lockfile")
		return nil
	}

	errCount := 0
	for _, pkg := range lck.Packages {
		path := pkg.Specfile
		if path == "" {
			debug("missing specfile for %s, falling back to name/owner", pkg.FullName())
			path = pkg.FullName()
		}
		err = installPackage(path)
		if err != nil {
			errCount += 1
			log("! %s", err)
		}
	}

	if errCount > 0 {
		return fmt.Errorf("failed to install %d packages", errCount)
	}
	log("installed %d packages", len(lck.Packages))
	return nil
}

// Install installs a new package or updates an existing one.
func Install(args []string) error {
	if len(args) != 1 {
		return errors.New(installHelp)
	}

	printLocalRepo()

	path := args[0]
	err := installPackage(path)
	return err
}

// installPackage installs a package using a specfile from a given path.
func installPackage(path string) error {
	log("> installing %s...", path)

	pkg, err := readSpec(path)
	if err != nil {
		return err
	}

	if !hasNewVersion(pkg) {
		log("✓ already at the latest version")
		return nil
	}

	assetPath, err := buildAssetPath(pkg)
	if err != nil {
		return err
	}

	asset, err := downloadAsset(pkg, assetPath)
	if err != nil {
		return err
	}

	err = validateAsset(pkg, asset)
	if err != nil {
		return err
	}

	err = unpackAsset(pkg, asset)
	if err != nil {
		return err
	}

	err = installFiles(pkg, asset)
	if err != nil {
		return err
	}

	lck, err := readLockfile()
	if err != nil {
		return err
	}

	err = addToLockfile(lck, pkg)
	if err != nil {
		return err
	}

	dir := spec.Dir(workDir, pkg.Owner, pkg.Name)
	log("✓ installed package %s to %s", pkg.FullName(), dir)
	return nil
}
