package githubapi

import (
	"strings"
	"time"
)

type Release struct {
	Version     string    `json:"tag_name"`
	PublishedAt time.Time `json:"published_at"`
}

type SemVer struct {
	Major int
	Minor int
	Patch int
}

func IsNewRelease(latest, stored string) bool {
	latestVer := strings.Split(strings.TrimPrefix(latest, "v"), ".")
	storedVer := strings.Split(strings.TrimPrefix(stored, "v"), ".")

	if latestVer[0] > storedVer[0] {
		return true
	}

	if latestVer[1] > storedVer[1] {
		return true
	}

	if latestVer[2] > storedVer[2] {
		return true
	}

	return false
}

func ParseSemVer(ver string) SemVer {
	return SemVer{}
}
