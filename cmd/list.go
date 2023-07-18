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

	packages, err := gatherPackages()
	if err != nil {
		return err
	}

	lck, err := readLockfile()
	if err != nil {
		return err
	}

	err = addMissingToLockfile(lck, packages)
	if err != nil {
		return err
	}

	sortPackages(packages)
	printPackages(packages)
	return nil
}

// gatherPackages collects installed packages.
func gatherPackages() ([]*spec.Package, error) {
	pattern := filepath.Join(workDir, spec.DirName, "*", "*", spec.FileName)
	paths, _ := filepath.Glob(pattern)

	packages := []*spec.Package{}
	for _, path := range paths {
		pkg, err := spec.ReadLocal(path)
		if err != nil {
			return nil, fmt.Errorf("invalid package spec: %s", path)
		}
		packages = append(packages, pkg)
	}

	debug("gathered %d packages", len(packages))
	return packages, nil
}

// addMissingToLockfile adds missing packages to the lockfile.
func addMissingToLockfile(lck *lockfile.Lockfile, packages []*spec.Package) error {
	count := 0
	for _, pkg := range packages {
		if lck.Has(pkg.FullName()) {
			continue
		}
		lck.Add(pkg)
		count += 1
	}

	if count == 0 {
		return nil
	}

	err := lck.Save(workDir)
	if err != nil {
		return fmt.Errorf("failed to save lockfile: %w", err)
	}

	debug("added %d packages to the lockfile", count)
	return nil
}

// sortPackages sorts packages by full name.
func sortPackages(packages []*spec.Package) {
	sort.Slice(packages, func(i, j int) bool {
		return packages[i].FullName() < packages[j].FullName()
	})
}

// printPackages prints packages.
func printPackages(packages []*spec.Package) {
	printLocalRepo()
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

func printLocalRepo() {
	if workDir == "." {
		log("(local repository)")
	}
}
