// Package semver implements comparison of semantic version strings.
// The general form of a semantic version string is
//
//	MAJOR[.MINOR[.PATCH[-PRERELEASE][+BUILD]]]
//
// Follows SemVer 2.0.0 (see semver.org) with some exceptions:
//   - allows the leading 'v' (e.g. v1.2.3)
//   - treats MAJOR as MAJOR.0.0
//   - treats MAJOR.MINOR as MAJOR.MINOR.0
package semver

import (
	"regexp"
)

// https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string
var re = regexp.MustCompile(`^(0|[1-9]\d*)(?:\.(0|[1-9]\d*))?(?:\.(0|[1-9]\d*))?(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

type version struct {
	maj, min, pat string
	pre, build    string
}

func makeVersion() version {
	return version{maj: "0", min: "0", pat: "0"}
}

func (v version) compare(other version) int {
	if c := compareInt(v.maj, other.maj); c != 0 {
		return c
	}
	if c := compareInt(v.min, other.min); c != 0 {
		return c
	}
	if c := compareInt(v.pat, other.pat); c != 0 {
		return c
	}
	return comparePre(v.pre, other.pre)
}

// Compare returns an integer comparing two versions according to
// semantic version precedence:
//   - 0 if v1 == v2
//   - -1 if v1 < v2
//   - +1 if v1 > v2
//
// An invalid version string is considered less than a valid one.
// Two invalid versions strings are considered equal.
func Compare(v1, v2 string) int {
	ver1, ok1 := parse(v1)
	ver2, ok2 := parse(v2)
	if !ok1 && !ok2 {
		return 0
	}
	if !ok1 {
		return -1
	}
	if !ok2 {
		return +1
	}
	return ver1.compare(ver2)
}

func parse(v string) (version, bool) {
	ver := makeVersion()
	if v == "" {
		return ver, false
	}

	// leading 'v' is optional
	if v[0] == 'v' {
		v = v[1:]
	}

	parts := re.FindStringSubmatch(v)
	if parts == nil {
		// does not match the pattern at all
		return ver, false
	}
	parts = parts[1:]

	// maj, min, pat, pre, build = 0, 1, 2, 3, 4
	if (parts[3] != "" || parts[4] != "") && (parts[1] == "" || parts[2] == "") {
		// prerelease and build are only allowed
		// for full versions (maj.min.pat)
		return ver, false
	}

	// fill version fields
	ver.maj = parts[0]
	if parts[1] != "" {
		ver.min = parts[1]
	}
	if parts[2] != "" {
		ver.pat = parts[2]
	}
	ver.pre = parts[3]
	ver.build = parts[4]
	return ver, true
}

// From here on uses the Go implementation:
// https://cs.opensource.google/go/x/mod/+/master:semver/semver.go
func compareInt(x, y string) int {
	if x == y {
		return 0
	}
	if len(x) < len(y) {
		return -1
	}
	if len(x) > len(y) {
		return +1
	}
	if x < y {
		return -1
	} else {
		return +1
	}
}

func comparePre(x, y string) int {
	// "When major, minor, and patch are equal, a pre-release version has
	// lower precedence than a normal version.
	// Example: 1.0.0-alpha < 1.0.0.
	// Precedence for two pre-release versions with the same major, minor,
	// and patch version MUST be determined by comparing each dot separated
	// identifier from left to right until a difference is found as follows:
	// identifiers consisting of only digits are compared numerically and
	// identifiers with letters or hyphens are compared lexically in ASCII
	// sort order. Numeric identifiers always have lower precedence than
	// non-numeric identifiers. A larger set of pre-release fields has a
	// higher precedence than a smaller set, if all of the preceding
	// identifiers are equal.
	// Example: 1.0.0-alpha < 1.0.0-alpha.1 < 1.0.0-alpha.beta <
	// 1.0.0-beta < 1.0.0-beta.2 < 1.0.0-beta.11 < 1.0.0-rc.1 < 1.0.0."
	if x == y {
		return 0
	}
	if x == "" {
		return +1
	}
	if y == "" {
		return -1
	}
	x, y = "-"+x, "-"+y
	for x != "" && y != "" {
		x = x[1:] // skip - or .
		y = y[1:] // skip - or .
		var dx, dy string
		dx, x = nextIdent(x)
		dy, y = nextIdent(y)
		if dx != dy {
			ix := isNum(dx)
			iy := isNum(dy)
			if ix != iy {
				if ix {
					return -1
				} else {
					return +1
				}
			}
			if ix {
				if len(dx) < len(dy) {
					return -1
				}
				if len(dx) > len(dy) {
					return +1
				}
			}
			if dx < dy {
				return -1
			} else {
				return +1
			}
		}
	}
	if x == "" {
		return -1
	} else {
		return +1
	}
}

func nextIdent(x string) (dx, rest string) {
	i := 0
	for i < len(x) && x[i] != '.' {
		i++
	}
	return x[:i], x[i:]
}

func isNum(v string) bool {
	i := 0
	for i < len(v) && '0' <= v[i] && v[i] <= '9' {
		i++
	}
	return i == len(v)
}
