package info

import (
	"errors"
	"strings"

	"sqlpkg.org/cli/cmd"
	"sqlpkg.org/cli/fileio"
	"sqlpkg.org/cli/spec"
)

const infoHelp = "usage: sqlpkg info package"

// Info prints information about the package (installed or not).
func Info(args []string) error {
	if len(args) != 1 {
		return errors.New(infoHelp)
	}

	path := args[0]
	pkg, err := cmd.FindSpec(path)
	if err != nil {
		cmd.Debug(err.Error())
		cmd.Log("package not found")
		return nil
	}

	lines := prepareInfo(pkg)
	cmd.Log(strings.Join(lines, "\n"))

	return nil
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

// isInstalled checks if there is a local package installed.
func isInstalled(pkg *spec.Package) bool {
	path := spec.Path(cmd.WorkDir, pkg.Owner, pkg.Name)
	return fileio.Exists(path)
}
