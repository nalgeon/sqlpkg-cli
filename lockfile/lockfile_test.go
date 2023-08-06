package lockfile

import (
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"sqlpkg.org/cli/fileio"
	"sqlpkg.org/cli/spec"
)

func TestPath(t *testing.T) {
	got := Path("testdata")
	if got != filepath.Join("testdata", FileName) {
		t.Errorf("Path: unexpected value %v", got)
	}
}

func TestLockfile_Has(t *testing.T) {
	pkg := &spec.Package{Owner: "nalgeon", Name: "example", Version: "0.1.0"}
	lck := NewLockfile()
	lck.Packages[pkg.FullName()] = pkg
	{
		const name = "nalgeon/example"
		ok := lck.Has(name)
		if !ok {
			t.Errorf("Has(%s) expected true, got false", name)
		}
	}
	{
		const name = "nalgeon/text"
		ok := lck.Has(name)
		if ok {
			t.Errorf("Has(%s) expected false, got true", name)
		}
	}
}

func TestLockfile_Add(t *testing.T) {
	pkg := &spec.Package{Owner: "nalgeon", Name: "example", Version: "0.1.0"}
	lck := NewLockfile()

	t.Run("add", func(t *testing.T) {
		lck.Add(pkg)
		if len(lck.Packages) != 1 {
			t.Errorf("Add: unexpected package count %v", len(lck.Packages))
		}
		got := lck.Packages[pkg.FullName()]
		if got.FullName() != pkg.FullName() {
			t.Errorf("Add: unexpected package %v", got)
		}
	})
	t.Run("update", func(t *testing.T) {
		lck.Add(pkg)
		upd := &spec.Package{Owner: "nalgeon", Name: "example", Version: "0.2.0"}
		lck.Add(upd)
		if len(lck.Packages) != 1 {
			t.Errorf("Add: unexpected package count %v", len(lck.Packages))
		}
		got := lck.Packages[upd.FullName()]
		if got.Version != upd.Version {
			t.Errorf("Add: unexpected version %v", got.Version)
		}
	})
}

func TestLockfile_Remove(t *testing.T) {
	pkg := &spec.Package{Owner: "nalgeon", Name: "example", Version: "0.1.0"}
	lck := NewLockfile()

	t.Run("remove", func(t *testing.T) {
		lck.Add(pkg)
		lck.Remove(pkg)
		if len(lck.Packages) != 0 {
			t.Errorf("Remove: unexpected package count %v", len(lck.Packages))
		}
	})
	t.Run("does not exist", func(t *testing.T) {
		lck.Remove(pkg)
		if len(lck.Packages) != 0 {
			t.Errorf("Remove: unexpected package count %v", len(lck.Packages))
		}
	})
}

func TestLockfile_Range(t *testing.T) {
	pkg1 := &spec.Package{Owner: "nalgeon", Name: "example", Version: "0.1.0"}
	pkg2 := &spec.Package{Owner: "sqlite", Name: "stmt", Version: "0.21.5"}
	lck := NewLockfile()
	lck.Packages[pkg1.FullName()] = pkg1
	lck.Packages[pkg2.FullName()] = pkg2

	t.Run("range", func(t *testing.T) {
		names := []string{}
		lck.Range(func(fullName string, pkg *spec.Package) bool {
			name := fmt.Sprintf("%s/%s@%s", pkg.Owner, pkg.Name, pkg.Version)
			names = append(names, name)
			return true
		})
		sort.Strings(names)
		if !reflect.DeepEqual(names, []string{"nalgeon/example@0.1.0", "sqlite/stmt@0.21.5"}) {
			t.Errorf("Range: unexpected packages %v", names)
		}
	})
	t.Run("break", func(t *testing.T) {
		names := []string{}
		lck.Range(func(fullName string, pkg *spec.Package) bool {
			name := fmt.Sprintf("%s/%s@%s", pkg.Owner, pkg.Name, pkg.Version)
			names = append(names, name)
			return false
		})
		if len(names) != 1 {
			t.Errorf("Range: unexpected package count %v", len(names))
		}
	})
}

func TestSave(t *testing.T) {
	pkg1 := &spec.Package{Owner: "nalgeon", Name: "example", Version: "0.1.0"}
	pkg2 := &spec.Package{Owner: "sqlite", Name: "stmt", Version: "0.21.5"}
	lck := NewLockfile()
	lck.Add(pkg1)
	lck.Add(pkg2)
	dir := t.TempDir()

	err := lck.Save(dir)
	if err != nil {
		t.Fatalf("Save: unexpected error %v", err)
	}
	got, err := fileio.ReadJSON[Lockfile](filepath.Join(dir, FileName))
	if err != nil {
		t.Fatalf("fileio.ReadJSON: unexpected error %v", err)
	}
	if len(got.Packages) != 2 {
		t.Errorf("Save: unexpected package count %v", len(got.Packages))
	}
	got1 := got.Packages[pkg1.FullName()]
	if got1.FullName() != pkg1.FullName() {
		t.Errorf("Save: unexpected package %v", got1)
	}
	got2 := got.Packages[pkg2.FullName()]
	if got2.FullName() != pkg2.FullName() {
		t.Errorf("Save: unexpected package %v", got2)
	}
}
