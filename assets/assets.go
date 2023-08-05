// Package assets manages package assets (hmm).
package assets

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/hex"
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"sqlpkg.org/cli/fileio"
	"sqlpkg.org/cli/httpx"
)

// An Asset is an archive of package files for a specific platform.
type Asset struct {
	Name     string
	Path     string
	Size     int64
	Checksum []byte
}

func (a *Asset) Dir() string {
	return filepath.Dir(a.Path)
}

// Validate compares the asset checksum against the provided checksum string.
func (a *Asset) Validate(checksumStr string) (bool, error) {
	algo, str, ok := strings.Cut(checksumStr, "-")
	if !ok || algo != "sha256" {
		return false, errors.New("unsupported checksum algorithm")
	}
	checksum, err := hex.DecodeString(str)
	if err != nil {
		return false, errors.New("failed to decode checksum string")
	}
	return areEqual(a.Checksum, checksum), nil
}

// Download downloads an asset from the remote url to the local dir.
func Download(dir, rawURL string) (asset *Asset, err error) {
	url, err := url.Parse(rawURL)
	if err != nil {
		return nil, errors.New("invalid url")
	}

	name := filepath.Base(url.Path)
	path := filepath.Join(dir, name)
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body, err := httpx.GetBody(rawURL, "application/octet-stream")
	if err != nil {
		return nil, err
	}
	defer body.Close()

	size, err := io.Copy(file, body)
	if err != nil {
		return nil, err
	}

	checksum, err := fileio.CalcChecksum(path)
	if err != nil {
		return nil, err
	}

	return &Asset{name, path, size, checksum}, nil
}

// Copy copies an asset from the local path to the local dir.
func Copy(dir, path string) (asset *Asset, err error) {
	_, name := filepath.Split(path)
	dstPath := filepath.Join(dir, name)

	size, err := fileio.CopyFile(path, dstPath)
	if err != nil {
		return nil, err
	}

	checksum, err := fileio.CalcChecksum(dstPath)
	if err != nil {
		return nil, err
	}

	return &Asset{name, dstPath, int64(size), checksum}, nil
}

// Unpack unpacks an asset from the given path to the same dir
// where the asset resides. If pattern is provided, unpacks
// only the files that match it. Returns the number of unpacked files.
func Unpack(path, pattern string) (int, error) {
	dir, _ := filepath.Split(path)
	if strings.HasSuffix(path, ".zip") {
		return unpackZip(path, pattern, dir)
	}
	if strings.HasSuffix(path, ".tar.gz") || strings.HasSuffix(path, ".tgz") {
		return unpackTarGz(path, pattern, dir)
	}
	return 0, nil
}

// unpackZip unpackes a zip archive.
func unpackZip(path, pattern, dir string) (int, error) {
	archive, err := zip.OpenReader(path)
	if err != nil {
		return 0, err
	}
	defer archive.Close()

	count := 0
	for _, f := range archive.File {
		if f.FileInfo().IsDir() {
			// ignore dirs
			continue
		}
		if pattern != "" {
			matched, _ := filepath.Match(pattern, f.Name)
			if !matched {
				continue
			}
		}

		dstPath := filepath.Join(dir, f.Name)
		dstFile, err := os.Create(dstPath)
		if err != nil {
			return 0, err
		}

		file, err := f.Open()
		if err != nil {
			return 0, err
		}

		_, err = io.Copy(dstFile, file)
		if err != nil {
			file.Close()
			return 0, err
		}
		dstFile.Close()
		count += 1
	}

	return count, nil
}

// unpackTarGz unpackes a .tar.gz archive.
func unpackTarGz(path, pattern, dir string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}

	gzip, err := gzip.NewReader(file)
	if err != nil {
		return 0, err
	}
	defer gzip.Close()

	rdr := tar.NewReader(gzip)

	count := 0
	for {
		header, err := rdr.Next()

		if err == io.EOF {
			return count, nil
		}
		if err != nil {
			return 0, err
		}

		if header.Typeflag != tar.TypeReg {
			// ignore dirs
			continue
		}

		if pattern != "" {
			matched, _ := filepath.Match(pattern, header.Name)
			if !matched {
				continue
			}
		}

		dstPath := filepath.Join(dir, header.Name)
		dstFile, err := os.Create(dstPath)
		if err != nil {
			return 0, err
		}

		// copy over contents
		_, err = io.Copy(dstFile, rdr)
		if err != nil {
			return 0, err
		}

		dstFile.Close()
		count += 1
	}
}

// areEqual checks if two slices are equal.
func areEqual[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i, v := range s1 {
		if v != s2[i] {
			return false
		}
	}

	return true
}
