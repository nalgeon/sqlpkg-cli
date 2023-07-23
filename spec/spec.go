// Package spec manages the package spec file (sqlpkg.json).
package spec

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"sqlpkg.org/cli/httpx"
)

// DirName is the name of the folder with packages
const DirName = ".sqlpkg"

// FileName is the package spec filename.
const FileName = "sqlpkg.json"

// downloadBase determines default asset url for known providers.
var downloadBase = map[string]string{
	"github.com": "{repository}/releases/download/{version}",
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
	Assets      Assets   `json:"assets"`
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
	version := p.Version
	if version == "latest" {
		// "latest" is a placeholder, so keep it as a variable
		// to replace later when the package is actually installed
		version = "{latest}"
	}
	p.Assets.Path.Value = stringFormat(p.Assets.Path.Value, map[string]any{
		"repository": p.Repository,
		"owner":      p.Owner,
		"name":       p.Name,
		"version":    version,
	})
	for platform, file := range p.Assets.Files {
		p.Assets.Files[platform] = stringFormat(file, map[string]any{
			"version": version,
		})
	}
}

// ReplaceLatest forces a specific package version instead of the "latest" placeholder.
func (p *Package) ReplaceLatest(version string) {
	if p.Version != "latest" {
		return
	}
	p.Version = version
	p.Assets.Path.Value = strings.Replace(p.Assets.Path.Value, "{latest}", version, 1)
	for platform, file := range p.Assets.Files {
		p.Assets.Files[platform] = strings.Replace(file, "{latest}", version, 1)
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
	hostname := httpx.Hostname(repoUrl)
	return downloadBase[hostname]
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
