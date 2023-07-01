package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/nalgeon/sqlpkg-cli/cmd"
)

func execCommand() error {
	if len(os.Args) < 2 {
		return cmd.Help(nil)
	}

	flag.BoolVar(&cmd.IsVerbose, "v", false, "verbose output")
	flag.Parse()

	args := flag.Args()
	command, args := args[0], args[1:]

	switch command {
	case "init":
		return cmd.Init(args)
	case "install":
		return cmd.Install(args)
	case "uninstall":
		return cmd.Uninstall(args)
	case "update":
		if len(args) == 1 {
			return cmd.Update(args)
		}
		return cmd.UpdateAll(args)
	case "list":
		return cmd.List(args)
	case "info":
		return cmd.Info(args)
	case "help":
		return cmd.Help(args)
	default:
		return errors.New("unknown command")
	}
}

func main() {
	err := execCommand()
	if err != nil {
		fmt.Println("!", err)
		os.Exit(1)
	}
}
