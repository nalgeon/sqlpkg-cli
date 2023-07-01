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
	Homepage    string   `json:"homepage"`
	Repository  string   `json:"repository"`
	Authors     []string `json:"authors"`
	License     string   `json:"license"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	Symbols     []string `json:"symbols"`
	Assets      `json:"assets"`

	Path string `json:"-"`
}

// Assets are archives of package files, each for a specific platform.
type Assets struct {
	Path    *AssetPath        `json:"path"`
	Pattern string            `json:"pattern"`
	Files   map[string]string `json:"files"`
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

// A ReadFunc if a function that reads package spec from a given path.
type ReadFunc func(path string) (*Package, error)

// ReadLocal reads package spec from a local file.
func ReadLocal(path string) (*Package, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var pkg Package
	err = json.Unmarshal(data, &pkg)
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

// ReadRemote reads package spec from a remote url.
func ReadRemote(url string) (*Package, error) {
	var pkg Package
	err := httpx.GetJSON(url, &pkg)
	if err != nil {
		return nil, err
	}
	return &pkg, nil
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
