package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"sqlpkg.org/cli/cmd/help"
	"sqlpkg.org/cli/cmd/info"
	init_ "sqlpkg.org/cli/cmd/init"
	"sqlpkg.org/cli/cmd/install"
	"sqlpkg.org/cli/cmd/list"
	"sqlpkg.org/cli/cmd/uninstall"
	"sqlpkg.org/cli/cmd/update"
	"sqlpkg.org/cli/cmd/which"
	"sqlpkg.org/cli/logx"
)

var version = "main"

func parseArgs() (command string, args []string) {
	if len(os.Args) < 2 {
		return "", nil
	}

	var isVerbose bool
	flag.BoolVar(&isVerbose, "v", false, "verbose output")
	flag.Parse()

	logx.SetVerbose(isVerbose)
	args = flag.Args()
	command, args = args[0], args[1:]
	return
}

func execCommand(command string, args []string) error {
	if command == "" {
		return help.Help(nil)
	}

	switch command {
	case "init":
		return init_.Init(args)
	case "install":
		if len(args) == 0 {
			return install.InstallAll(args)
		}
		return install.Install(args)
	case "uninstall":
		return uninstall.Uninstall(args)
	case "update":
		if len(args) == 0 {
			return update.UpdateAll(args)
		}
		return update.Update(args)
	case "list":
		return list.List(args)
	case "info":
		return info.Info(args)
	case "which":
		return which.Which(args)
	case "help":
		return help.Help(args)
	case "version":
		fmt.Println(version)
		return nil
	default:
		return errors.New("unknown command")
	}
}

func main() {
	command, args := parseArgs()
	err := execCommand(command, args)
	if err != nil {
		fmt.Println("!", err)
		os.Exit(1)
	}
}
