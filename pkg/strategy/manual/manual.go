package manual

import (
	"github.com/Masterminds/semver/v3"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
)

type Strategy struct {
	Version string
}

func (s Strategy) ReadVersion() (*semver.Version, error) {
	log.Logger().Debugf("Using manual version %s", s.Version)
	return semver.NewVersion(s.Version)
}

func (s Strategy) BumpVersion(_ semver.Version) (*semver.Version, error) {
	log.Logger().Debugf("Using manual version %s", s.Version)
	return semver.NewVersion(s.Version)
}
