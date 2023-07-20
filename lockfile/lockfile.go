// Package lockfile manages the lockfile (sqlpkg.lock).
package lockfile

import (
	"encoding/json"
	"os"
	"path/filepath"

	"sqlpkg.org/cli/spec"
)

// FileName is the lockfile filename.
const FileName = "sqlpkg.lock"

// Path returns the path to the lockfile.
func Path(basePath string) string {
	return filepath.Join(basePath, FileName)
}

// A Lockfile describes a collection of specific package versions.
type Lockfile struct {
	Packages map[string]*spec.Package `json:"packages"`
}

// NewLockfile creates an empty lockfile.
func NewLockfile() *Lockfile {
	packages := map[string]*spec.Package{}
	return &Lockfile{packages}
}

// Has checks if a package is in the lockfile.
func (lck *Lockfile) Has(fullName string) bool {
	_, ok := lck.Packages[fullName]
	return ok
}

// Range iterates over packages from the lockfile.
func (lck *Lockfile) Range(fn func(fullName string, pkg *spec.Package) bool) {
	for fullName, pkg := range lck.Packages {
		ok := fn(fullName, pkg)
		if !ok {
			break
		}
	}
}

// Add adds a package to the lockfile.
func (lck *Lockfile) Add(pkg *spec.Package) {
	p := spec.Package{
		Owner:    pkg.Owner,
		Name:     pkg.Name,
		Version:  pkg.Version,
		Specfile: pkg.Specfile,
		Assets:   pkg.Assets,
	}
	lck.Packages[pkg.FullName()] = &p
}

// Remove removes a package from the lockfile.
func (lck *Lockfile) Remove(pkg *spec.Package) {
	delete(lck.Packages, pkg.FullName())
}

// Save writes the lockfile to the specified directory.
func (lck *Lockfile) Save(dir string) error {
	data, err := json.MarshalIndent(lck, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(Path(dir), data, 0644)
}
