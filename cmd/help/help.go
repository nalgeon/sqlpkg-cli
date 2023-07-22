package help

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"sqlpkg.org/cli/logx"
)

const helpHelp = "usage: sqlpkg help"

var commands = []string{
	"install", "uninstall", "update", "list", "init", "info", "help", "version",
}

var commandsHelp = map[string]string{
	"help":      "Display help",
	"info":      "Display package information",
	"init":      "Create a local repository",
	"install":   "Install packages",
	"list":      "List installed packages",
	"uninstall": "Uninstall a package",
	"update":    "Update installed packages",
	"version":   "Display version",
}

// Help prints available commands.
func Help(args []string) error {
	if len(args) != 0 {
		return errors.New(helpHelp)
	}

	logx.Log("`sqlpkg` is an SQLite package manager. Use it to install or update SQLite extensions.\n")
	logx.Log("Commands:")
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 0, ' ', 0)
	for _, cmd := range commands {
		fmt.Fprintln(w, cmd, "\t", commandsHelp[cmd])
	}
	w.Flush()

	return nil
}
