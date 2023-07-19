package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/nalgeon/sqlpkg-cli/internal/assets"
	"github.com/nalgeon/sqlpkg-cli/internal/checksums"
	"github.com/nalgeon/sqlpkg-cli/internal/fileio"
	"github.com/nalgeon/sqlpkg-cli/internal/lockfile"
	"github.com/nalgeon/sqlpkg-cli/internal/spec"
	"golang.org/x/mod/semver"
)

// workDir is the current working directory.
var workDir string

// userHomeDir is the user's home directory.
var userHomeDir string

// init determines the working directory.
// It is either the .sqlpkg directory (if present) or ~/.sqlpkg otherwise.
func init() {
	if fileio.Exists(spec.DirName) {
		workDir = "."
		return
	}
	var err error
	userHomeDir, err = os.UserHomeDir()
	if err != nil {
		workDir = "."
		return
	}
	workDir = userHomeDir
}

// readSpec reads package spec.
func readSpec(path string) (*spec.Package, error) {
	pkg, err := spec.Read(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read package spec: %w", err)
	}
	pkg.ExpandVars()
	debug("found package spec at %s", pkg.Specfile)
	debug("read package %s, version = %s", pkg.FullName(), pkg.Version)
	return pkg, nil
}

// readChecksums reads package asset checksums from the checksum file.
func readChecksums(pkg *spec.Package) error {
	path := pkg.Assets.Path.Join(checksums.FileName)
	if !checksums.Exists(path.Value, path.IsRemote) {
		debug("missing spec checksum file")
		return nil
	}
	sums, err := checksums.Read(path.Value, path.IsRemote)
	if err != nil {
		return fmt.Errorf("failed to read checksum file: %w", err)
	}
	debug("read %d checksums", len(sums))
	pkg.Assets.Checksums = sums
	return nil
}

// isInstalled checks if there is a local package installed.
func isInstalled(pkg *spec.Package) bool {
	path := spec.Path(workDir, pkg.Owner, pkg.Name)
	return fileio.Exists(path)
}

// hasNewVersion checks if the remote package is newer than the local one.
func hasNewVersion(pkg *spec.Package) bool {
	oldPath := spec.Path(workDir, pkg.Owner, pkg.Name)
	if !fileio.Exists(oldPath) {
		return true
	}

	oldPkg, err := spec.ReadLocal(oldPath)
	if err != nil {
		return true
	}
	debug("local package version = %s", oldPkg.Version)

	if oldPkg.Version == "" {
		// not explicitly versioned, always assume there is a later version
		return true
	}

	if oldPkg.Version == pkg.Version {
		return false
	}

	if semver.Compare(oldPkg.Version, pkg.Version) < 0 {
		return false
	}

	return true
}

// buildAssetPath constructs an URL to download package asset.
func buildAssetPath(pkg *spec.Package) (*spec.AssetPath, error) {
	debug("checking remote asset for platform %s-%s", runtime.GOOS, runtime.GOARCH)
	debug("asset base path = %s", pkg.Assets.Path)

	assetPath, err := pkg.AssetPath(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return nil, fmt.Errorf("unsupported platform: %s-%s", runtime.GOOS, runtime.GOARCH)
	}

	if !assetPath.Exists() {
		return nil, fmt.Errorf("asset does not exist: %s", assetPath)
	}

	return assetPath, nil
}

// downloadAsset downloads package asset.
func downloadAsset(pkg *spec.Package, assetPath *spec.AssetPath) (*assets.Asset, error) {
	debug("downloading %s", assetPath)
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
	debug("downloaded %s (%.2f Kb)", asset.Name, sizeKb)
	return asset, nil
}

// validateAsset checks if the asset is valid.
func validateAsset(pkg *spec.Package, asset *assets.Asset) error {
	checksumStr, ok := pkg.Assets.Checksums[asset.Name]
	if !ok {
		debug("spec is missing asset checksum")
		return nil
	}

	ok, err := asset.Validate(checksumStr)
	if err != nil {
		return fmt.Errorf("failed to validate asset: %w", err)
	}

	if !ok {
		return fmt.Errorf("asset checksum is invalid")
	}

	debug("asset checksum is valid")
	return nil
}

