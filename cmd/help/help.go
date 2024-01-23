package help

import (
	"errors"
	"fmt"
	"sort"
	"text/tabwriter"

	"sqlpkg.org/cli/logx"
)

const help = "usage: sqlpkg help"

var commandsHelp = map[string]string{
	"help":      "Display help",
	"info":      "Display package information",
	"init":      "Init project scope",
	"install":   "Install packages",
	"list":      "List installed packages",
	"uninstall": "Uninstall package",
	"update":    "Update installed packages",
	"version":   "Display version",
	"which":     "Display path to extension file",
}

// Help prints available commands.
func Help(args []string) error {
	if len(args) != 0 {
		return errors.New(help)
	}

	logx.Log("sqlpkg is a package manager for installing and updating SQLite extensions.\n")
	logx.Log("USAGE")
	logx.Log("  sqlpkg [global-options] <command> [arguments]\n")
	logx.Log("GLOBAL OPTIONS")
	logx.Log("  -v  verbose output\n")
	logx.Log("COMMANDS")

	w := tabwriter.NewWriter(logx.Output(), 0, 4, 0, ' ', 0)
	for _, cmd := range sortedCommands() {
		fmt.Fprintln(w, "  ", cmd, "\t", commandsHelp[cmd])
	}
	w.Flush()

	return nil
}

// sortedCommands returns a slice of all commands sorted alphabetically.
func sortedCommands() []string {
	list := make([]string, 0, len(commandsHelp))
	for cmd := range commandsHelp {
		list = append(list, cmd)
	}
	sort.Strings(list)
	return list
}
