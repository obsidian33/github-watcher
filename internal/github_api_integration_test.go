package githubapi_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	githubapi "github.com/obsidian33/github-watcher/internal"
)

const baseURL = "https://api.github.com"

func TestGitHubAPI(t *testing.T) {

	t.Run("Get latest release", func(t *testing.T) {
		// https://docs.github.com/en/rest/releases/releases?apiVersion=2022-11-28#get-the-latest-release
		url := fmt.Sprintf("%s/repos/Azure/kubelogin/releases/latest", baseURL)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatalf("Failed to make request %v", err)
		}

		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("x-GitHub-Api-Version", "2022-11-28")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200 OK, got %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)

		var g githubapi.Release
		if err := json.Unmarshal(body, &g); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		fmt.Printf("Latest release version: %s\nPublished at: %s\n", g.Version, g.PublishedAt.String())
	})

}