// unpackAsset unpacks package asset.
func unpackAsset(pkg *spec.Package, asset *assets.Asset) error {
	nFiles, err := assets.Unpack(asset.Path, pkg.Assets.Pattern)
	if err != nil {
		return fmt.Errorf("failed to unpack asset: %w", err)
	}
	if nFiles == 0 {
		debug("not an archive, skipping unpack: %s", asset.Name)
		return nil
	}
	err = os.Remove(asset.Path)
	if err != nil {
		return fmt.Errorf("failed to delete asset after unpacking: %w", err)
	}
	debug("unpacked %d files from %s", nFiles, asset.Name)
	return nil
}

// installFiles installes unpacked package files.
func installFiles(pkg *spec.Package, asset *assets.Asset) error {
	pkgDir := spec.Dir(workDir, pkg.Owner, pkg.Name)
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

// dequarantineFiles removes the macOS quarantine flag
// from all *.dylib files in the package directory.
func dequarantineFiles(pkg *spec.Package) error {
	if runtime.GOOS != "darwin" {
		return nil
	}

	pattern := filepath.Join(spec.Dir(workDir, pkg.Owner, pkg.Name), "*.dylib")
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

	debug("removed %d files from quarantine", len(paths))
	return nil
}

// readLockfile reads lockfile from the work directory.
func readLockfile() (*lockfile.Lockfile, error) {
	path := lockfile.Path(workDir)
	if !fileio.Exists(path) {
		debug("created new lockfile")
		return lockfile.NewLockfile(), nil
	}

	lck, err := lockfile.ReadLocal(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read lockfile: %w", err)
	}

	debug("read existing lockfile")
	return lck, nil
}

// addToLockfile adds package to the lockfile.
func addToLockfile(lck *lockfile.Lockfile, pkg *spec.Package) error {
	lck.Add(pkg)
	err := lck.Save(workDir)
	if err != nil {
		return fmt.Errorf("failed to save lockfile: %w", err)
	}

	debug("added package to the lockfile")
	return nil
}

// removeFromLockfile removes package from the lockfile.
func removeFromLockfile(lck *lockfile.Lockfile, fullName string) error {
	pkg, ok := lck.Packages[fullName]
	if !ok {
		debug("package not listed in the lockfile")
		return nil
	}

	lck.Remove(pkg)
	err := lck.Save(workDir)
	if err != nil {
		return fmt.Errorf("failed to save lockfile: %w", err)
	}

	debug("removed package from the lockfile")
	return nil
}

func removePackageDir(fullName string) error {
	dir, err := getDirByFullName(fullName)
	if err != nil {
		return err
	}

	debug("checking dir: %s", dir)
	if !fileio.Exists(dir) {
		debug("package dir not found")
		return errors.New("package is not installed")
	}

	debug("deleting dir: %s", dir)
	err = os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("failed to delete package dir: %w", err)
	}

	debug("deleted package dir")
	return nil
}

// getDirByFullName expands an owner-name package pair to a full package dir.
func getDirByFullName(fullName string) (string, error) {
	parts := strings.Split(fullName, "/")
	if len(parts) != 2 {
		return "", errors.New("invalid package name")
	}
	path := spec.Dir(workDir, parts[0], parts[1])
	return path, nil
}

// getPathFullName expands an owner-name package pair to a full sqlpkg.json path.
func getPathByFullName(fullName string) (string, error) {
	parts := strings.Split(fullName, "/")
	if len(parts) != 2 {
		return "", errors.New("invalid package name")
	}
	path := spec.Path(workDir, parts[0], parts[1])
	return path, nil
}

var IsVerbose bool

// log prints a message to the screen.
func log(message string, args ...any) {
	if len(args) == 0 {
		fmt.Println(message)
	} else {
		fmt.Printf(message+"\n", args...)
	}
}

// debug prints a message to the screen if the verbose mode is on.
func debug(message string, args ...any) {
	if !IsVerbose {
		return
	}
	log(".."+message, args...)
}
