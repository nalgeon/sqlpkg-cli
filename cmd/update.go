package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/nalgeon/sqlpkg-cli/internal/metadata"
)

const updateHelp = "usage: sqlpkg update [package]"

// UpdateAll updates installed packages to latest versions.
func UpdateAll(args []string) error {
	if len(args) != 0 {
		return errors.New(updateHelp)
	}

	pattern := fmt.Sprintf("%s/%s/*/*/%s", workDir, metadata.DirName, metadata.FileName)
	paths, _ := filepath.Glob(pattern)

	if len(paths) == 0 {
		fmt.Println("no packages installed")
		return nil
	}

	count := 0
	for _, path := range paths {
		pkg, err := metadata.ReadLocal(path)
		if err != nil {
			log("! invalid package %s: %s", path, err)
			continue
		}

		log("> updating %s...", pkg.FullName())
		updated, err := updatePackage(pkg)
		if err != nil {
			log("! error updating %s: %s", pkg.FullName(), err)
			continue
		}
		if !updated {
			log("✓ already at the latest version")
			continue
		}
		log("✓ updated package %s to %s", pkg.FullName(), pkg.Version)
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

	pkg, err := metadata.ReadLocal(path)
	if err != nil {
		return fmt.Errorf("invalid package: %w", err)
	}

	log("> updating %s...", pkg.FullName())
	updated, err := updatePackage(pkg)
	if err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}

	if updated {
		log("✓ updated package %s to %s", pkg.FullName(), pkg.Version)
	} else {
		log("✓ already at the latest version")
	}
	return nil
}

// updatePackage updates a package.
// Returns true if the package was actually updated, false otherwise
// (already at the latest version or encountered an error).
func updatePackage(pkg *metadata.Package) (bool, error) {
	cmd := new(command)
	cmd.readMetadata(pkg.FullName())
	if !cmd.hasNewVersion() {
		return false, nil
	}
	assetUrl := cmd.buildAssetPath()
	asset := cmd.downloadAsset(assetUrl)
	cmd.unpackAsset(asset)
	cmd.installFiles()
	return cmd.err != nil, cmd.err
}
