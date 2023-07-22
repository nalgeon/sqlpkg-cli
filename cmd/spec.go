// Commands that manage package spec files.
package cmd

import (
	"fmt"

	"sqlpkg.org/cli/checksums"
	"sqlpkg.org/cli/logx"
	"sqlpkg.org/cli/spec"
)

// ReadSpec reads package spec.
func ReadSpec(path string) (*spec.Package, error) {
	pkg, err := spec.Read(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read package spec: %w", err)
	}
	pkg.ExpandVars()
	logx.Debug("found package spec at %s", pkg.Specfile)
	logx.Debug("read package %s, version = %s", pkg.FullName(), pkg.Version)
	return pkg, nil
}

// FindSpec loads the package spec, giving preference to already installed packages.
func FindSpec(path string) (*spec.Package, error) {
	pkg := ReadInstalledSpec(path)
	if pkg != nil {
		return pkg, nil
	}

	logx.Debug("package is not installed")
	pkg, err := ReadSpec(path)
	return pkg, err
}

// ReadInstalledSpec loads the package spec for an installed package (if any).
func ReadInstalledSpec(fullName string) *spec.Package {
	path, err := GetPathByFullName(fullName)
	if err != nil {
		return nil
	}

	pkg, err := spec.ReadLocal(path)
	if err != nil {
		return nil
	}

	logx.Debug("found installed package")
	return pkg
}

// ReadChecksums reads package asset checksums from the checksum file.
func ReadChecksums(pkg *spec.Package) error {
	path := pkg.Assets.Path.Join(checksums.FileName)
	if !checksums.Exists(path.Value, path.IsRemote) {
		logx.Debug("missing spec checksum file")
		return nil
	}
	sums, err := checksums.Read(path.Value, path.IsRemote)
	if err != nil {
		return fmt.Errorf("failed to read checksum file: %w", err)
	}
	logx.Debug("read %d checksums", len(sums))
	pkg.Assets.Checksums = sums
	return nil
}
