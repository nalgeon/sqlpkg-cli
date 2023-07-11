package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/nalgeon/sqlpkg-cli/internal/lockfile"
	"github.com/nalgeon/sqlpkg-cli/internal/spec"
)

const updateHelp = "usage: sqlpkg update [package]"

// UpdateAll updates installed packages to latest versions.
func UpdateAll(args []string) error {
	if len(args) != 0 {
		return errors.New(updateHelp)
	}

	printLocalRepo()

	pattern := fmt.Sprintf("%s/%s/*/*/%s", workDir, spec.DirName, spec.FileName)
	paths, _ := filepath.Glob(pattern)

	if len(paths) == 0 {
		fmt.Println("no packages installed")
		return nil
	}

	lck, err := readLockfile()
	if err != nil {
		return err
	}

	count := 0
	for _, path := range paths {
		pkg, err := spec.ReadLocal(path)
		if err != nil {
			log("! invalid package %s: %s", path, err)
			continue
		}

		log("> updating %s...", pkg.FullName())
		updPkg, err := updatePackage(lck, pkg.FullName())
		if err != nil {
			log("! error updating %s: %s", pkg.FullName(), err)
			continue
		}
		if updPkg == nil {
			log("✓ already at the latest version")
			continue
		}
		log("✓ updated package %s to %s", updPkg.FullName(), updPkg.Version)
		count += 1
	}

	log("updated %d packages", count)
	return nil
}

// Update updates a specific package to the latest version.
func Update(args []string) error {
	if len(args) != 1 {
		return errors.New(updateHelp)
	}

	fullName := args[0]
	path, err := getPathByFullName(fullName)
	if err != nil {
		return err
	}

	pkg, err := spec.ReadLocal(path)
	if err != nil {
		return fmt.Errorf("invalid package: %w", err)
	}

	lck, err := readLockfile()
	if err != nil {
		return err
	}

	log("> updating %s...", pkg.FullName())
	updPkg, err := updatePackage(lck, pkg.FullName())
	if err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}

	if updPkg == nil {
		log("✓ already at the latest version")
		return nil
	}

	log("✓ updated package %s to %s", updPkg.FullName(), updPkg.Version)
	return nil
}

// updatePackage updates a package.
// Returns true if the package was actually updated, false otherwise
// (already at the latest version or encountered an error).
func updatePackage(lck *lockfile.Lockfile, fullName string) (*spec.Package, error) {
	pkg, err := readSpec(fullName)
	if err != nil {
		return nil, err
	}
	if !hasNewVersion(pkg) {
		return nil, nil
	}

	assetUrl, err := buildAssetPath(pkg)
	if err != nil {
		return nil, err
	}

	asset, err := downloadAsset(pkg, assetUrl)
	if err != nil {
		return nil, err
	}

	err = validateAsset(pkg, asset)
	if err != nil {
		return nil, err
	}

	err = unpackAsset(pkg, asset)
	if err != nil {
		return nil, err
	}

	err = installFiles(pkg, asset)
	if err != nil {
		return nil, err
	}

	err = addToLockfile(lck, pkg)
	if err != nil {
		return nil, err
	}

	return pkg, nil
}
