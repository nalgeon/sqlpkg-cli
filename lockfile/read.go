package lockfile

import (
	"sqlpkg.org/cli/fileio"
)

// ReadLocal reads the lockfile from a local file.
func ReadLocal(path string) (lck *Lockfile, err error) {
	return fileio.ReadJSON[Lockfile](path)
}
