// Package fileio provides high-level file operations.
package fileio

import (
	"crypto/sha256"
	"io"
	"os"
)

// CreateDir creates an empty directory.
// If the directory already exists, deletes it and creates a new one.
func CreateDir(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	return nil
}

// MoveDir moves the source directory to the destination.
// If the destination already exists, deletes it before moving the source.
func MoveDir(src, dst string) error {
	err := os.MkdirAll(dst, 0755)
	if err != nil {
		return err
	}
	err = os.RemoveAll(dst)
	if err != nil {
		return err
	}
	err = os.Rename(src, dst)
	if err != nil {
		return err
	}
	return nil
}

// CopyFile copies a single file from source to destination.
func CopyFile(src, dst string) (int, error) {
	data, err := os.ReadFile(src)
	if err != nil {
		return 0, err
	}
	err = os.WriteFile(dst, data, 0644)
	return len(data), err
}

// Exists checks if the specified path exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// CalcChecksum calculates the SHA-256 checksum of a file.
func CalcChecksum(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}
