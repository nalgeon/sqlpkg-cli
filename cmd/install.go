package cmd

import (
	"errors"
)

const installHelp = "usage: sqlpkg install package"

// Install installs a new package or updates an existing one.
func Install(args []string) error {
	if len(args) != 1 {
		return errors.New(installHelp)
	}

	path := args[0]
	log("> installing %s...", path)

	cmd := new(command)
	cmd.readSpec(path)
	if !cmd.hasNewVersion() {
		log("✓ already at the latest version")
		return nil
	}
	assetPath := cmd.buildAssetPath()
	asset := cmd.downloadAsset(assetPath)
	cmd.unpackAsset(asset)
	cmd.installFiles()
	if cmd.err != nil {
		return cmd.err
	}

	log("✓ installed package %s to %s", cmd.pkg.FullName(), cmd.dir)
	return cmd.err
}
