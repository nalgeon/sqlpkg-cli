package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/nalgeon/sqlpkg-cli/internal/metadata"
)

const listHelp = "usage: sqlpkg list"

// List prints all installed packages.
func List(args []string) error {
	if len(args) != 0 {
		return errors.New(listHelp)
	}

	pattern := fmt.Sprintf("%s/%s/*/*/%s", workDir, metadata.DirName, metadata.FileName)
	paths, _ := filepath.Glob(pattern)

	if len(paths) == 0 {
		log("no packages installed")
		return nil
	}

	if workDir == "." {
		log("(local repository)")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 0, ' ', 0)
	for _, path := range paths {
		pkg, err := metadata.ReadLocal(path)
		if err != nil {
			return fmt.Errorf("invalid package spec: %s", path)
		}
		fmt.Fprintln(w, pkg.FullName(), "\t", pkg.Description)
	}
	w.Flush()

	return nil
}
