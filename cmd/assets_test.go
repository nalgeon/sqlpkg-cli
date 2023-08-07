package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"sqlpkg.org/cli/assets"
	"sqlpkg.org/cli/fileio"
	"sqlpkg.org/cli/httpx"
	"sqlpkg.org/cli/spec"
)

func TestBuildAssetPath(t *testing.T) {
	httpx.Mock()
	t.Run("exists", func(t *testing.T) {
		pkg := &spec.Package{
			Owner: "nalgeon", Name: "example", Version: "0.1.0",
			Assets: spec.Assets{
				Path: &spec.AssetPath{Value: "https://antonz.org", IsRemote: true},
				Files: map[string]string{
					"darwin-arm64": "example-darwin.zip",
					"linux-amd64":  "example-linux.zip",
				},
			},
		}

		path, err := BuildAssetPath(pkg)
		if err != nil {
			t.Fatalf("BuildAssetPath: unexpected error %v", err)
		}
		want := fmt.Sprintf("https://antonz.org/example-%s.zip", runtime.GOOS)
		if path.Value != want {
			t.Errorf("BuildAssetPath: unexpected Value %q", path.Value)
		}
		if !path.IsRemote {
			t.Errorf("BuildAssetPath: unexpected IsRemote %v", path.IsRemote)
		}
	})
	t.Run("unsupported platform", func(t *testing.T) {
		pkg := &spec.Package{
			Owner: "nalgeon", Name: "example", Version: "0.1.0",
		}

		_, err := BuildAssetPath(pkg)
		if err == nil {
			t.Fatal("BuildAssetPath: expected error, got nil")
		}
	})
}

func TestDownloadAsset(t *testing.T) {
	httpx.Mock()
	t.Run("http", func(t *testing.T) {
		pkg := &spec.Package{
			Owner: "nalgeon", Name: "example", Version: "0.1.0",
			Assets: spec.Assets{
				Path: &spec.AssetPath{Value: "https://antonz.org", IsRemote: true},
				Files: map[string]string{
					"darwin-arm64": "example-darwin.zip",
					"linux-amd64":  "example-linux.zip",
				},
			},
		}
		path := &spec.AssetPath{
			Value:    fmt.Sprintf("%s/example-%s.zip", pkg.Assets.Path, runtime.GOOS),
			IsRemote: true,
		}

		asset, err := DownloadAsset(pkg, path)
		if err != nil {
			t.Fatalf("DownloadAsset: unexpected error %v", err)
		}
		if asset.Name != fmt.Sprintf("example-%s.zip", runtime.GOOS) {
			t.Errorf("DownloadAsset: unexpected Name %v", asset.Name)
		}
		if filepath.Base(asset.Path) != asset.Name {
			t.Errorf("DownloadAsset: unexpected Path %v", asset.Path)
		}
		if !fileio.Exists(asset.Path) {
			t.Error("DownloadAsset: file does not exist")
		}
		if asset.Size != 128 && asset.Size != 137 {
			t.Errorf("DownloadAsset: unexpected asset size %v", asset.Size)
		}
		if !reflect.DeepEqual(asset.Checksum[:6], []byte{0x17, 0xe2, 0xf2, 0xf9, 0x71, 0x93}) && !reflect.DeepEqual(asset.Checksum[:6], []byte{0x58, 0x38, 0x10, 0x0a, 0x12, 0x9a}) {
			t.Errorf("DownloadAsset: unexpected asset checksum %v", asset.Checksum[:6])
		}
	})
	t.Run("file", func(t *testing.T) {
		pkg := &spec.Package{
			Owner: "nalgeon", Name: "example", Version: "0.1.0",
			Assets: spec.Assets{
				Path: &spec.AssetPath{Value: "./testdata", IsRemote: false},
				Files: map[string]string{
					"darwin-arm64": "example-darwin.zip",
					"linux-amd64":  "example-linux.zip",
				},
			},
		}
		path := &spec.AssetPath{
			Value:    fmt.Sprintf("%s/example-%s.zip", pkg.Assets.Path, runtime.GOOS),
			IsRemote: true,
		}

		asset, err := DownloadAsset(pkg, path)
		if err != nil {
			t.Fatalf("DownloadAsset: unexpected error %v", err)
		}
		if asset.Name != fmt.Sprintf("example-%s.zip", runtime.GOOS) {
			t.Errorf("DownloadAsset: unexpected Name %v", asset.Name)
		}
		if filepath.Base(asset.Path) != asset.Name {
			t.Errorf("DownloadAsset: unexpected Path %v", asset.Path)
		}
		if !fileio.Exists(asset.Path) {
			t.Error("DownloadAsset: file does not exist")
		}
		if asset.Size != 128 && asset.Size != 137 {
			t.Errorf("DownloadAsset: unexpected asset size %v", asset.Size)
		}
		if !reflect.DeepEqual(asset.Checksum[:6], []byte{0x17, 0xe2, 0xf2, 0xf9, 0x71, 0x93}) && !reflect.DeepEqual(asset.Checksum[:6], []byte{0x58, 0x38, 0x10, 0x0a, 0x12, 0x9a}) {
			t.Errorf("DownloadAsset: unexpected asset checksum %v", asset.Checksum[:6])
		}
	})
	t.Run("not found", func(t *testing.T) {
		pkg := &spec.Package{
			Owner: "nalgeon", Name: "other", Version: "0.1.0",
			Assets: spec.Assets{
				Path: &spec.AssetPath{Value: "https://antonz.org", IsRemote: true},
				Files: map[string]string{
					"darwin-arm64": "other-darwin.zip",
					"linux-amd64":  "other-linux.zip",
				},
			},
		}
		path := &spec.AssetPath{
			Value:    fmt.Sprintf("%s/other-%s.zip", pkg.Assets.Path, runtime.GOOS),
			IsRemote: true,
		}

		_, err := DownloadAsset(pkg, path)
		if err == nil {
			t.Fatal("DownloadAsset: expected error, got nil")
		}
	})
}

