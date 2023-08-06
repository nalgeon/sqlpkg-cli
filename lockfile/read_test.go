package lockfile

import (
	"path/filepath"
	"testing"

	"sqlpkg.org/cli/spec"
)

func TestRead(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		path := filepath.Join("testdata", "sqlpkg.lock")
		got, err := ReadLocal(path)
		if err != nil {
			t.Fatalf("ReadLocal: unexpected error %v", err)
		}

		want := NewLockfile()
		pkg1 := &spec.Package{Owner: "nalgeon", Name: "example", Version: "0.1.0"}
		pkg2 := &spec.Package{Owner: "sqlite", Name: "stmt", Version: "0.21.5"}
		want.Add(pkg1)
		want.Add(pkg2)

		if len(got.Packages) != len(want.Packages) {
			t.Errorf("Read: unexpacted package count %v", len(got.Packages))
		}
		got1 := got.Packages[pkg1.FullName()]
		if got1.FullName() != pkg1.FullName() {
			t.Errorf("Read: unexpected package %v", got1)
		}
		got2 := got.Packages[pkg2.FullName()]
		if got2.FullName() != pkg2.FullName() {
			t.Errorf("Read: unexpected package %v", got2)
		}
	})
	t.Run("failure", func(t *testing.T) {
		path := filepath.Join("testdata", "missing", "sqlpkg.json")
		_, err := ReadLocal(path)
		if err == nil {
			t.Fatal("Read: expected error, got nil")
		}
	})
}
