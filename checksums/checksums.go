// Package checksums loads asset checksums from a file.
package checksums

import (
	"errors"
	"os"
	"strings"

	"sqlpkg.org/cli/fileio"
	"sqlpkg.org/cli/httpx"
)

// FileName is the checksum filename.
const FileName = "checksums.txt"

var ErrInvalidFile = errors.New("invalid checksum file")
var ErrInvalidSum = errors.New("invalid checksum value")

// Exists checks if a checksum file exists at the given path.
func Exists(path string, isRemote bool) bool {
	if isRemote {
		return httpx.Exists(path)
	} else {
		return fileio.Exists(path)
	}
}

// Read loads asset checksums from a local or remote file into a map,
// where keys are filenames and values are checksums.
func Read(path string, isRemote bool) (map[string]string, error) {
	read := inferReader(isRemote)
	data, err := read(path)
	if err != nil {
		return nil, err
	}
	return parse(data)
}

// A readFunc if a function that reads a file from a given path.
type readFunc func(path string) ([]byte, error)

// inferReader returns a proper reader function for a path,
// which can be a local file path or a remote url path.
func inferReader(isRemote bool) readFunc {
	if isRemote {
		return httpx.GetBytes
	} else {
		return os.ReadFile
	}
}

// parse parses checksum data into a map,
// where keys are filenames and values are checksums.
//
// Expects checksum data in the following format:
// 5072e5737...(sha-256 checksum)  sqlean-linux-x86.zip
// f86f443ac...(sha-256 checksum)  sqlean-macos-arm64.zip
// 8c0dc4fde...(sha-256 checksum)  sqlean-macos-x86.zip
// 0eead5873...(sha-256 checksum)  sqlean-win-x64.zip
func parse(data []byte) (map[string]string, error) {
	lines := strings.Split(string(data), "\n")
	sums := make(map[string]string, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			// want `checksum filename` line format
			return nil, ErrInvalidFile
		}
		if len(parts[0]) != 64 {
			// want sha-256 checksum
			return nil, ErrInvalidSum
		}
		sums[parts[1]] = "sha256-" + parts[0]
	}
	return sums, nil
}
