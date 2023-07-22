// Test helpers.
package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"sqlpkg.org/cli/logx"
)

func SetupTestRepo(t *testing.T) (string, string) {
	repoDir := filepath.Join(WorkDir, ".sqlpkg")
	err := os.RemoveAll(repoDir)
	if err != nil {
		t.Fatalf("SetupTestRepo: %v", err)
	}
	lockPath := filepath.Join(WorkDir, "sqlpkg.lock")
	err = os.RemoveAll(lockPath)
	if err != nil {
		t.Fatalf("SetupTestRepo: %v", err)
	}
	return repoDir, lockPath
}

func CopyTestRepo(t *testing.T) {
	{
		cmd := []string{"cp", "-r", "./testdata/.sqlpkg", "."}
		err := exec.Command(cmd[0], cmd[1:]...).Run()
		if err != nil {
			t.Fatalf("%s: %v", strings.Join(cmd, " "), err)
		}
	}
	{
		cmd := []string{"cp", "./testdata/sqlpkg.lock", "."}
		err := exec.Command(cmd[0], cmd[1:]...).Run()
		if err != nil {
			t.Fatalf("%s: %v", strings.Join(cmd, " "), err)
		}
	}
}

func TeardownTestRepo(t *testing.T, repoDir, lockPath string) {
	err := os.RemoveAll(repoDir)
	if err != nil {
		t.Fatalf("TeardownTestRepo: %v", err)
	}
	err = os.RemoveAll(lockPath)
	if err != nil {
		t.Fatalf("TeardownTestRepo: %v", err)
	}
}

func SetupTestLogger() *logx.Memory {
	memory := logx.NewMemory("log")
	logx.SetOutput(memory)
	logx.SetVerbose(true)
	return memory
}
