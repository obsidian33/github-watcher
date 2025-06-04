package githubapi_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	. "github.com/obsidian33/github-watcher/internal"
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
		{"2.1.34", "v2.1.0", true},
		{"v2.1.4", "v2.1.5", false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprint(c), func(t *testing.T) {
			got := IsNewRelease(c.latest, c.stored)
			if got != c.want {
				t.Errorf("latest %q, stored %q is new release %t", c.latest, c.stored, got)
			}
		})
	}
}

type fakeRoundTripper struct {
	respBody string
	status   int
}

func (f *fakeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: f.status,
		Body:       io.NopCloser(strings.NewReader(f.respBody)),
		Header:     make(http.Header),
	}, nil
}

func TestGetLatestRelease(t *testing.T) {
	json := `{"tag_name": "v2.1.34", "published_at": "2023-10-01T12:00:00Z"}`
	client := &http.Client{
		Transport: &fakeRoundTripper{respBody: json, status: http.StatusOK},
	}

	got, err := GetLatestRelease(client, "Azure/kubelogin")
	if err != nil {
		t.Fatalf("unecpected error: %v", err)
	}

	want := Release{Version: "v2.1.34", PublishedAt: time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)}
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
