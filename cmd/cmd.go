// Package cmd implements sqlpkg commands logic.
package cmd

import (
	"errors"
	"os"
	"strings"

	"sqlpkg.org/cli/fileio"
	"sqlpkg.org/cli/logx"
	"sqlpkg.org/cli/spec"
)

// WorkDir is the current working directory.
var WorkDir string

// userHomeDir is the user's home directory.
var userHomeDir string

func init() {
	inferWorkDir()
}

// GetDirByFullName expands an owner-name package pair to a full package dir.
func GetDirByFullName(fullName string) (string, error) {
	parts := strings.Split(fullName, "/")
	if len(parts) != 2 {
		return "", errors.New("invalid package name")
	}
	path := spec.Dir(WorkDir, parts[0], parts[1])
	return path, nil
}

// GetPathByFullName expands an owner-name package pair to a full sqlpkg.json path.
func GetPathByFullName(fullName string) (string, error) {
	parts := strings.Split(fullName, "/")
	if len(parts) != 2 {
		return "", errors.New("invalid package name")
	}
	path := spec.Path(WorkDir, parts[0], parts[1])
	return path, nil
}

// PrintScope prints information about the current scope (project/global).
func PrintScope() {
	if WorkDir == "." {
		logx.Log("(project scope)")
	}
}

// inferWorkDir determines the working directory.
// It is either the .sqlpkg directory (if present) or ~/.sqlpkg otherwise.
func inferWorkDir() {
	if fileio.Exists(spec.DirName) {
		WorkDir = "."
		return
	}
	var err error
	userHomeDir, err = os.UserHomeDir()
	if err != nil {
		WorkDir = "."
		return
	}
	WorkDir = userHomeDir
}