func TestValidateAsset(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		pkg := &spec.Package{
			Owner: "nalgeon", Name: "example", Version: "0.1.0",
			Assets: spec.Assets{
				Path: &spec.AssetPath{Value: "https://antonz.org", IsRemote: true},
				Files: map[string]string{
					"darwin-arm64": "example-darwin.zip",
					"linux-amd64":  "example-linux.zip",
				},
				Checksums: map[string]string{
					"example-darwin.zip": "sha256-17e2f2f97193",
					"example-linux.zip":  "sha256-5838100a129a",
				},
			},
		}
		asset := &assets.Asset{
			Name:     "example-darwin.zip",
			Path:     "./testdata/example-darwin.zip",
			Size:     137,
			Checksum: []byte{0x17, 0xe2, 0xf2, 0xf9, 0x71, 0x93},
		}

		err := ValidateAsset(pkg, asset)
		if err != nil {
			t.Errorf("ValidateAsset: unexpected error %v", err)
		}
	})
	t.Run("missing", func(t *testing.T) {
		pkg := &spec.Package{
			Owner: "nalgeon", Name: "example", Version: "0.1.0",
			Assets: spec.Assets{
				Path: &spec.AssetPath{Value: "https://antonz.org", IsRemote: true},
				Files: map[string]string{
					"darwin-arm64": "example-darwin.zip",
					"linux-amd64":  "example-linux.zip",
				},
			},
		}
		asset := &assets.Asset{
			Name:     "example-darwin.zip",
			Path:     "./testdata/example-darwin.zip",
			Size:     137,
			Checksum: []byte{0x17, 0xe2, 0xf2, 0xf9, 0x71, 0x93},
		}

		err := ValidateAsset(pkg, asset)
		if err != nil {
			t.Errorf("ValidateAsset: unexpected error %v", err)
		}
	})
	t.Run("invalid", func(t *testing.T) {
		pkg := &spec.Package{
			Owner: "nalgeon", Name: "example", Version: "0.1.0",
			Assets: spec.Assets{
				Path: &spec.AssetPath{Value: "https://antonz.org", IsRemote: true},
				Files: map[string]string{
					"darwin-arm64": "example-darwin.zip",
					"linux-amd64":  "example-linux.zip",
				},
				Checksums: map[string]string{
					"example-darwin.zip": "sha256-17e2f2f97193",
					"example-linux.zip":  "sha256-5838100a129a",
				},
			},
		}
		asset := &assets.Asset{
			Name:     "example-darwin.zip",
			Path:     "./testdata/example-darwin.zip",
			Size:     137,
			Checksum: []byte{0x51, 0x52, 0x53, 0x54, 0x55, 0x56},
		}

		err := ValidateAsset(pkg, asset)
		if err == nil {
			t.Fatal("ValidateAsset: expected error, got nil")
		}
	})
}

