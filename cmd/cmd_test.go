package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"sqlpkg.org/cli/logx"
	"sqlpkg.org/cli/spec"
)

func TestWorkDir(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("os.UserHomeDir: unexpected error %v", err)
	}

	t.Run("home dir", func(t *testing.T) {
		inferWorkDir()
		if WorkDir != home {
			t.Errorf("WorkDir: unexpected value %q", WorkDir)
		}
	})
	t.Run("project dir", func(t *testing.T) {
		err = os.Mkdir(spec.DirName, 0755)
		if err != nil {
			t.Fatalf("os.Mkdir: unexpected error %v", err)
		}
		defer os.Remove(spec.DirName)

		inferWorkDir()
		if WorkDir != "." {
			t.Errorf("WorkDir: unexpected value %q", WorkDir)
		}
	})
}

func TestGetPathByFullName(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		path, err := GetPathByFullName("nalgeon/example")
		if err != nil {
			t.Fatalf("GetPathByFullName: unexpected error %v", err)
		}
		if path != filepath.Join(WorkDir, spec.DirName, "nalgeon", "example", spec.FileName) {
			t.Errorf("GetPathByFullName: unexpected value %q", path)
		}
	})
	t.Run("invalid", func(t *testing.T) {
		_, err := GetPathByFullName("nalgeon")
		if err == nil {
			t.Fatal("GetPathByFullName: expected error, fot nil")
		}
	})
}

func TestGetDirByFullName(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		dir, err := GetDirByFullName("nalgeon/example")
		if err != nil {
			t.Fatalf("GetDirByFullName: unexpected error %v", err)
		}
		if dir != filepath.Join(WorkDir, spec.DirName, "nalgeon", "example") {
			t.Errorf("GetDirByFullName: unexpected value %q", dir)
		}
	})
	t.Run("invalid", func(t *testing.T) {
		_, err := GetDirByFullName("nalgeon")
		if err == nil {
			t.Fatal("GetDirByFullName: expected error, fot nil")
		}
	})
}

func TestPrintLocalRepo(t *testing.T) {
	mem := logx.Mock()
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("os.UserHomeDir: unexpected error %v", err)
	}

	t.Run("home dir", func(t *testing.T) {
		WorkDir = home
		PrintLocalRepo()
		if len(mem.Lines) != 0 {
			t.Fatalf("PrintLocalRepo: unexpected line count %v", len(mem.Lines))
		}
	})
	t.Run("project dir", func(t *testing.T) {
		WorkDir = "."
		PrintLocalRepo()
		if len(mem.Lines) != 1 {
			t.Fatalf("PrintLocalRepo: unexpected line count %v", len(mem.Lines))
		}
		if !mem.Has("(local repository)") {
			t.Errorf("PrintLocalRepo: unexpected output %v", mem.Lines)
		}
	})
}
