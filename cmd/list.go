package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"text/tabwriter"

	"github.com/nalgeon/sqlpkg-cli/internal/lockfile"
	"github.com/nalgeon/sqlpkg-cli/internal/spec"
)

const listHelp = "usage: sqlpkg list"

// List prints all installed packages.
func List(args []string) error {
	if len(args) != 0 {
		return errors.New(listHelp)
	}

	lck, err := readLockfile()
	if err != nil {
		return err
	}
	debug("found %d packages in lockfile", len(lck.Packages))

	addedCount, err := addInstalledPackages(lck)
	if err != nil {
		return err
	}
	if addedCount > 0 {
		debug("found %d installed packages not in lockfile", addedCount)
		lck.Save(workDir)
	}

	packages, err := gatherLockfilePackages(lck)
	if err != nil {
		return err
	}

	sort.Slice(packages, func(i, j int) bool {
		return packages[i].FullName() < packages[j].FullName()
	})
	printPackages(packages)
	return nil
}

// addInstalledPackages checks for installed packages
// missing from the lockfile and adds them there.
func addInstalledPackages(lck *lockfile.Lockfile) (int, error) {
	pattern := fmt.Sprintf("%s/%s/*/*/%s", workDir, spec.DirName, spec.FileName)
	paths, _ := filepath.Glob(pattern)

	count := 0
	for _, path := range paths {
		pkg, err := spec.ReadLocal(path)
		if err != nil {
			return 0, fmt.Errorf("invalid package spec: %s", path)
		}
		if !lck.Has(pkg.FullName()) {
			lck.Add(pkg)
			count += 1
		}
	}
	return count, nil
}

// gatherLockfilePackages returns packages listed in the lockfile.
func gatherLockfilePackages(lck *lockfile.Lockfile) ([]*spec.Package, error) {
	packages := make([]*spec.Package, 0, len(lck.Packages))
	for fullName, pkg := range lck.Packages {
		if !isInstalled(pkg) {
			err := fmt.Errorf("package %s listed in specfile but not installed", fullName)
			return nil, err
		}
		packages = append(packages, pkg)
	}
	return packages, nil
}

// printPackages prints packages.
func printPackages(packages []*spec.Package) {
	if workDir == "." {
		log("(local repository)")
	}
	if len(packages) == 0 {
		log("no packages installed")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 0, ' ', 0)
	defer w.Flush()

	for _, pkg := range packages {
		fmt.Fprintln(w, pkg.FullName(), "\t", pkg.Description)
	}
}
