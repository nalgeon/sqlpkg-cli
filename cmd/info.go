package cmd

import (
	"errors"
	"strings"

	"github.com/nalgeon/sqlpkg-cli/internal/spec"
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
		return err
	}

	if pkg.Description != "" {
		log(pkg.Description)
	}
	if pkg.Repository != "" {
		log(pkg.Repository)
	}
	if len(pkg.Authors) != 0 {
		authors := strings.Join(pkg.Authors, ", ")
		log("by %s", authors)
	}
	if pkg.Version != "" {
		log("version: %s", pkg.Version)
	}
	if pkg.License != "" {
		log("license: %s", pkg.License)
	}
	if isInstalled(pkg) {
		log("✓ installed")
	} else {
		log("✘ not installed")
	}

	return nil
}

// findSpec loads the package spec, giving preference to already installed packages.
func findSpec(path string) (*spec.Package, error) {
	installPath, err := getPathByFullName(path)
	if err == nil {
		pkg, err := spec.ReadLocal(installPath)
		if err == nil {
			debug("found installed package")
			return pkg, nil
		}
	}

	debug("installed package not found")
	pkg, err := readSpec(path)
	return pkg, err
}
