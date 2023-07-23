// Commands that manage the package version.
package cmd

import (
	"fmt"

	"golang.org/x/mod/semver"
	"sqlpkg.org/cli/fileio"
	"sqlpkg.org/cli/github"
	"sqlpkg.org/cli/httpx"
	"sqlpkg.org/cli/logx"
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

// HasNewVersion checks if the remote package is newer than the local one.
func HasNewVersion(pkg *spec.Package) bool {
	oldPath := spec.Path(WorkDir, pkg.Owner, pkg.Name)
	if !fileio.Exists(oldPath) {
		return true
	}

	oldPkg, err := spec.ReadLocal(oldPath)
	if err != nil {
		return true
	}
	logx.Debug("local package version = %s", oldPkg.Version)

	if oldPkg.Version == "" {
		// not explicitly versioned, always assume there is a later version
		return true
	}

	if oldPkg.Version == pkg.Version {
		return false
	}

	return compareVersions(oldPkg.Version, pkg.Version) < 0
}

// compareVersions compares package versions.
// Returns 0 if v == w, -1 if v < w, or +1 if v > w.
func compareVersions(v, w string) int {
	if v == "" || w == "" {
		return 0
	}
	// add the leading 'v' if needed
	if v[0] != 'v' {
		v = "v" + v
	}
	if w[0] != 'v' {
		w = "v" + w
	}
	return semver.Compare(v, w)
}
