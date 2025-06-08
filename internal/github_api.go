package githubapi

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

const bsaeURL = "https://api.github.com"

type Release struct {
	Version Semver `json:"tag_name"`
}

func GetLatestRelease(client *http.Client, repo string) (Release, error) {
	body, err := GetLatestReleaseJSON(client, repo)
	if err != nil {
		return Release{}, err
	}

	var rel Release
	if err := json.Unmarshal(body, &rel); err != nil {
		return Release{}, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	return rel, nil
}

// https://docs.github.com/en/rest/releases/releases?apiVersion=2022-11-28#get-the-latest-release
func GetLatestReleaseJSON(client *http.Client, repo string) ([]byte, error) {
	url := fmt.Sprintf("%s/repos/%s/releases/latest", bsaeURL, repo)
	body, err := get(client, url)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest release: %w", err)
	}
	return body, nil
}

func GetNuspecVersion(client *http.Client, repo, path string) (Semver, error) {
	body, err := GetRepositoryContentJSON(client, repo, path)
	if err != nil {
		return Semver{}, fmt.Errorf("failed to get nuspec content: %w", err)
	}

	var response struct {
		Content string `json:"content"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return Semver{}, fmt.Errorf("failed to unmarshal nuspec content: %w", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(response.Content)
	if err != nil {
		return Semver{}, fmt.Errorf("failed to decode base64 content: %w", err)
	}

	var nuspec struct {
		Metadata struct {
			Version Semver `xml:"version"`
		} `xml:"metadata"`
	}
	xml.Unmarshal(decoded, &nuspec)

	return nuspec.Metadata.Version, nil
}

// https://docs.github.com/en/rest/repos/contents?apiVersion=2022-11-28#get-repository-content
func GetRepositoryContentJSON(client *http.Client, repo, path string) ([]byte, error) {
	url := fmt.Sprintf("%s/repos/%s/contents/%s", bsaeURL, repo, path)
	body, err := get(client, url)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository content: %w", err)
	}
	return body, nil
}

func get(client *http.Client, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("x-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return body, nil
}

type WorkflowDispatchInput struct {
	Package string `json:"package"`
}

type WorkflowDispatchPayload struct {
	Ref    string                `json:"ref"`
	Inputs WorkflowDispatchInput `json:"inputs"`
}

// https://docs.github.com/en/rest/actions/workflows?apiVersion=2022-11-28#create-a-workflow-dispatch-event
func WorkflowDispatch(client *http.Client, githubToken, repo, workflowID string) error {
	url := fmt.Sprintf("%s/repos/%s/actions/workflows/%s/dispatches", bsaeURL, repo, workflowID)

	payload := WorkflowDispatchPayload{
		Ref: "main",
		Inputs: WorkflowDispatchInput{
			Package: "azure-kubelogin",
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("x-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Authorization", "Bearer "+githubToken)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d, %s", resp.StatusCode, body)
	}

	return nil
}
