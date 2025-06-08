package githubapi_test

import (
	"net/http"
	"os"
	"testing"

	approvals "github.com/approvals/go-approval-tests"
	. "github.com/obsidian33/github-watcher/internal"
	reporters "github.com/obsidian33/go-approval-reporters"
)

func TestGitHubAPI(t *testing.T) {

	t.Run("Get latest release", func(t *testing.T) {
		got, err := GetLatestReleaseJSON(&http.Client{}, "Azure/kubelogin")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		approvals.VerifyJSONBytes(t, got)
	})

	t.Run("Get repository content", func(t *testing.T) {
		got, err := GetRepositoryContentJSON(
			&http.Client{},
			"obsidian33/chocolatey-packages",
			"azure-kubelogin/azure-kubelogin.nuspec",
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		approvals.VerifyJSONBytes(t, got)
	})

	t.Run("Workflow dispatch", func(t *testing.T) {
		err := WorkflowDispatch(
			&http.Client{},
			os.Getenv("GITHUB_TOKEN"),
			"obsidian33/chocolatey-packages",
			"chocolatey-pacakge-dispatch.yaml",
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

}

func TestMain(m *testing.M) {
	r := approvals.UseReporter(reporters.NewDeltaDiffReporter())
	defer r.Close()

	approvals.UseFolder("test-approvals")
	os.Exit(m.Run())
}
