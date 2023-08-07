package main

import (
	"os"
	"reflect"
	"testing"
)

func Test_parseArgs(t *testing.T) {
	tests := []struct {
		in   []string
		cmd  string
		args []string
	}{
		{
			[]string{"sqlpkg"},
			"", nil,
		},
		{
			[]string{"sqlpkg", "help"},
			"help", []string{},
		},
		{
			[]string{"sqlpkg", "install", "nalgeon/example"},
			"install", []string{"nalgeon/example"},
		},
		{
			[]string{"sqlpkg", "-v", "install", "nalgeon/example"},
			"install", []string{"nalgeon/example"},
		},
	}
	for _, test := range tests {
		os.Args = test.in
		cmd, args := parseArgs()
		if cmd != test.cmd {
			t.Errorf("parseArgs(%v) expected cmd = %s, got %s", test.in, test.cmd, cmd)
		}
		if !reflect.DeepEqual(args, test.args) {
			t.Errorf("parseArgs(%v) expected args = %v, got %v", test.in, test.args, args)
		}
	}
}
