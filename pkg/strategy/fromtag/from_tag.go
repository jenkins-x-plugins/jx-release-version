package fromtag

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
)

var (
	ErrNoTags       = errors.New("the git repository has no tags")
	ErrNoSemverTags = errors.New("the git repository has no semver tags")
)

type Strategy struct {
	Dir        string
	TagPattern string
	FetchTags  bool
}

func (s Strategy) ReadVersion() (*semver.Version, error) {
	var (
		dir = s.Dir
		err error
	)
	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current working directory: %w", err)
		}
	}

	var tagRegexp *regexp.Regexp
	if len(s.TagPattern) > 0 {
		tagRegexp, err = regexp.Compile(s.TagPattern)
		if err != nil {
			return nil, fmt.Errorf("failed to compile tag pattern %q: %w", s.TagPattern, err)
		}
	}

	repo, err := git.PlainOpen(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository at %q: %w", dir, err)
	}

	if s.FetchTags {
		log.Logger().Debug("Fetching tags from origin")
		err = repo.Fetch(&git.FetchOptions{
			RemoteName: "origin",
			Progress:   os.Stdout,
			RefSpecs:   []config.RefSpec{config.RefSpec("refs/tags/*:refs/tags/*")},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to fetch tags from origin at %q: %w", dir, err)
		}
	}

	tagIterator, err := repo.Tags()
	if err != nil {
		return nil, fmt.Errorf("failed to list tags from git repository at %q: %w", dir, err)
	}

	var (
		tags     int
		versions []semver.Version
	)
	err = tagIterator.ForEach(func(ref *plumbing.Reference) error {
		tags++
		tag := ref.Name().Short()
		if tagRegexp != nil && !tagRegexp.MatchString(tag) {
			log.Logger().Debugf("Skipping tag %q not matching pattern %q", tag, s.TagPattern)
			return nil
		}
		v, err := semver.NewVersion(tag)
		if err != nil {
			log.Logger().Debugf("Skipping non-semver tag %q (%s)", tag, err)
			return nil
		}
		versions = append(versions, *v)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterator over tags from git repository at %q: %w", dir, err)
	}
	if tags == 0 {
		return nil, ErrNoTags
	}
	if len(versions) == 0 && len(s.TagPattern) == 0 {
		return nil, ErrNoSemverTags
	}
	if len(versions) == 0 {
		return nil, fmt.Errorf("no semver tags with pattern %q found", s.TagPattern)
	}
	log.Logger().Debugf("Found %d semver tags with pattern %q", len(versions), s.TagPattern)

	sort.SliceStable(versions, func(i, j int) bool {
		return versions[i].GreaterThan(&versions[j])
	})

	return &versions[0], nil
}
