// Package spec manages the package spec file (sqlpkg.json).
package spec

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/nalgeon/sqlpkg-cli/internal/fileio"
	"github.com/nalgeon/sqlpkg-cli/internal/httpx"
)

// e.g. github.com/nalgeon/sqlean
var reGithub = regexp.MustCompile(`^github.com/[\w\-_.]+/[\w\-_.]+$`)

// e.g. nalgeon/sqlean
var reOwnerName = regexp.MustCompile(`^[\w\-_.]+/[\w\-_.]+$`)

// Read retrieves package spec file from the specified path.
// Path can be one of the following:
//   - owner-name pair: nalgeon/sqlean
//   - github repo: github.com/nalgeon/sqlean
//   - custom url: https://antonz.org/stuff/whatever/sqlean.json
//   - local path: /Users/anton/Desktop/sqlean.json
func Read(path string) (pkg *Package, err error) {
	errs := []error{}
	paths := expandPath(path)
	for _, path := range paths {
		readFunc := inferReader(path)
		pkg, err = readFunc(path)
		if err == nil {
			pkg.Specfile = path
			return pkg, nil
		} else {
			errs = append(errs, fmt.Errorf("%s: %w", path, err))
		}
	}
	return pkg, errors.Join(errs...)
}

// expandPath generates possible paths to the package spec file.
func expandPath(path string) []string {
	if reGithub.MatchString(path) {
		// try reading from the main branch of the github repository
		return []string{fmt.Sprintf("https://%s/raw/main/%s", path, FileName)}
	}
	if reOwnerName.MatchString(path) {
		// can be a local path or an owner-name pair, which in turn can point
		// to the author's repo or to the sqlpkg's registry
		return []string{
			path,
			fmt.Sprintf("https://github.com/%s/raw/main/%s", path, FileName),
			fmt.Sprintf("https://github.com/nalgeon/sqlpkg/raw/main/pkg/%s.json", path),
		}
	}
	return []string{path}
}

// ReadLocal reads package spec from a local file.
func ReadLocal(path string) (pkg *Package, err error) {
	return fileio.ReadJSON[Package](path)
}

// ReadRemote reads package spec from a remote url.
func ReadRemote(path string) (pkg *Package, err error) {
	return httpx.GetJSON[Package](path)
}

// A ReadFunc if a function that reads package spec from a given path.
type ReadFunc func(path string) (*Package, error)

// inferReader returns a proper reader function for a path,
// which can be a local file path or a remote url path.
func inferReader(path string) ReadFunc {
	if httpx.IsURL(path) {
		return ReadRemote
	} else {
		return ReadLocal
	}
}
