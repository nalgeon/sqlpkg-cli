package cmd

import (
	"os"
	"testing"

	"sqlpkg.org/cli/spec"
)

func TestReadLockfile(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		SetupTestRepo(t)
		defer TeardownTestRepo(t)
		CopyTestRepo(t)

		lck, err := ReadLockfile()
		if err != nil {
			t.Fatalf("ReadLockfile: unexpected error %v", err)
		}
		if len(lck.Packages) != 2 {
			t.Fatalf("ReadLockfile: unexpected package count %v", len(lck.Packages))
		}
	})
	t.Run("new", func(t *testing.T) {
		SetupTestRepo(t)
		defer TeardownTestRepo(t)

		lck, err := ReadLockfile()
		if err != nil {
			t.Fatalf("ReadLockfile: unexpected error %v", err)
		}
		if len(lck.Packages) != 0 {
			t.Fatalf("ReadLockfile: unexpected package count %v", len(lck.Packages))
		}
	})
	t.Run("invalid", func(t *testing.T) {
		_, lockPath := SetupTestRepo(t)
		defer TeardownTestRepo(t)
		err := os.WriteFile(lockPath, []byte("invalid"), 0644)
		if err != nil {
			t.Fatalf("os.WriteFile: %v", err)
		}

		_, err = ReadLockfile()
		if err == nil {
			t.Fatal("ReadLockfile: expected error, got nil")
		}
	})
}

func TestAddToLockfile(t *testing.T) {
	SetupTestRepo(t)
	defer TeardownTestRepo(t)
	CopyTestRepo(t)

	lck, err := ReadLockfile()
	if err != nil {
		t.Fatalf("ReadLockfile: %v", err)
	}

	pkg := &spec.Package{Owner: "nalgeon", Name: "text", Version: "0.5.0"}
	err = AddToLockfile(lck, pkg)
	if err != nil {
		t.Fatalf("AddToLockfile: unexpected error %v", err)
	}

	lck, err = ReadLockfile()
	if err != nil {
		t.Fatalf("ReadLockfile: %v", err)
	}

	if len(lck.Packages) != 3 {
		t.Fatalf("AddToLockfile: unexpected package count %v", len(lck.Packages))
	}
	if !lck.Has("nalgeon/text") {
		t.Errorf("AddToLockfile: unexpected packages %v", lck.Packages)
	}
}

func TestRemoveFromLockfile(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		SetupTestRepo(t)
		defer TeardownTestRepo(t)
		CopyTestRepo(t)

		lck, err := ReadLockfile()
		if err != nil {
			t.Fatalf("ReadLockfile: %v", err)
		}

		err = RemoveFromLockfile(lck, "nalgeon/example")
		if err != nil {
			t.Fatalf("RemoveFromLockfile: unexpected error %v", err)
		}

		lck, err = ReadLockfile()
		if err != nil {
			t.Fatalf("ReadLockfile: %v", err)
		}

		if len(lck.Packages) != 1 {
			t.Fatalf("RemoveFromLockfile: unexpected package count %v", len(lck.Packages))
		}
		if lck.Has("nalgeon/example") {
			t.Errorf("RemoveFromLockfile: unexpected packages %v", lck.Packages)
		}
	})
	t.Run("not founr", func(t *testing.T) {
		SetupTestRepo(t)
		defer TeardownTestRepo(t)
		CopyTestRepo(t)

		lck, err := ReadLockfile()
		if err != nil {
			t.Fatalf("ReadLockfile: %v", err)
		}

		err = RemoveFromLockfile(lck, "nalgeon/missing")
		if err != nil {
			t.Fatalf("RemoveFromLockfile: unexpected error %v", err)
		}

		lck, err = ReadLockfile()
		if err != nil {
			t.Fatalf("ReadLockfile: %v", err)
		}

		if len(lck.Packages) != 2 {
			t.Fatalf("RemoveFromLockfile: unexpected package count %v", len(lck.Packages))
		}
	})

}
