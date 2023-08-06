package github

import (
	"testing"

	"sqlpkg.org/cli/httpx"
)

func TestGetLatestTag(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		httpx.Mock("valid")
		tag, err := GetLatestTag("nalgeon", "sqlean")
		if err != nil {
			t.Fatalf("GetLatestTag: unexpected error %v", err)
		}
		if tag != "0.21.6" {
			t.Errorf("GetLatestTag: unexpected tag %v", tag)
		}
	})
	t.Run("invalid", func(t *testing.T) {
		httpx.Mock()
		_, err := GetLatestTag("nalgeon", "sqlean")
		if err == nil {
			t.Fatal("GetLatestTag: expected error, got nil")
		}
	})
}

func TestParseRepoUrl(t *testing.T) {
	type test struct {
		url         string
		owner, repo string
	}
	valid := []test{
		{"https://github.com/nalgeon/sqlean", "nalgeon", "sqlean"},
		{"https://github.com/nalgeon/sqlean/", "nalgeon", "sqlean"},
		{"https://github.com/asg017/sqlite-vss", "asg017", "sqlite-vss"},
	}
	for _, test := range valid {
		owner, repo, err := ParseRepoUrl(test.url)
		if err != nil {
			t.Errorf("ParseRepoUrl(%s): unexpected error %v", test.url, err)
			continue
		}
		if owner != test.owner {
			t.Errorf("ParseRepoUrl(%s): unexpected owner %v", test.url, test.owner)
		}
		if repo != test.repo {
			t.Errorf("ParseRepoUrl(%s): unexpected name %v", test.url, test.repo)
		}
	}

	invalid := []string{
		"https://github.com/nalgeon",
		"https://antonz.org",
	}
	for _, url := range invalid {
		_, _, err := ParseRepoUrl(url)
		if err == nil {
			t.Errorf("ParseRepoUrl(%s): expected error, got nil", url)
		}
	}
}
