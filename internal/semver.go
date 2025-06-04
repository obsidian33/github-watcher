package githubapi

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrSyntax = errors.New("syntax error")

type Semver struct{ Major, Minor, Patch int }

func ParseSemver(ver string) (Semver, error) {
	parts := strings.Split(ver, ".")
	if len(parts) != 3 {
		return Semver{}, fmt.Errorf("expected 3 parts, got %d: %w", len(parts), ErrSyntax)
	}

	if strings.HasPrefix(parts[0], "v") || strings.HasPrefix(parts[0], "V") {
		parts[0] = parts[0][1:]
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return Semver{}, fmt.Errorf("invalid major version %q: %w", parts[0], ErrSyntax)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return Semver{}, fmt.Errorf("invalid minor version %q: %w", parts[1], ErrSyntax)
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return Semver{}, fmt.Errorf("invalid patch version %q: %w", parts[2], ErrSyntax)
	}

	return Semver{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

func (s Semver) String() string {
	return fmt.Sprintf("%d.%d.%d", s.Major, s.Minor, s.Patch)
}

func (s Semver) GreaterThan(other Semver) bool {
	if s.Major != other.Major {
		return s.Major > other.Major
	}
	if s.Minor != other.Minor {
		return s.Minor > other.Minor
	}
	return s.Patch > other.Patch
}
