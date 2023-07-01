package cmd

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"
)

const helpHelp = "usage: sqlpkg help"

var commandsHelp = map[string]string{
	"init":      "Create a local repository",
	"install":   "Install a package",
	"uninstall": "Uninstall a package",
	"update":    "Update installed packages",
	"list":      "List installed packages",
	"info":      "Display package information",
}

// Help prints available commands.
func Help(args []string) error {
	if len(args) != 0 {
		return errors.New(helpHelp)
	}

	log("`sqlpkg` is an SQLite package manager. Use it to install or update SQLite extensions.")
	log("Commands:")
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 0, ' ', 0)
	for cmd, descr := range commandsHelp {
		fmt.Fprintln(w, cmd, "\t", descr)
	}
	w.Flush()

	return nil
}
