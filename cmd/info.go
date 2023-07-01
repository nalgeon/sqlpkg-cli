package cmd

import (
	"errors"
	"strings"
)

const infoHelp = "usage: sqlpkg info package"

// Info prints information about the package (installed or not).
func Info(args []string) error {
	if len(args) != 1 {
		return errors.New(infoHelp)
	}

	path := args[0]

	cmd := new(command)
	cmd.readMetadata(path)
	if cmd.err != nil {
		return errors.New("package not found")
	}

	if cmd.pkg.Description != "" {
		log(cmd.pkg.Description)
	}
	if cmd.pkg.Repository != "" {
		log(cmd.pkg.Repository)
	}
	if len(cmd.pkg.Authors) != 0 {
		authors := strings.Join(cmd.pkg.Authors, ", ")
		log("by %s", authors)
	}
	if cmd.pkg.Version != "" {
		log("version: %s", cmd.pkg.Version)
	}
	if cmd.pkg.License != "" {
		log("license: %s", cmd.pkg.License)
	}
	if cmd.isInstalled() {
		log("✓ installed")
	} else {
		log("✘ not installed")
	}

	return nil
}
