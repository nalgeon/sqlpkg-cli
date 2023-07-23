package spec

import (
	"path/filepath"
	"strings"

	"sqlpkg.org/cli/fileio"
	"sqlpkg.org/cli/httpx"
)

// An AssetPath describes a local file path or a remote URL.
type AssetPath struct {
	Value    string
	IsRemote bool
}

// Exists checks if the asset actually exists at the said path.
func (p *AssetPath) Exists() bool {
	if p.IsRemote {
		return httpx.Exists(p.Value)
	} else {
		return fileio.Exists(p.Value)
	}
}

// Join appends a filename to the path.
func (p *AssetPath) Join(fileName string) *AssetPath {
	if p.IsRemote {
		return &AssetPath{p.Value + "/" + fileName, true}
	} else {
		return &AssetPath{filepath.Join(p.Value, fileName), false}
	}
}

// MarshalText implements the encoding.TextMarshaler interface.
func (p *AssetPath) MarshalText() ([]byte, error) {
	return []byte(p.Value), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (p *AssetPath) UnmarshalText(text []byte) error {
	p.Value = string(text)
	p.IsRemote = httpx.IsURL(p.Value) || strings.HasPrefix(p.Value, "{repository}")
	return nil
}

// String implements the fmt.Stringer interface.
func (p *AssetPath) String() string {
	return p.Value
}
