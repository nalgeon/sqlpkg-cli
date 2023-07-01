// Package metadata manages the package metadata (the sqlpkg.json file).
package metadata

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
)

// e.g. github.com/nalgeon/sqlean
var reGithub = regexp.MustCompile(`^github.com/[\w\-_.]+/[\w\-_.]+$`)

// e.g. nalgeon/sqlean
var reOwnerName = regexp.MustCompile(`^[\w\-_.]+/[\w\-_.]+$`)

// Read retrieves package metadata from the specified path.
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
			pkg.Path = path
			return pkg, nil
		} else {
			errs = append(errs, fmt.Errorf("%s: %w", path, err))
		}
	}
	return pkg, errors.Join(errs...)
}

// expandPath generates possible paths to the package metadata file.
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

// inferReader returns a proper reader function for a path,
// which can be a local file path or a remote url path.
func inferReader(path string) ReadFunc {
	if isURL(path) {
		return ReadRemote
	} else {
		return ReadLocal
	}
}

// isURL checks if the path is an url.
func isURL(path string) bool {
	u, err := url.Parse(path)
	if err != nil {
		return false
	}
	if u.Scheme == "" {
		return false
	}
	return true
}
