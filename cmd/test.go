// Test helpers.
package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// SetupTestRepo returns paths to the .sqlpkg folder and the lockfile
// as if they reside in the working directory.
func SetupTestRepo(t *testing.T) (repoDir string, lockPath string) {
	WorkDir = "."
	repoDir = filepath.Join(WorkDir, ".sqlpkg")
	err := os.RemoveAll(repoDir)
	if err != nil {
		t.Fatalf("SetupTestRepo: %v", err)
	}
	lockPath = filepath.Join(WorkDir, "sqlpkg.lock")
	err = os.RemoveAll(lockPath)
	if err != nil {
		t.Fatalf("SetupTestRepo: %v", err)
	}
	return repoDir, lockPath
}

// CopyTestRepo copies the .sqlpkg folder and the lockfile
// from the testdata folder to the working directory.
func CopyTestRepo(t *testing.T, path ...string) {
	basePath := filepath.Join("testdata", filepath.Join(path...))

	{
		// copy `.sqlpkg` contents
		path := filepath.Join(basePath, ".sqlpkg")
		cmd := []string{"cp", "-r", path, "."}
		err := exec.Command(cmd[0], cmd[1:]...).Run()
		if err != nil {
			t.Fatalf("%s: %v", strings.Join(cmd, " "), err)
		}
	}
	{
		// copy lockfile
		path := filepath.Join(basePath, "sqlpkg.lock")
		cmd := []string{"cp", path, "."}
		err := exec.Command(cmd[0], cmd[1:]...).Run()
		if err != nil {
			t.Fatalf("%s: %v", strings.Join(cmd, " "), err)
		}
	}
}

// TeardownTestRepo deletes the .sqlpkg folder and the lockfile
// from the working directory.
func TeardownTestRepo(t *testing.T) {
	repoDir := filepath.Join(WorkDir, ".sqlpkg")
	err := os.RemoveAll(repoDir)
	if err != nil {
		t.Fatalf("TeardownTestRepo: %v", err)
	}

	lockPath := filepath.Join(WorkDir, "sqlpkg.lock")
	err = os.RemoveAll(lockPath)
	if err != nil {
		t.Fatalf("TeardownTestRepo: %v", err)
	}
}
