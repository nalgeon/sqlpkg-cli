package semver

import "testing"

// Uses test cases from the Go implementation:
// https://cs.opensource.google/go/x/mod/+/master:semver/semver_test.go
var tests = []struct {
	in  string
	out string
}{
	{"", ""},
	{"bad", ""},
	{"v1-alpha.beta.gamma", ""},
	{"v1-pre", ""},
	{"v1+meta", ""},
	{"v1-pre+meta", ""},
	{"v1.2-pre", ""},
	{"v1.2+meta", ""},
	{"v1.2-pre+meta", ""},
	{"v1.0.0-alpha", "v1.0.0-alpha"},
	{"v1.0.0-alpha.1", "v1.0.0-alpha.1"},
	{"v1.0.0-alpha.beta", "v1.0.0-alpha.beta"},
	{"v1.0.0-beta", "v1.0.0-beta"},
	{"v1.0.0-beta.2", "v1.0.0-beta.2"},
	{"v1.0.0-beta.11", "v1.0.0-beta.11"},
	{"v1.0.0-rc.1", "v1.0.0-rc.1"},
	{"v1", "v1.0.0"},
	{"v1.0", "v1.0.0"},
	{"v1.0.0", "v1.0.0"},
	{"v1.2", "v1.2.0"},
	{"v1.2.0", "v1.2.0"},
	{"v1.2.3-456", "v1.2.3-456"},
	{"v1.2.3-456.789", "v1.2.3-456.789"},
	{"v1.2.3-456-789", "v1.2.3-456-789"},
	{"v1.2.3-456a", "v1.2.3-456a"},
	{"v1.2.3-pre", "v1.2.3-pre"},
	{"v1.2.3-pre+meta", "v1.2.3-pre"},
	{"v1.2.3-pre.1", "v1.2.3-pre.1"},
	{"v1.2.3-zzz", "v1.2.3-zzz"},
	{"v1.2.3", "v1.2.3"},
	{"v1.2.3+meta", "v1.2.3"},
	{"v1.2.3+meta-pre", "v1.2.3"},
	{"v1.2.3+meta-pre.sha.256a", "v1.2.3"},
	{"v1.25.0", "v1.25.0"},
}

func TestCompare(t *testing.T) {
	for i, ti := range tests {
		for j, tj := range tests {
			cmp := Compare(ti.in, tj.in)
			var want int
			if ti.out == tj.out {
				want = 0
			} else if i < j {
				want = -1
			} else {
				want = +1
			}
			if cmp != want {
				t.Errorf("Compare(%q, %q) = %d, want %d", ti.in, tj.in, cmp, want)
			}
		}
	}
}
