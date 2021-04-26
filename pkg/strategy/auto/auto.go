package auto

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/jenkins-x-plugins/jx-release-version/v2/pkg/strategy/fromtag"
	"github.com/jenkins-x-plugins/jx-release-version/v2/pkg/strategy/semantic"
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

	if err == fromtag.ErrNoTags || err == fromtag.ErrNoSemverTags {
		log.Logger().Debugf("Using fake version 0.0.0 because %s", err)
		return semver.MustParse("0.0.0"), nil
	}

	return nil, fmt.Errorf("failed to read previous version from tags: %w", err)
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

	return nil, fmt.Errorf("failed to bump version %s using semantic strategy: %w", previous.String(), err)
}
