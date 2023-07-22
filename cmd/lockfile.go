// Commands that manage the lockfile.
package cmd

import (
	"fmt"

	"sqlpkg.org/cli/fileio"
	"sqlpkg.org/cli/lockfile"
	"sqlpkg.org/cli/logx"
	"sqlpkg.org/cli/spec"
)

// ReadLockfile reads lockfile from the work directory.
func ReadLockfile() (*lockfile.Lockfile, error) {
	path := lockfile.Path(WorkDir)
	if !fileio.Exists(path) {
		logx.Debug("created new lockfile")
		return lockfile.NewLockfile(), nil
	}

	lck, err := lockfile.ReadLocal(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read lockfile: %w", err)
	}

	logx.Debug("read existing lockfile")
	return lck, nil
}

// AddToLockfile adds package to the lockfile.
func AddToLockfile(lck *lockfile.Lockfile, pkg *spec.Package) error {
	lck.Add(pkg)
	err := lck.Save(WorkDir)
	if err != nil {
		return fmt.Errorf("failed to save lockfile: %w", err)
	}

	logx.Debug("added package to the lockfile")
	return nil
}

// RemoveFromLockfile removes package from the lockfile.
func RemoveFromLockfile(lck *lockfile.Lockfile, fullName string) error {
	pkg, ok := lck.Packages[fullName]
	if !ok {
		logx.Debug("package not listed in the lockfile")
		return nil
	}

	lck.Remove(pkg)
	err := lck.Save(WorkDir)
	if err != nil {
		return fmt.Errorf("failed to save lockfile: %w", err)
	}

	logx.Debug("removed package from the lockfile")
	return nil
}
