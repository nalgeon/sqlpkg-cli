package spec

import (
	"path"
	"path/filepath"
	"testing"
)

func TestAssetPath_Exists(t *testing.T) {
	t.Run("local", func(t *testing.T) {
		p := &AssetPath{Value: filepath.Join("testdata", "sqlpkg.json"), IsRemote: false}
		if !p.Exists() {
			t.Errorf("Exists: expected true for %v", p.Value)
		}
		p = &AssetPath{Value: filepath.Join("testdata", "null.json"), IsRemote: false}
		if p.Exists() {
			t.Errorf("Exists: expected false for %v", p.Value)
		}
	})
}

func TestAssetPath_Join(t *testing.T) {
	t.Run("local", func(t *testing.T) {
		p := &AssetPath{Value: "testdata", IsRemote: false}
		p = p.Join("sqlpkg.json")
		if p.Value != filepath.Join("testdata", "sqlpkg.json") {
			t.Errorf("Join: unexpected Value = %v", p.Value)
		}
		if p.IsRemote {
			t.Errorf("Join: unexpected IsRemote = %v", p.IsRemote)
		}
	})
	t.Run("remote", func(t *testing.T) {
		p := &AssetPath{Value: "testdata", IsRemote: true}
		p = p.Join("sqlpkg.json")
		if p.Value != path.Join("testdata", "sqlpkg.json") {
			t.Errorf("Join: unexpected Value = %v", p.Value)
		}
		if !p.IsRemote {
			t.Errorf("Join: unexpected IsRemote = %v", p.IsRemote)
		}
	})
}

func TestAssetPath_MarshalText(t *testing.T) {
	t.Run("local", func(t *testing.T) {
		src := filepath.Join("testdata", "sqlpkg.json")
		p := &AssetPath{Value: src, IsRemote: false}
		val, err := p.MarshalText()
		if err != nil {
			t.Errorf("MarshalText: unexpected error %v", err)
		}
		if string(val) != src {
			t.Errorf("MarshalText: unexpected value %s", val)
		}
	})
	t.Run("remote", func(t *testing.T) {
		src := path.Join("testdata", "sqlpkg.json")
		p := &AssetPath{Value: src, IsRemote: true}
		val, err := p.MarshalText()
		if err != nil {
			t.Errorf("MarshalText: unexpected error %v", err)
		}
		if string(val) != src {
			t.Errorf("MarshalText: unexpected value %s", val)
		}
	})
}

func TestAssetPath_UnmarshalText(t *testing.T) {
	t.Run("local", func(t *testing.T) {
		src := filepath.Join("testdata", "sqlpkg.json")
		p := new(AssetPath)
		err := p.UnmarshalText([]byte(src))
		if err != nil {
			t.Errorf("UnmarshalText: unexpected error %v", err)
		}
		if p.Value != src {
			t.Errorf("UnmarshalText: unexpected Value = %v", p.Value)
		}
		if p.IsRemote {
			t.Errorf("UnmarshalText: unexpected IsRemote = %v", p.IsRemote)
		}
	})
	t.Run("remote", func(t *testing.T) {
		src := "https://antonz.org/sqlpkg.json"
		p := new(AssetPath)
		err := p.UnmarshalText([]byte(src))
		if err != nil {
			t.Errorf("UnmarshalText: unexpected error %v", err)
		}
		if p.Value != src {
			t.Errorf("UnmarshalText: unexpected Value = %v", p.Value)
		}
		if !p.IsRemote {
			t.Errorf("UnmarshalText: unexpected IsRemote = %v", p.IsRemote)
		}
	})
	t.Run("repository", func(t *testing.T) {
		src := "{repository}/raw/main/sqlpkg.json"
		p := new(AssetPath)
		err := p.UnmarshalText([]byte(src))
		if err != nil {
			t.Errorf("UnmarshalText: unexpected error %v", err)
		}
		if p.Value != src {
			t.Errorf("UnmarshalText: unexpected Value = %v", p.Value)
		}
		if !p.IsRemote {
			t.Errorf("UnmarshalText: unexpected IsRemote = %v", p.IsRemote)
		}
	})
}

func TestAssetPath_String(t *testing.T) {
	t.Run("local", func(t *testing.T) {
		src := filepath.Join("testdata", "sqlpkg.json")
		p := &AssetPath{Value: src, IsRemote: false}
		if p.String() != src {
			t.Errorf("String: unexpected value %v", p.String())
		}
	})
	t.Run("remote", func(t *testing.T) {
		src := filepath.Join("testdata", "sqlpkg.json")
		p := &AssetPath{Value: src, IsRemote: true}
		if p.String() != src {
			t.Errorf("String: unexpected value %v", p.String())
		}
	})
}
