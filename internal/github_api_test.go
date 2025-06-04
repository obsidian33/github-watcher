package githubapi_test

import (
	"fmt"
	"testing"

	githubapi "github.com/obsidian33/github-watcher/internal"
)

func TestIsNewRelease(t *testing.T) {
	cases := []struct {
		latest, stored string
		want           bool
	}{
		{"v2.0.0", "v1.0.0", true},
		{"v2.0.0", "v2.0.0", false},
		{"v1.0.0", "v2.0.0", false},
		{"v2.1.0", "v2.0.0", true},
		{"v2.1.34", "v2.1.0", true},
		{"v2.1.4", "v2.1.5", false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprint(c), func(t *testing.T) {
			got := githubapi.IsNewRelease(c.latest, c.stored)
			if got != c.want {
				t.Errorf("latest %q, stored %q is new release %t", c.latest, c.stored, got)
			}
		})
	}
}

func TestParseSemVer(t *testing.T) {
	t.Run("parse valid version", func(t *testing.T) {
		ver := "v2.1.34"

		got := githubapi.ParseSemVer(ver)
		want := githubapi.SemVer{
			Major: 2,
			Minor: 1,
			Patch: 34,
		}

		if got != want {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})
}
