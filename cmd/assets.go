// Commands that manage package assets.
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"sqlpkg.org/cli/assets"
	"sqlpkg.org/cli/fileio"
	"sqlpkg.org/cli/logx"
	"sqlpkg.org/cli/spec"
)

// BuildAssetPath constructs an URL to download package asset.
func BuildAssetPath(pkg *spec.Package) (*spec.AssetPath, error) {
	logx.Debug("checking remote asset for platform %s-%s", runtime.GOOS, runtime.GOARCH)
	logx.Debug("asset base path = %s", pkg.Assets.Path)

	assetPath, err := pkg.AssetPath(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return nil, fmt.Errorf("unsupported platform: %s-%s", runtime.GOOS, runtime.GOARCH)
	}

	if !assetPath.Exists() {
		return nil, fmt.Errorf("asset does not exist: %s", assetPath)
	}

	return assetPath, nil
}

// DownloadAsset downloads package asset.
func DownloadAsset(pkg *spec.Package, assetPath *spec.AssetPath) (*assets.Asset, error) {
	logx.Debug("downloading %s", assetPath)
	dir := spec.Dir(os.TempDir(), pkg.Owner, pkg.Name)
	err := fileio.CreateDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	var asset *assets.Asset
	if assetPath.IsRemote {
		asset, err = assets.Download(dir, assetPath.Value)
	} else {
		asset, err = assets.Copy(dir, assetPath.Value)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to download asset: %w", err)
	}

	sizeKb := float64(asset.Size) / 1024
	logx.Debug("downloaded %s (%.2f Kb)", asset.Name, sizeKb)
	return asset, nil
}

// ValidateAsset checks if the asset is valid.
func ValidateAsset(pkg *spec.Package, asset *assets.Asset) error {
	checksumStr, ok := pkg.Assets.Checksums[asset.Name]
	if !ok {
		logx.Debug("spec is missing asset checksum")
		return nil
	}

	ok, err := asset.Validate(checksumStr)
	if err != nil {
		return fmt.Errorf("failed to validate asset: %w", err)
	}

	if !ok {
		return fmt.Errorf("asset checksum is invalid")
	}

	logx.Debug("asset checksum is valid")
	return nil
}

// UnpackAsset unpacks package asset.
func UnpackAsset(pkg *spec.Package, asset *assets.Asset) error {
	nFiles, err := assets.Unpack(asset.Path, pkg.Assets.Pattern)
	if err != nil {
		return fmt.Errorf("failed to unpack asset: %w", err)
	}
	if nFiles == 0 {
		logx.Debug("not an archive, skipping unpack: %s", asset.Name)
		return nil
	}
	err = os.Remove(asset.Path)
	if err != nil {
		return fmt.Errorf("failed to delete asset after unpacking: %w", err)
	}
	logx.Debug("unpacked %d files from %s", nFiles, asset.Name)
	return nil
}

// InstallFiles installes unpacked package files.
func InstallFiles(pkg *spec.Package, asset *assets.Asset) error {
	pkgDir := spec.Dir(WorkDir, pkg.Owner, pkg.Name)
	err := fileio.MoveDir(asset.Dir(), pkgDir)
	if err != nil {
		return fmt.Errorf("failed to copy downloaded files: %w", err)
	}

	err = pkg.Save(pkgDir)
	if err != nil {
		return fmt.Errorf("failed to write package spec: %w", err)
	}

	return nil
}

// DequarantineFiles removes the macOS quarantine flag
// from all *.dylib files in the package directory.
func DequarantineFiles(pkg *spec.Package) error {
	if runtime.GOOS != "darwin" {
		return nil
	}

	pattern := filepath.Join(spec.Dir(WorkDir, pkg.Owner, pkg.Name), "*.dylib")
	paths, _ := filepath.Glob(pattern)
	if len(paths) == 0 {
		return nil
	}

	var allErr error
	for _, path := range paths {
		err := fileio.Dequarantine(path)
		allErr = errors.Join(allErr, err)
	}
	if allErr != nil {
		return fmt.Errorf("failed to dequarantine files: %w", allErr)
	}

	logx.Debug("removed %d files from quarantine", len(paths))
	return nil
}
