package github

import (
	"fmt"

	"github.com/nalgeon/sqlpkg-cli/internal/httpx"
)

const base_url = "https://api.github.com"

type Release struct {
	TagName string `json:"tag_name"`
}

func GetLatestVersion(owner, repo string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases/latest", base_url, owner, repo)
	var rel Release
	err := httpx.GetJSON(url, &rel)
	if err != nil {
		return "", err
	}
	return rel.TagName, nil
}
