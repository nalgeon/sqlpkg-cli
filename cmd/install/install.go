package install

import (
	"errors"
	"fmt"

	"sqlpkg.org/cli/cmd"
	"sqlpkg.org/cli/spec"
)

const installHelp = "usage: sqlpkg install [package]"

// InstallAll installs all packages from the lockfile.
func InstallAll(args []string) error {
	cmd.PrintLocalRepo()

	lck, err := cmd.ReadLockfile()
	if err != nil {
		return err
	}
	cmd.Debug("loaded the lockfile with %d packages", len(lck.Packages))

	if len(lck.Packages) == 0 {
		cmd.Log("no packages found in the lockfile")
		return nil
	}

	errCount := 0
	for _, pkg := range lck.Packages {
		err = installLockedPackage(pkg)
		if err != nil {
			errCount += 1
			cmd.Log("! %s", err)
		}
	}

	if errCount > 0 {
		return fmt.Errorf("failed to install %d packages", errCount)
	}
	cmd.Log("installed %d packages", len(lck.Packages))
	return nil
}

// Install installs a new package or updates an existing one.
func Install(args []string) error {
	if len(args) != 1 {
		return errors.New(installHelp)
	}

	cmd.PrintLocalRepo()

	path := args[0]
	err := installPackage(path)
	return err
}

// installPackage installs a package using a specfile from a given path.
func installPackage(path string) error {
	cmd.Log("> installing %s...", path)

	pkg, err := cmd.ReadSpec(path)
	if err != nil {
		return err
	}

	err = cmd.ResolveVersion(pkg)
	if err != nil {
		return err
	}

	if !cmd.HasNewVersion(pkg) {
		cmd.Log("✓ already at the latest version")
		return nil
	}

	err = cmd.ReadChecksums(pkg)
	if err != nil {
		return err
	}

	assetPath, err := cmd.BuildAssetPath(pkg)
	if err != nil {
		return err
	}

	asset, err := cmd.DownloadAsset(pkg, assetPath)
	if err != nil {
		return err
	}

	err = cmd.ValidateAsset(pkg, asset)
	if err != nil {
		return err
	}

	err = cmd.UnpackAsset(pkg, asset)
	if err != nil {
		return err
	}

	err = cmd.InstallFiles(pkg, asset)
	if err != nil {
		return err
	}

	err = cmd.DequarantineFiles(pkg)
	if err != nil {
		return err
	}

	lck, err := cmd.ReadLockfile()
	if err != nil {
		return err
	}

	err = cmd.AddToLockfile(lck, pkg)
	if err != nil {
		return err
	}

	dir := spec.Dir(cmd.WorkDir, pkg.Owner, pkg.Name)
	cmd.Log("✓ installed package %s to %s", pkg.FullName(), dir)
	return nil
}

// installLockedPackage installs a specific version of a package from the lockfile.
func installLockedPackage(lckPkg *spec.Package) error {
	path := lckPkg.Specfile
	if path == "" {
		cmd.Debug("missing specfile for %s, falling back to name/owner", lckPkg.FullName())
		path = lckPkg.FullName()
	}

	cmd.Log("> installing %s...", path)

	pkg, err := cmd.ReadSpec(path)
	if err != nil {
		return err
	}

	// lock the version
	cmd.Debug("locked version = %s", lckPkg.Version)
	pkg.Version = lckPkg.Version
	pkg.Assets = lckPkg.Assets

	if !cmd.HasNewVersion(pkg) {
		cmd.Log("✓ already at the %s version", pkg.Version)
		return nil
	}

	assetPath, err := cmd.BuildAssetPath(pkg)
	if err != nil {
		return err
	}

	asset, err := cmd.DownloadAsset(pkg, assetPath)
	if err != nil {
		return err
	}

	err = cmd.ValidateAsset(pkg, asset)
	if err != nil {
		return err
	}

	err = cmd.UnpackAsset(pkg, asset)
	if err != nil {
		return err
	}

	err = cmd.InstallFiles(pkg, asset)
	if err != nil {
		return err
	}

	err = cmd.DequarantineFiles(pkg)
	if err != nil {
		return err
	}

	// no need to add the package to the lockfile,
	// it's already there

	dir := spec.Dir(cmd.WorkDir, pkg.Owner, pkg.Name)
	cmd.Log("✓ installed package %s to %s", pkg.FullName(), dir)
	return nil
}
