package update

import (
	"errors"
	"fmt"
	"path/filepath"

	"sqlpkg.org/cli/cmd"
	"sqlpkg.org/cli/lockfile"
	"sqlpkg.org/cli/spec"
)

const updateHelp = "usage: sqlpkg update [package]"

// UpdateAll updates installed packages to latest versions.
func UpdateAll(args []string) error {
	if len(args) != 0 {
		return errors.New(updateHelp)
	}

	cmd.PrintLocalRepo()

	pattern := filepath.Join(cmd.WorkDir, spec.DirName, "*", "*", spec.FileName)
	paths, _ := filepath.Glob(pattern)

	if len(paths) == 0 {
		fmt.Println("no packages installed")
		return nil
	}

	lck, err := cmd.ReadLockfile()
	if err != nil {
		return err
	}

	count := 0
	for _, path := range paths {
		pkg, err := spec.ReadLocal(path)
		if err != nil {
			cmd.Log("! invalid package %s: %s", path, err)
			continue
		}
		cmd.Debug("found local spec from %s", path)
		cmd.Debug("read package %s, version = %s", pkg.FullName(), pkg.Version)

		cmd.Log("> updating %s...", pkg.FullName())
		updPkg, err := updatePackage(lck, getSpecPath(pkg))
		if err != nil {
			cmd.Log("! error updating %s: %s", pkg.FullName(), err)
			continue
		}
		if updPkg == nil {
			cmd.Log("✓ already at the latest version")
			continue
		}
		updVersion := updPkg.Version
		if updVersion == "" {
			updVersion = "latest version"
		}
		cmd.Log("✓ updated package %s to %s", updPkg.FullName(), updVersion)
		count += 1
	}

	cmd.Log("updated %d packages", count)
	return nil
}

// Update updates a specific package to the latest version.
func Update(args []string) error {
	if len(args) != 1 {
		return errors.New(updateHelp)
	}

	fullName := args[0]
	path, err := cmd.GetPathByFullName(fullName)
	if err != nil {
		return err
	}

	pkg, err := spec.ReadLocal(path)
	if err != nil {
		return fmt.Errorf("invalid package: %w", err)
	}
	cmd.Debug("found local spec from %s", path)
	cmd.Debug("read package %s, version = %s", pkg.FullName(), pkg.Version)

	lck, err := cmd.ReadLockfile()
	if err != nil {
		return err
	}

	cmd.Log("> updating %s...", pkg.FullName())
	updPkg, err := updatePackage(lck, getSpecPath(pkg))
	if err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}

	if updPkg == nil {
		cmd.Log("✓ already at the latest version")
		return nil
	}

	cmd.Log("✓ updated package %s to %s", updPkg.FullName(), updPkg.Version)
	return nil
}

// updatePackage updates a package.
// Returns true if the package was actually updated, false otherwise
// (already at the latest version or encountered an error).
func updatePackage(lck *lockfile.Lockfile, path string) (*spec.Package, error) {
	cmd.Debug("using spec path: %s", path)
	pkg, err := cmd.ReadSpec(path)
	if err != nil {
		return nil, err
	}

	err = cmd.ResolveVersion(pkg)
	if err != nil {
		return nil, err
	}

	if !cmd.HasNewVersion(pkg) {
		return nil, nil
	}

	err = cmd.ReadChecksums(pkg)
	if err != nil {
		return nil, err
	}

	assetUrl, err := cmd.BuildAssetPath(pkg)
	if err != nil {
		return nil, err
	}

	asset, err := cmd.DownloadAsset(pkg, assetUrl)
	if err != nil {
		return nil, err
	}

	err = cmd.ValidateAsset(pkg, asset)
	if err != nil {
		return nil, err
	}

	err = cmd.UnpackAsset(pkg, asset)
	if err != nil {
		return nil, err
	}

	err = cmd.InstallFiles(pkg, asset)
	if err != nil {
		return nil, err
	}

	err = cmd.DequarantineFiles(pkg)
	if err != nil {
		return nil, err
	}

	err = cmd.AddToLockfile(lck, pkg)
	if err != nil {
		return nil, err
	}

	return pkg, nil
}

// getSpecPath returns a remote package spec path.
func getSpecPath(pkg *spec.Package) string {
	if pkg.Specfile != "" {
		return pkg.Specfile
	}
	// in older specs the .Specfile may be empty
	return pkg.FullName()
}
