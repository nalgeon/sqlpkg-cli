package spec

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/nalgeon/sqlpkg-cli/internal/fileio"
	"github.com/nalgeon/sqlpkg-cli/internal/httpx"
)

// DirName is the name of the folder with packages
const DirName = ".sqlpkg"

// FileName is the package spec filename.
const FileName = "sqlpkg.json"

// downloadBase determines default asset url for known providers.
var downloadBase = map[string]string{
	"github.com": "{repository}/releases/download/{version}",
}

// An Asset describes a local file path or a remote URL.
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

func (p *AssetPath) MarshalText() ([]byte, error) {
	return []byte(p.Value), nil
}

func (p *AssetPath) UnmarshalText(text []byte) error {
	p.Value = string(text)
	p.IsRemote = httpx.IsURL(p.Value) || strings.HasPrefix(p.Value, "{repository}")
	return nil
}

func (p *AssetPath) String() string {
	return p.Value
}

// A Package describes the package spec.
type Package struct {
	Owner       string   `json:"owner"`
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Homepage    string   `json:"homepage,omitempty"`
	Repository  string   `json:"repository,omitempty"`
	Specfile    string   `json:"specfile,omitempty"`
	Authors     []string `json:"authors,omitempty"`
	License     string   `json:"license,omitempty"`
	Description string   `json:"description,omitempty"`
	Keywords    []string `json:"keywords,omitempty"`
	Symbols     []string `json:"symbols,omitempty"`
	Assets      `json:"assets"`
}

// Assets are archives of package files, each for a specific platform.
type Assets struct {
	Path      *AssetPath        `json:"path"`
	Pattern   string            `json:"pattern,omitempty"`
	Files     map[string]string `json:"files"`
	Checksums map[string]string `json:"checksums,omitempty"`
}

// FullName is an owner-name pair that uniquely identifies the package.
func (p *Package) FullName() string {
	return p.Owner + "/" + p.Name
}

// ExpandVars substitutes variables in Assets with real values.
func (p *Package) ExpandVars() {
	if p.Assets.Path == nil || p.Assets.Path.Value == "" {
		p.Assets.Path = &AssetPath{inferAssetUrl(p.Repository), true}
	}
	p.Assets.Path.Value = stringFormat(p.Assets.Path.Value, map[string]any{
		"repository": p.Repository,
		"owner":      p.Owner,
		"name":       p.Name,
		"version":    p.Version,
	})
	for platform, file := range p.Assets.Files {
		p.Assets.Files[platform] = stringFormat(file, map[string]any{
			"version": p.Version,
		})
	}
}

// AssetPath determines the package url for a specific platform (OS + architecture).
func (p *Package) AssetPath(os, arch string) (*AssetPath, error) {
	platform := os + "-" + arch
	asset, ok := p.Assets.Files[platform]
	if !ok {
		return nil, errors.New("platform is not supported")
	}
	if p.Assets.Path == nil || p.Assets.Path.Value == "" {
		return nil, errors.New("asset path is not set")
	}
	path := p.Assets.Path.Join(asset)
	return path, nil
}

// Save writes the package spec file to the specified directory.
func (p *Package) Save(dir string) error {
	data, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		return err
	}
	path := filepath.Join(dir, FileName)
	return os.WriteFile(path, data, 0644)
}

// Dir returns the package directory.
func Dir(basePath, owner, name string) string {
	return filepath.Join(basePath, DirName, owner, name)
}

// Path returns the path to the package spec file.
func Path(basePath, owner, name string) string {
	return filepath.Join(basePath, DirName, owner, name, FileName)
}

// inferAssetUrl determines an asset url given the package repository url.
func inferAssetUrl(repoUrl string) string {
	url, err := url.Parse(repoUrl)
	if err != nil {
		return ""
	}
	return downloadBase[url.Hostname()]
}

// stringFormat formats a string according to the map of values.
// E.g. stringFormat("hello, {name}", map[string]string{"name": "world"})
// -> "hello, world"
func stringFormat(s string, mapping map[string]any) string {
	args := make([]string, 0, len(mapping)*2)
	for key, val := range mapping {
		args = append(args, "{"+key+"}")
		args = append(args, fmt.Sprint(val))
	}
	return strings.NewReplacer(args...).Replace(s)
}
