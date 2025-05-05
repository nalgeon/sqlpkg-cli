package assets

import (
	"path/filepath"
	"reflect"
	"testing"

	"sqlpkg.org/cli/fileio"
	"sqlpkg.org/cli/httpx"
)

func TestAsset_Dir(t *testing.T) {
	asset := Asset{
		Name:     "example.zip",
		Path:     "/opt/assets/example.zip",
		Size:     42,
		Checksum: []byte{51, 52, 53},
	}
	dir := asset.Dir()
	if dir != "/opt/assets" {
		t.Errorf("Asset.Dir: unexpected value %v", dir)
	}
}

func TestAsset_Validate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		checkStr := "sha256-6b731afc7621"
		asset := Asset{
			Name:     "example.zip",
			Path:     "/opt/assets/example.zip",
			Size:     246,
			Checksum: []byte{0x6b, 0x73, 0x1a, 0xfc, 0x76, 0x21},
		}
		ok, err := asset.Validate(checkStr)
		if err != nil {
			t.Fatalf("Validate: unexpected error %v", err)
		}
		if !ok {
			t.Errorf("Validate: unexpected value %v", ok)
		}
	})
	t.Run("invalid", func(t *testing.T) {
		checkStr := "sha256-6b731afc7621"
		asset := Asset{
			Name:     "example.zip",
			Path:     "/opt/assets/example.zip",
			Size:     246,
			Checksum: []byte{0x51, 0x52, 0x53, 0x54, 0x55, 0x56},
		}
		ok, err := asset.Validate(checkStr)
		if err != nil {
			t.Fatalf("Validate: unexpected error %v", err)
		}
		if ok {
			t.Errorf("Validate: unexpected value %v", ok)
		}
	})
	t.Run("unsupported algo", func(t *testing.T) {
		checkStr := "sha512-6b731afc7621"
		asset := Asset{
			Name:     "example.zip",
			Path:     "/opt/assets/example.zip",
			Size:     246,
			Checksum: []byte{0x6b, 0x73, 0x1a, 0xfc, 0x76, 0x21},
		}
		ok, err := asset.Validate(checkStr)
		if err == nil {
			t.Errorf("Validate: expected error, got nil")
		}
		if ok {
			t.Errorf("Validate: unexpected value %v", ok)
		}
	})
}

func TestDownload(t *testing.T) {
	httpx.Mock()
	dir := t.TempDir()
	t.Run("valid", func(t *testing.T) {
		asset, err := Download(dir, "https://antonz.org/example.zip")
		if err != nil {
			t.Fatalf("Download: unexpected error %v", err)
		}
		if asset.Name != "example.zip" {
			t.Errorf("Download: unexpected Name %v", asset.Name)
		}
		if asset.Path != filepath.Join(dir, "example.zip") {
			t.Errorf("Download: unexpected Path %v", asset.Path)
		}
		if asset.Size != 246 {
			t.Errorf("Download: unexpected Size %v", asset.Size)
		}
		checksum := []byte{0x6b, 0x73, 0x1a, 0xfc, 0x76, 0x21}
		if !reflect.DeepEqual(asset.Checksum[:6], checksum) {
			t.Errorf("Download: unexpected Checksum %v", asset.Checksum[:6])
		}
	})
	t.Run("missing", func(t *testing.T) {
		_, err := Download(dir, "https://antonz.org/missing.zip")
		if err == nil {
			t.Fatal("Download: expected error, got nil")
		}
	})
}

func TestCopy(t *testing.T) {
	path := filepath.Join("testdata", "example.zip")
	dir := t.TempDir()
	asset, err := Copy(dir, path)
	if err != nil {
		t.Fatalf("Copy: unexpected error %v", err)
	}
	if asset.Name != "example.zip" {
		t.Errorf("Download: unexpected Name %v", asset.Name)
	}
	if asset.Path != filepath.Join(dir, "example.zip") {
		t.Errorf("Download: unexpected Path %v", asset.Path)
	}
	if asset.Size != 246 {
		t.Errorf("Download: unexpected Size %v", asset.Size)
	}
	checksum := []byte{0x6b, 0x73, 0x1a, 0xfc, 0x76, 0x21}
	if !reflect.DeepEqual(asset.Checksum[:6], checksum) {
		t.Errorf("Download: unexpected Checksum %v", asset.Checksum[:6])
	}
}