func TestUnpackAsset(t *testing.T) {
	t.Run("archived", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "example-darwin.zip")
		_, err := fileio.CopyFile(filepath.Join("testdata", "example-darwin.zip"), path)
		if err != nil {
			t.Fatalf("fileio.CopyFile: unexpected error %v", err)
		}

		pkg := &spec.Package{
			Owner: "nalgeon", Name: "example", Version: "0.1.0",
		}
		asset := &assets.Asset{
			Name: "example-darwin.zip", Path: path,
		}

		err = UnpackAsset(pkg, asset)
		if err != nil {
			t.Fatalf("UnpackAsset: unexpected error %v", err)
		}

		if fileio.Exists(filepath.Join(dir, asset.Name)) {
			t.Fatal("UnpackAsset: asset archive is not deleted")
		}

		if !fileio.Exists(filepath.Join(dir, "example.dylib")) {
			t.Fatal("UnpackAsset: unpacked asset does not exit")
		}
		data, err := os.ReadFile(filepath.Join(dir, "example.dylib"))
		if err != nil {
			t.Fatalf("os.ReadFile: unexpected error %v", err)
		}
		if string(data) != "example.dylib" {
			t.Errorf("UnpackAsset: unexpected unpacked contents: %q", string(data))
		}
	})
	t.Run("unpacked", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "example.dylib")
		_, err := fileio.CopyFile(filepath.Join("testdata", "example.dylib"), path)
		if err != nil {
			t.Fatalf("fileio.CopyFile: unexpected error %v", err)
		}

		pkg := &spec.Package{
			Owner: "nalgeon", Name: "example", Version: "0.1.0",
		}
		asset := &assets.Asset{
			Name: "example.dylib", Path: path,
		}

		err = UnpackAsset(pkg, asset)
		if err != nil {
			t.Fatalf("UnpackAsset: unexpected error %v", err)
		}

		if !fileio.Exists(filepath.Join(dir, asset.Name)) {
			t.Fatal("UnpackAsset: asset does not exit")
		}

		data, err := os.ReadFile(filepath.Join(dir, asset.Name))
		if err != nil {
			t.Fatalf("os.ReadFile: unexpected error %v", err)
		}
		if string(data) != "example.dylib" {
			t.Errorf("UnpackAsset: unexpected contents: %q", string(data))
		}
	})
}

func TestInstallFiles(t *testing.T) {
	SetupTestRepo(t)
	defer TeardownTestRepo(t)

	dir := t.TempDir()
	path := filepath.Join(dir, "example.dylib")
	_, err := fileio.CopyFile(filepath.Join("testdata", "example.dylib"), path)
	if err != nil {
		t.Fatalf("fileio.CopyFile: unexpected error %v", err)
	}

	pkg := &spec.Package{
		Owner: "nalgeon", Name: "example", Version: "0.1.0",
	}
	asset := &assets.Asset{
		Name: "example.dylib", Path: path,
	}

	err = InstallFiles(pkg, asset)
	if err != nil {
		t.Fatalf("InstallFiles: unexpected error %v", err)
	}
	if !fileio.Exists(spec.Path(WorkDir, pkg.Owner, pkg.Name)) {
		t.Errorf("InstallFiles: package spec is not installed")
	}
	installedPath := filepath.Join(spec.Dir(WorkDir, pkg.Owner, pkg.Name), asset.Name)
	if !fileio.Exists(installedPath) {
		t.Errorf("InstallFiles: package asset is not installed")
	}
}
