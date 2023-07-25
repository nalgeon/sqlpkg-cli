package which

import (
	"errors"
	"path/filepath"
	"runtime"
	"strings"

	"sqlpkg.org/cli/cmd"
	"sqlpkg.org/cli/fileio"
	"sqlpkg.org/cli/logx"
	"sqlpkg.org/cli/spec"
)

const help = "usage: sqlpkg which package"

// maps the OS name to the file extension
var fileExt = map[string]string{
	"darwin":  ".dylib",
	"linux":   ".so",
	"windows": ".dll",
}

// Which prints a path to the extension file.
func Which(args []string) error {
	if len(args) != 1 {
		return errors.New(help)
	}

	parts := strings.Split(args[0], "/")
	if len(parts) != 2 {
		return errors.New("invalid package name")
	}

	owner, name := parts[0], parts[1]
	pkgDir := spec.Dir(cmd.WorkDir, owner, name)
	if !fileio.Exists(pkgDir) {
		return errors.New("package is not installed")
	}

	path := findExact(pkgDir, name, runtime.GOOS)
	if path != "" {
		logx.Log(path)
		return nil
	}

	paths := findByExt(pkgDir, name, runtime.GOOS)
	if len(paths) == 0 {
		return errors.New("extension file is not found")
	}

	logx.Log("exact match not found")
	logx.Log("possible matches:")
	for _, path := range paths {
		logx.Log(path)
	}

	return nil
}

// findExact returns a path to the extension file
// if the extension file has the same name as the package itself.
func findExact(pkgDir, name, os string) string {
	{
		// e.g., text.dylib
		pattern := filepath.Join(pkgDir, name+fileExt[os])
		paths, _ := filepath.Glob(pattern)
		if len(paths) != 0 {
			return paths[0]
		}
	}
	{
		// e.g., text0.dylib
		pattern := filepath.Join(pkgDir, name+"[0-9]"+fileExt[os])
		paths, _ := filepath.Glob(pattern)
		if len(paths) != 0 {
			return paths[0]
		}
	}
	{
		// e.g., libtext.dylib
		pattern := filepath.Join(pkgDir, "lib"+name+fileExt[os])
		paths, _ := filepath.Glob(pattern)
		if len(paths) != 0 {
			return paths[0]
		}
	}
	// no exact match
	return ""
}

// findByExt returns paths to files in the package dir
// that have an expected extension (e.g. textext.dylib)
func findByExt(pkgDir, name, os string) []string {
	file := "*" + fileExt[os]
	pattern := filepath.Join(pkgDir, file)
	paths, _ := filepath.Glob(pattern)
	return paths
}
