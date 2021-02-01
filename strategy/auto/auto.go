package auto

import (
	"github.com/Masterminds/semver/v3"
	"github.com/jenkins-x-plugins/jx-release-version/v2/strategy/fromtag"
	"github.com/jenkins-x-plugins/jx-release-version/v2/strategy/semantic"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
)

type Strategy struct {
	FromTagStrategy  fromtag.Strategy
	SemanticStrategy semantic.Strategy
}

func (s Strategy) ReadVersion() (*semver.Version, error) {
	log.Logger().Debug("Trying to read the previous version from the git tags first...")
	v, err := s.FromTagStrategy.ReadVersion()
	if err == nil {
		return v, nil
	}

	if err == fromtag.ErrNoTags {
		log.Logger().Debug("The git repository has no tags yet - returning fake version 0.0.0")
		return semver.MustParse("0.0.0"), nil
	}

	return nil, err
}

func (s Strategy) BumpVersion(previous semver.Version) (*semver.Version, error) {
	log.Logger().Debug("Trying to bump the version using semantic release first...")
	v, err := s.SemanticStrategy.BumpVersion(previous)
	if err == nil {
		return v, nil
	}

	if err == semantic.ErrPreviousVersionTagNotFound {
		log.Logger().Debugf("The git repository has no tag for the previous version %s - fallback to incrementing the patch component of the previous version", previous.String())
		next := previous.IncPatch()
		return &next, nil
	}

	return nil, err
}
