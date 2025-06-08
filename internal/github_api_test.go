package githubapi_test

import (
	"encoding/base64"
	"io"
	"net/http"
	"strings"
	"testing"

	. "github.com/obsidian33/github-watcher/internal"
)

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
	t.Run("happy path", func(t *testing.T) {
		json := `{"tag_name": "v2.1.34"}`
		client := &http.Client{
			Transport: &fakeRoundTripper{respBody: json, status: http.StatusOK},
		}

		got, err := GetLatestRelease(client, "Azure/kubelogin")
		if err != nil {
			t.Fatalf("unecpected error: %v", err)
		}

		want := Release{
			Version: Semver{
				Major: 2,
				Minor: 1,
				Patch: 34,
			},
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestGetNuspecVersion(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		xml := `<?xml version="1.0" encoding="utf-8"?>
		<package xmlns="http://schemas.microsoft.com/packaging/2015/06/nuspec.xsd">
			<metadata>
				<version>0.2.8</version>
			</metadata>
		</package>`
		content := base64.StdEncoding.EncodeToString([]byte(xml))
		json := `{"content": "` + content + `"}`
		client := &http.Client{
			Transport: &fakeRoundTripper{respBody: json, status: http.StatusOK},
		}

		got, err := GetNuspecVersion(client, "owner/repo", "path/to/a.nuspec")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := Semver{Major: 0, Minor: 2, Patch: 8}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestWorkflowDispatch(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		client := &http.Client{
			Transport: &fakeRoundTripper{status: http.StatusNoContent},
		}

		err := WorkflowDispatch(client, "TOKEN", "owner/repo", "workflow_id")

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
