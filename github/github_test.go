package github

import (
	"fmt"
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
		num         int
		owner, repo string
	}
	valid := map[string]test{
		"https://github.com/nalgeon/sqlean":    {0, "nalgeon", "sqlean"},
		"https://github.com/nalgeon/sqlean/":   {1, "nalgeon", "sqlean"},
		"https://github.com/asg017/sqlite-vss": {2, "asg017", "sqlite-vss"},
	}
	for url, test := range valid {
		name := fmt.Sprintf("valid_%d", test.num)
		t.Run(name, func(t *testing.T) {
			owner, repo, err := ParseRepoUrl(url)
			if err != nil {
				t.Fatalf("ParseRepoUrl: unexpected error %v", err)
			}
			if owner != test.owner {
				t.Errorf("ParseRepoUrl: unexpected owner %v", test.owner)
			}
			if repo != test.repo {
				t.Errorf("ParseRepoUrl: unexpected name %v", test.repo)
			}
		})
	}

	invalid := []string{
		"https://github.com/nalgeon",
		"https://antonz.org",
	}
	for idx, url := range invalid {
		name := fmt.Sprintf("valid_%d", idx)
		t.Run(name, func(t *testing.T) {
			_, _, err := ParseRepoUrl(url)
			if err == nil {
				t.Fatal("ParseRepoUrl: expected error, got nil")
			}
		})
	}
}
