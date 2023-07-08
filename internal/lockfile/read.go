package lockfile

import (
	"github.com/nalgeon/sqlpkg-cli/internal/fileio"
)

// ReadLocal reads the lockfile from a local file.
func ReadLocal(path string) (lck *Lockfile, err error) {
	return fileio.ReadJSON[Lockfile](path)
}
