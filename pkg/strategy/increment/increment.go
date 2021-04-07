package increment

import (
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
)

type Strategy struct {
	ComponentToIncrement string
}

func (s Strategy) BumpVersion(previous semver.Version) (*semver.Version, error) {
	var next semver.Version
	switch strings.ToLower(s.ComponentToIncrement) {
	case "major":
		log.Logger().Debug("Incrementing major component")
		next = previous.IncMajor()
	case "minor":
		log.Logger().Debug("Incrementing minor component")
		next = previous.IncMinor()
	default:
		log.Logger().Debug("Incrementing patch component")
		next = previous.IncPatch()
	}
	return &next, nil
}
