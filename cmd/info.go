package cmd

import (
	"errors"
	"strings"

	"sqlpkg.org/cli/spec"
)

const infoHelp = "usage: sqlpkg info package"

// Info prints information about the package (installed or not).
func Info(args []string) error {
	if len(args) != 1 {
		return errors.New(infoHelp)
	}

	path := args[0]
	pkg, err := findSpec(path)
	if err != nil {
		debug(err.Error())
		log("package not found")
		return nil
	}

	lines := prepareInfo(pkg)
	log(strings.Join(lines, "\n"))

	return nil
}

// findSpec loads the package spec, giving preference to already installed packages.
func findSpec(path string) (*spec.Package, error) {
	pkg := readInstalledSpec(path)
	if pkg != nil {
		return pkg, nil
	}

	debug("installed package not found")
	pkg, err := readSpec(path)
	return pkg, err
}

// readInstalledSpec loads the package spec for an installed package (if any).
func readInstalledSpec(fullName string) *spec.Package {
	path, err := getPathByFullName(fullName)
	if err != nil {
		return nil
	}

	pkg, err := spec.ReadLocal(path)
	if err != nil {
		return nil
	}

	debug("found installed package")
	return pkg
}

// prepareInfo returns detailed package description.
func prepareInfo(pkg *spec.Package) []string {
	lines := []string{}

	header := pkg.FullName()
	if pkg.Version != "" {
		header += "@" + pkg.Version
	}
	if len(pkg.Authors) != 0 {
		authors := strings.Join(pkg.Authors, ", ")
		header += " by " + authors
	}
	lines = append(lines, header)

	if pkg.Description != "" {
		lines = append(lines, pkg.Description)
	}
	if pkg.Repository != "" {
		lines = append(lines, pkg.Repository)
	}
	if pkg.License != "" {
		lines = append(lines, "license: "+pkg.License)
	}
	if isInstalled(pkg) {
		lines = append(lines, "✓ installed")
	} else {
		lines = append(lines, "✘ not installed")
	}
	return lines
}
