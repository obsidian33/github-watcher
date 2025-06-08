package main

import (
	"log"
	"net/http"
	"os"

	. "github.com/obsidian33/github-watcher/internal"
)

func main() {
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		log.Fatal("GITHUB_TOKEN is not set")
	}

	client := &http.Client{}
	rel, err := GetLatestRelease(client, "Azure/kubelogin")
	if err != nil {
		log.Fatalf("failed to get latest release: %v", err)
	}

	nuspecVer, err := GetNuspecVersion(
		client,
		"obsidian33/chocolatey-packages",
		"azure-kubelogin/azure-kubelogin.nuspec",
	)
	if err != nil {
		log.Fatalf("failed to get nuspec version: %v", err)
	}

	if !rel.Version.GreaterThan(nuspecVer) {
		log.Printf("No new release available. Current version: %s, Nuspec version: %s", rel.Version, nuspecVer)
		return
	}

	log.Printf("New release available: %s (Nuspec version: %s)", rel.Version, nuspecVer)
	err = WorkflowDispatch(client, githubToken, "obsidian22/chocolatey-packages", "chocolatey-pacakge-dispatch.yaml")
	if err != nil {
		log.Fatalf("failed to trigger workflow dispatch: %v", err)
	}
}
