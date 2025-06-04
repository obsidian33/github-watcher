package githubapi

import (
	"log"
	"net/http"
	"time"
)

type Release struct {
	Version     string    `json:"tag_name"`
	PublishedAt time.Time `json:"published_at"`
}

func IsNewRelease(latest, stored string) bool {
	latestSemver, err := ParseSemver(latest)
	if err != nil {
		log.Print(err)
		return false
	}

	storedSemver, err := ParseSemver(stored)
	if err != nil {
		log.Print(err)
		return false
	}

	return latestSemver.GreaterThan(storedSemver)
}

func GetLatestRelease(client *http.Client, repo string) (Release, error) {
	return Release{}, nil
}
