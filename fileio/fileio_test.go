package fileio

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestExists(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		path := createFile(t, t.TempDir(), "file.txt")
		if !Exists(path) {
			t.Fatalf("Exists: %s does not exist", filepath.Base(path))
		}
	})
	t.Run("does not exist", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "file.txt")
		if Exists(path) {
			t.Fatalf("Exists: %s should not exist", filepath.Base(path))
		}
	})
}

func TestCreateDir(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "new")
		err := CreateDir(dir)
		if err != nil {
			t.Fatalf("CreateDir: unexpected error %v", err)
		}
		if !Exists(dir) {
			t.Fatal("CreateDir: dir does not exist")
		}
	})
	t.Run("existing", func(t *testing.T) {
		dir := createDir(t, "existing")
		createFile(t, dir, "file.txt")

		err := CreateDir(dir)
		if err != nil {
			t.Fatalf("CreateDir: unexpected error %v", err)
		}
		if !Exists(dir) {
			t.Fatal("CreateDir: dir does not exist")
		}
		files, err := filepath.Glob(filepath.Join(dir, "*"))
		if err != nil {
			t.Fatalf("filepath.Glob: unexpected error %v", err)
		}
		if len(files) != 0 {
			t.Fatalf("CreateDir: expected empty dir, got %d files", len(files))
		}
	})
}

func TestMoveDir(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		src := createDir(t, "src")
		createFile(t, src, "file.txt")

		dst := filepath.Join(t.TempDir(), "dst")
		err := MoveDir(src, dst)
		if err != nil {
			t.Fatalf("MoveDir: unexpected error %v", err)
		}

		if Exists(src) {
			t.Fatal("MoveDir: src dir exists")
		}
		if !Exists(dst) {
			t.Fatal("MoveDir: dst dir does not exist")
		}

		files, err := filepath.Glob(filepath.Join(dst, "*"))
		if err != nil {
			t.Fatalf("filepath.Glob: unexpected error %v", err)
		}
		if len(files) != 1 {
			t.Fatalf("MoveDir: expected 1 file in dst dir, got %d", len(files))
		}
	})
	t.Run("existing", func(t *testing.T) {
		src := createDir(t, "src")
		createFile(t, src, "src.txt")
		dst := createDir(t, "dst")
		createFile(t, dst, "dst.txt")

		err := MoveDir(src, dst)
		if err != nil {
			t.Fatalf("MoveDir: unexpected error %v", err)
		}

		if Exists(src) {
			t.Fatal("MoveDir: src dir exists")
		}
		if !Exists(dst) {
			t.Fatal("MoveDir: dst dir does not exist")
		}

		files, err := filepath.Glob(filepath.Join(dst, "*"))
		if err != nil {
			t.Fatalf("filepath.Glob: unexpected error %v", err)
		}
		if len(files) != 1 {
			t.Fatalf("MoveDir: expected 1 file in dst dir, got %d", len(files))
		}
		if filepath.Base(files[0]) != "src.txt" {
			t.Fatal("MoveDir: expected src.txt file in dst dir")
		}
	})
	t.Run("does not exist", func(t *testing.T) {
		src := filepath.Join(t.TempDir(), "src")
		dst := filepath.Join(t.TempDir(), "dst")
		err := MoveDir(src, dst)
		if err == nil {
			t.Fatal("MoveDir: expected error, got nil")
		}
	})
}

func TestCopyFile(t *testing.T) {
	t.Run("copy", func(t *testing.T) {
		dir := t.TempDir()
		src := createFile(t, dir, "src.txt")
		dst := filepath.Join(dir, "dst.txt")

		size, err := CopyFile(src, dst)
		if err != nil {
			t.Fatalf("CopyFile: unexpected error %v", err)
		}
		if !Exists(dst) {
			t.Fatal("CopyFile: dst file does not exist")
		}
		if size != 3 {
			t.Fatalf("CopyFile: expected dst file size = 3, got %d", size)
		}
	})
	t.Run("overwrite", func(t *testing.T) {
		dir := t.TempDir()
		src := createFileWithContents(t, dir, "src.txt", "123456")
		dst := createFile(t, dir, "dst.txt")

		size, err := CopyFile(src, dst)
		if err != nil {
			t.Fatalf("CopyFile: unexpected error %v", err)
		}
		if !Exists(dst) {
			t.Fatal("CopyFile: dst file does not exist")
		}
		if size != 6 {
			t.Fatalf("CopyFile: expected dst file size = 6, got %d", size)
		}
	})
	t.Run("does not exist", func(t *testing.T) {
		dir := t.TempDir()
		src := filepath.Join(dir, "src.txt")
		dst := filepath.Join(dir, "dst.txt")

		_, err := CopyFile(src, dst)
		if err == nil {
			t.Fatal("CopyFile: expected error, got nil")
		}
	})
}

func TestReadJSON(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		type Bag struct{ Size int }
		dir := t.TempDir()
		path := createFileWithContents(t, dir, "bag.json", `{"size":3}`)

		got, err := ReadJSON[Bag](path)
		if err != nil {
			t.Fatalf("ReadJSON: unexpected error %v", err)
		}
		want := Bag{Size: 3}
		if *got != want {
			t.Fatalf("ReadJSON: expected %+v, got %+v", want, got)
		}
	})
	t.Run("invalid", func(t *testing.T) {
		type Bag struct{ Size int }
		dir := t.TempDir()
		path := createFileWithContents(t, dir, "bag.json", `hello`)

		_, err := ReadJSON[Bag](path)
		if err == nil {
			t.Fatal("ReadJSON: expected error, got nil")
		}
	})
	t.Run("does not exist", func(t *testing.T) {
		type Bag struct{ Size int }
		dir := t.TempDir()
		path := filepath.Join(dir, "bag.json")

		_, err := ReadJSON[Bag](path)
		if err == nil {
			t.Fatal("ReadJSON: expected error, got nil")
		}
	})
}

func TestCalcChecksum(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		dir := t.TempDir()
		path := createFileWithContents(t, dir, "example.txt", "example.txt")

		sum, err := CalcChecksum(path)
		if err != nil {
			t.Fatalf("CalcChecksum: unexpected error %v", err)
		}
		if len(sum) != 256/8 {
			t.Fatalf("CalcChecksum: unxpected length %d", len(sum))
		}
		if !reflect.DeepEqual(sum[:6], []byte{0xe7, 0xcb, 0x63, 0x23, 0x59, 0xa2}) {
			t.Fatalf("CalcChecksum: unexpected value %v", sum[:6])
		}
	})
	t.Run("does not exist", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "missing.txt")

		_, err := CalcChecksum(path)
		if err == nil {
			t.Fatal("CalcChecksum: expected error, got nil")
		}
	})
}

func createDir(t *testing.T, name string) string {
	parent := t.TempDir()
	dir := filepath.Join(parent, name)
	err := CreateDir(dir)
	if err != nil {
		t.Fatalf("CreateDir: unexpected error %v", err)
	}
	return dir
}

func createFile(t *testing.T, dir string, name string) string {
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte("123"), 0644)
	if err != nil {
		t.Fatalf("os.WriteFile: unexpected error %v", err)
	}
	return path
}

func createFileWithContents(t *testing.T, dir string, name string, text string) string {
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(text), 0644)
	if err != nil {
		t.Fatalf("os.WriteFile: unexpected error %v", err)
	}
	return path
}
