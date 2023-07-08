package cmd

import (
	"errors"

	"github.com/nalgeon/sqlpkg-cli/internal/spec"
)

const installHelp = "usage: sqlpkg install package"

// Install installs a new package or updates an existing one.
func Install(args []string) error {
	if len(args) != 1 {
		return errors.New(installHelp)
	}

	path := args[0]
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
