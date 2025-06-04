package githubapi_test

import (
	"errors"
	"fmt"
	"testing"

	. "github.com/obsidian33/github-watcher/internal"
)

type parseSemVerTest struct {
	in  string
	out Semver
	err error
}

func TestParseSemver(t *testing.T) {
	cases := []struct {
		version string
		want    Semver
		err     error
	}{
		{"v2.1.34", Semver{Major: 2, Minor: 1, Patch: 34}, nil},
		{"V2.1.34", Semver{Major: 2, Minor: 1, Patch: 34}, nil},
		{"1.10.1", Semver{Major: 1, Minor: 10, Patch: 1}, nil},
		{"2.10.1.0", Semver{}, ErrSyntax},
		{"foo", Semver{}, ErrSyntax},
		{"01.0.0", Semver{Major: 1, Minor: 0, Patch: 0}, nil},
	}

	for _, c := range cases {
		t.Run(c.version, func(t *testing.T) {
			got, err := ParseSemver(c.version)
			if got != c.want || !errors.Is(err, c.err) {
				t.Errorf("ParseSemVer(%q) = %v, %v want %v, %v",
					c.version, got, err, c.want, c.err)
			}
		})
	}

}

func TestGreaterThan(t *testing.T) {
	cases := []struct {
		a, b Semver
		want bool
	}{
		{
			Semver{Major: 1, Minor: 0, Patch: 0},
			Semver{Major: 0, Minor: 1, Patch: 0},
			true,
		},
		{
			Semver{Major: 1, Minor: 0, Patch: 0},
			Semver{Major: 1, Minor: 0, Patch: 0},
			false,
		},
		{
			Semver{Major: 1, Minor: 1, Patch: 0},
			Semver{Major: 1, Minor: 0, Patch: 0},
			true,
		},
		{
			Semver{Major: 1, Minor: 1, Patch: 10},
			Semver{Major: 1, Minor: 0, Patch: 0},
			true,
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%v is greater than %v", c.a, c.b), func(t *testing.T) {
			got := c.a.GreaterThan(c.b)
			if got != c.want {
				t.Errorf("got %v, want %v", got, c.want)
			}
		})
	}
}
