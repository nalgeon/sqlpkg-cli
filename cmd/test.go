// Test helpers.
package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func SetupTestRepo(t *testing.T) (string, string) {
	repoDir := filepath.Join(WorkDir, ".sqlpkg")
	err := os.RemoveAll(repoDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lockPath := filepath.Join(WorkDir, "sqlpkg.lock")
	err = os.RemoveAll(lockPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return repoDir, lockPath
}

func TeardownTestRepo(t *testing.T, repoDir, lockPath string) {
	err := os.RemoveAll(repoDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = os.RemoveAll(lockPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
