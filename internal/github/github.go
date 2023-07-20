package github

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/nalgeon/sqlpkg-cli/internal/httpx"
)

const Hostname = "github.com"
const apiUrl = "https://api.github.com"

type release struct {
	TagName string `json:"tag_name"`
}

// ParseRepoUrl extracts owner and repo names from the repo url.
func ParseRepoUrl(repoUrl string) (owner string, repo string, err error) {
	u, err := url.Parse(repoUrl)
	if err != nil {
		err = errors.New(repoUrl)
		return
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) != 2 {
		err = errors.New(repoUrl)
		return
	}
	owner = parts[0]
	repo = parts[1]
	return
}

// GetLatestTag fetches the latest release tag number for the repository.
func GetLatestTag(owner, repo string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases/latest", apiUrl, owner, repo)
	rel, err := httpx.GetJSON[release](url)
	if err != nil {
		return "", err
	}
	return rel.TagName, nil
}
