// Commands that manage the package version.
package cmd

import (
	"fmt"

	"sqlpkg.org/cli/fileio"
	"sqlpkg.org/cli/github"
	"sqlpkg.org/cli/httpx"
	"sqlpkg.org/cli/logx"
	"sqlpkg.org/cli/semver"
	"sqlpkg.org/cli/spec"
)

// ResolveVersion resolves the latest version if needed.
func ResolveVersion(pkg *spec.Package) error {
	if pkg.Version != "latest" {
		return nil
	}

	hostname := httpx.Hostname(pkg.Repository)
	if hostname != github.Hostname {
		logx.Debug("unknown provider %s, not resolving version", hostname)
		return nil
	}

	owner, repo, err := github.ParseRepoUrl(pkg.Repository)
	if err != nil {
		return fmt.Errorf("failed to parse repo url: %v", err)
	}

	version, err := github.GetLatestTag(owner, repo)
	if err != nil {
		return fmt.Errorf("failed to get latest tag: %w", err)
	}

	pkg.ReplaceLatest(version)
	logx.Debug("resolved latest version = %s", version)
	return nil
}

// HasNewVersion checks if the remote package is newer than the installed one.
func HasNewVersion(remotePkg *spec.Package) bool {
	installPath := spec.Path(WorkDir, remotePkg.Owner, remotePkg.Name)
	if !fileio.Exists(installPath) {
		return true
	}

	installedPkg, err := spec.ReadLocal(installPath)
	if err != nil {
		return true
	}
	logx.Debug("local package version = %s", installedPkg.Version)

	if installedPkg.Version == "" {
		// not explicitly versioned, always assume there is a later version
		return true
	}

	if installedPkg.Version == remotePkg.Version {
		return false
	}

	return semver.Compare(installedPkg.Version, remotePkg.Version) < 0
}
