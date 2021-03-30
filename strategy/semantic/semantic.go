package semantic

import (
	"errors"
	"os"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/zbindenren/cc"
)

var (
	ErrPreviousVersionTagNotFound = errors.New("the git repository has no tag for the previous version")
)

type Strategy struct {
	Dir             string
	StripPrerelease bool
}

func (s Strategy) BumpVersion(previous semver.Version) (*semver.Version, error) {
	var (
		dir = s.Dir
		err error
	)
	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	repo, err := git.PlainOpen(dir)
	if err != nil {
		return nil, err
	}

	tagCommit, err := s.extractTagCommit(repo, previous.String())
	if err != nil {
		return nil, err
	}

	summary, err := s.parseCommitsSince(repo, tagCommit)
	if err != nil {
		return nil, err
	}

	if s.StripPrerelease {
		previous, err = previous.SetPrerelease("")
		if err != nil {
			return nil, err
		}
	}

	var version semver.Version
	switch {
	case summary.breakingChanges:
		log.Logger().Debug("Found breaking changes - incrementing major component")
		version = previous.IncMajor()
	case summary.types["feat"]:
		log.Logger().Debug("Found at least 1 new feature - incrementing minor component")
		version = previous.IncMinor()
	default:
		log.Logger().Debug("Incrementing patch component")
		version = previous.IncPatch()
	}

	return &version, nil
}

func (s Strategy) extractTagCommit(repo *git.Repository, tagName string) (*object.Commit, error) {
	var tagCommit *object.Commit

	previousTagRef, err := repo.Tag(tagName)
	if err == git.ErrTagNotFound {
		previousTagRef, err = repo.Tag("v" + tagName)
		if err == git.ErrTagNotFound {
			return nil, ErrPreviousVersionTagNotFound
		}
	}
	if err != nil {
		return nil, err
	}

	previousTag, err := repo.TagObject(previousTagRef.Hash())
	if errors.Is(err, plumbing.ErrObjectNotFound) {
		// it's a lightweight tag, not an annotated tag
		tagCommit, err = repo.CommitObject(previousTagRef.Hash())
	}
	if err != nil {
		return nil, err
	}

	if previousTag != nil {
		tagCommit, err = previousTag.Commit()
		if err != nil {
			return nil, err
		}
	}

	log.Logger().Debugf("Previous version tag commit is %s", tagCommit.Hash)
	return tagCommit, nil
}

type conventionalCommitsSummary struct {
	types           map[string]bool
	breakingChanges bool
}

func (s Strategy) parseCommitsSince(repo *git.Repository, firstCommit *object.Commit) (*conventionalCommitsSummary, error) {
	summary := conventionalCommitsSummary{
		types: map[string]bool{},
	}

	log.Logger().Debugf("Iterating over all commits since %s", firstCommit.Committer.When)
	commitIterator, err := repo.Log(&git.LogOptions{
		Since: &firstCommit.Committer.When,
	})
	if err != nil {
		return nil, err
	}

	err = commitIterator.ForEach(func(commit *object.Commit) error {
		if commit.Hash == firstCommit.Hash {
			log.Logger().Debugf("Skipping first commit %s and stopping iteration", commit.Hash)
			return storer.ErrStop
		}
		log.Logger().Debugf("Parsing commit %s", commit.Hash)
		c, err := cc.Parse(commit.Message)
		if err != nil {
			log.Logger().WithError(err).Debugf("Skipping non-conventional commit %s", commit.Hash)
			return nil
		}
		summary.types[c.Header.Type] = true
		if len(c.BreakingMessage()) > 0 {
			summary.breakingChanges = true
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	log.Logger().Debugf("Summary of conventional commits since %s: %#v", firstCommit.Committer.When, summary)
	return &summary, nil
}
