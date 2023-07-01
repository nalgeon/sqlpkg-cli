package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/nalgeon/sqlpkg-cli/internal/assets"
	"github.com/nalgeon/sqlpkg-cli/internal/fileio"
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

type command struct {
	pkg *spec.Package
	dir string
	err error
}

// readSpec reads package spec.
func (cmd *command) readSpec(path string) {
	if cmd.err != nil {
		return
	}

	var err error
	cmd.pkg, err = spec.Read(path)
	if err != nil {
		cmd.err = fmt.Errorf("failed to read package spec: %w", err)
		return
	}
	cmd.pkg.ExpandVars()
	debug("found package spec at %s", cmd.pkg.Path)
	debug("read package %s, version = %s", cmd.pkg.FullName(), cmd.pkg.Version)
}

// isInstalled checks if there is a local package installed.
func (cmd *command) isInstalled() bool {
	path := spec.Path(workDir, cmd.pkg.Owner, cmd.pkg.Name)
	return fileio.Exists(path)
}

// hasNewVersion checks if the remote package is newer than the local one.
func (cmd *command) hasNewVersion() bool {
	if cmd.err != nil {
		return true
	}

	oldPath := spec.Path(workDir, cmd.pkg.Owner, cmd.pkg.Name)
	if !fileio.Exists(oldPath) {
		return true
	}

	oldPkg, err := spec.ReadLocal(oldPath)
	if err != nil {
		cmd.err = err
		return true
	}
	debug("local package version = %s", oldPkg.Version)

	if oldPkg.Version == "" {
		// not explicitly versioned, always assume there is a later version
		return true
	}

	if oldPkg.Version == cmd.pkg.Version {
		return false
	}

	if semver.Compare(oldPkg.Version, cmd.pkg.Version) < 0 {
		return false
	}

	return true
}

// buildAssetPath constructs an URL to download package asset.
func (cmd *command) buildAssetPath() *spec.AssetPath {
	if cmd.err != nil {
		return nil
	}
	debug("checking remote asset for platform %s-%s", runtime.GOOS, runtime.GOARCH)
	debug("asset base path = %s", cmd.pkg.Assets.Path)

	var err error
	assetPath, err := cmd.pkg.AssetPath(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		cmd.err = fmt.Errorf("unsupported platform: %s-%s", runtime.GOOS, runtime.GOARCH)
		return nil
	}

	if !assetPath.Exists() {
		cmd.err = fmt.Errorf("asset does not exist: %s", assetPath)
		return nil
	}

	return assetPath
}

// downloadAsset downloads package asset.
func (cmd *command) downloadAsset(assetPath *spec.AssetPath) *assets.Asset {
	if cmd.err != nil {
		return nil
	}

	debug("downloading %s", assetPath)
	cmd.dir = spec.Dir(os.TempDir(), cmd.pkg.Owner, cmd.pkg.Name)
	err := fileio.CreateDir(cmd.dir)
	if err != nil {
		cmd.err = fmt.Errorf("failed to create temp directory: %w", err)
		return nil
	}

	var asset *assets.Asset
	if assetPath.IsRemote {
		asset, err = assets.Download(cmd.dir, assetPath.Value)
	} else {
		asset, err = assets.Copy(cmd.dir, assetPath.Value)
	}
	if err != nil {
		cmd.err = fmt.Errorf("failed to download asset: %w", err)
		return nil
	}

	sizeKb := float64(asset.Size) / 1024
	debug("downloaded %s (%.2f Kb)", asset.Name, sizeKb)
	return asset
}

// unpackAsset unpacks package asset.
func (cmd *command) unpackAsset(asset *assets.Asset) {
	if cmd.err != nil {
		return
	}

	assetPath := filepath.Join(cmd.dir, asset.Name)
	nFiles, err := assets.Unpack(assetPath, cmd.pkg.Assets.Pattern)
	if err != nil {
		cmd.err = fmt.Errorf("failed to unpack asset: %w", err)
		return
	}
	if nFiles == 0 {
		debug("not an archive, skipping unpack: %s", asset.Name)
		return
	}
	err = os.Remove(assetPath)
	if err != nil {
		cmd.err = fmt.Errorf("failed to delete asset after unpacking: %w", err)
	}
	debug("unpacked %d files from %s", nFiles, asset.Name)
}

// installFiles installes unpacked package files.
func (cmd *command) installFiles() {
	if cmd.err != nil {
		return
	}

	pkgDir := spec.Dir(workDir, cmd.pkg.Owner, cmd.pkg.Name)
	err := fileio.MoveDir(cmd.dir, pkgDir)
	if err != nil {
		cmd.err = fmt.Errorf("failed to copy downloaded files: %w", err)
		return
	}

	err = cmd.pkg.Save(pkgDir)
	if err != nil {
		cmd.err = fmt.Errorf("failed to write package spec: %w", err)
		return
	}

	cmd.dir = pkgDir
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