func TestUnpack(t *testing.T) {
	t.Run("unzip", func(t *testing.T) {
		path := filepath.Join("testdata", "example.zip")
		dir := t.TempDir()
		asset, err := Copy(dir, path)
		if err != nil {
			t.Fatalf("Copy: unexpected error %v", err)
		}

		count, err := Unpack(asset.Path, "")
		if err != nil {
			t.Fatalf("Unpack: unexpected error %v", err)
		}
		if count != 2 {
			t.Errorf("Unpack: unexpected count %v", count)
		}
		if !fileio.Exists(filepath.Join(dir, "example.dylib")) {
			t.Error("Unpack: missing example.dylib")
		}
		if !fileio.Exists(filepath.Join(dir, "example.txt")) {
			t.Error("Unpack: missing example.txt")
		}
	})
	t.Run("unzip pattern", func(t *testing.T) {
		path := filepath.Join("testdata", "example.zip")
		dir := t.TempDir()
		asset, err := Copy(dir, path)
		if err != nil {
			t.Fatalf("Copy: unexpected error %v", err)
		}

		count, err := Unpack(asset.Path, "*.dylib")
		if err != nil {
			t.Fatalf("Unpack: unexpected error %v", err)
		}
		if count != 1 {
			t.Errorf("Unpack: unexpected count %v", count)
		}
		if !fileio.Exists(filepath.Join(dir, "example.dylib")) {
			t.Error("Unpack: missing example.dylib")
		}
		if fileio.Exists(filepath.Join(dir, "example.txt")) {
			t.Error("Unpack: unexpected example.txt")
		}
	})
	t.Run("untar", func(t *testing.T) {
		path := filepath.Join("testdata", "example.tar.gz")
		dir := t.TempDir()
		asset, err := Copy(dir, path)
		if err != nil {
			t.Fatalf("Copy: unexpected error %v", err)
		}

		count, err := Unpack(asset.Path, "")
		if err != nil {
			t.Fatalf("Unpack: unexpected error %v", err)
		}
		if count != 2 {
			t.Errorf("Unpack: unexpected count %v", count)
		}
		if !fileio.Exists(filepath.Join(dir, "example.dylib")) {
			t.Error("Unpack: missing example.dylib")
		}
		if !fileio.Exists(filepath.Join(dir, "example.txt")) {
			t.Error("Unpack: missing example.txt")
		}
	})
	t.Run("untar pattern", func(t *testing.T) {
		path := filepath.Join("testdata", "example.tar.gz")
		dir := t.TempDir()
		asset, err := Copy(dir, path)
		if err != nil {
			t.Fatalf("Copy: unexpected error %v", err)
		}

		count, err := Unpack(asset.Path, "*.dylib")
		if err != nil {
			t.Fatalf("Unpack: unexpected error %v", err)
		}
		if count != 1 {
			t.Errorf("Unpack: unexpected count %v", count)
		}
		if !fileio.Exists(filepath.Join(dir, "example.dylib")) {
			t.Error("Unpack: missing example.dylib")
		}
		if fileio.Exists(filepath.Join(dir, "example.txt")) {
			t.Error("Unpack: unexpected example.txt")
		}
	})
	t.Run("gunzip", func(t *testing.T) {
		path := filepath.Join("testdata", "example.so.gz")
		dir := t.TempDir()
		asset, err := Copy(dir, path)
		if err != nil {
			t.Fatalf("Copy: unexpected error %v", err)
		}

		count, err := Unpack(asset.Path, "")
		if err != nil {
			t.Fatalf("Unpack: unexpected error %v", err)
		}
		if count != 1 {
			t.Errorf("Unpack: unexpected count %v", count)
		}
		if !fileio.Exists(filepath.Join(dir, "example.so")) {
			t.Error("Unpack: missing example.so")
		}
	})
}
