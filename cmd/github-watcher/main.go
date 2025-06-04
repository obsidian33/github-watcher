package main

import (
	"log"
	"os"
)

func main() {
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		log.Fatal("GITHUB_TOKEN is not set")
	}

	releasePath := os.Getenv("RELEASE_VERSION_PATH")
	if releasePath == "" {
		log.Fatal("RELEASE_VERSION_PATH is not set")
	}

}
