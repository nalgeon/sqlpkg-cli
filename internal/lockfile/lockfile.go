// Package lockfile manages the lockfile (sqlpkg.lock).
package lockfile

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/nalgeon/sqlpkg-cli/internal/spec"
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

// Add adds a package to the lockfile.
func (lf *Lockfile) Add(pkg *spec.Package) {
	p := spec.Package{
		Owner:    pkg.Owner,
		Name:     pkg.Name,
		Version:  pkg.Version,
		Specfile: pkg.Specfile,
		Assets:   pkg.Assets,
	}
	lf.Packages[pkg.FullName()] = &p
}

// Remove removes a package from the lockfile.
func (lf *Lockfile) Remove(pkg *spec.Package) {
	delete(lf.Packages, pkg.FullName())
}

// Save writes the lockfile to the specified directory.
func (lf *Lockfile) Save(dir string) error {
	data, err := json.MarshalIndent(lf, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(Path(dir), data, 0644)
}
