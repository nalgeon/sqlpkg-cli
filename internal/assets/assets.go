// Package assets manages package assets (hmm).
package assets

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/nalgeon/sqlpkg-cli/internal/fileio"
	"github.com/nalgeon/sqlpkg-cli/internal/httpx"
)

// An Asset is an archive of package files for a specific platform.
type Asset struct {
	Name     string
	Size     int64
	Checksum []byte
}

// Download downloads an asset from the remote url to the local dir.
func Download(dir, rawURL string) (asset *Asset, err error) {
	url, err := url.Parse(rawURL)
	if err != nil {
		return nil, errors.New("invalid url")
	}

	name := filepath.Base(url.Path)
	file, err := os.Create(filepath.Join(dir, name))
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
	return &Asset{name, size, nil}, nil
}

// Copy copies an asset from the local path to the local dir.
func Copy(dir, path string) (asset *Asset, err error) {
	_, name := filepath.Split(path)
	dstPath := filepath.Join(dir, name)
	size, err := fileio.CopyFile(path, dstPath)
	if err != nil {
		return nil, err
	}
	return &Asset{name, int64(size), nil}, nil
}

// Unpack unpacks an asset from the given path to the same dir
// where the asset resides. If pattern is provided, unpacks
// only the files that match it.
func Unpack(path, pattern string) (int, error) {
	dir, _ := filepath.Split(path)
	if strings.HasSuffix(path, ".zip") {
		return unpackZip(path, pattern, dir)
	}
	if strings.HasSuffix(path, ".tar.gz") {
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
