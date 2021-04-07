package strategy

import (
	"github.com/Masterminds/semver/v3"
)

type VersionReader interface {
	ReadVersion() (*semver.Version, error)
}

type VersionBumper interface {
	BumpVersion(previous semver.Version) (*semver.Version, error)
}
