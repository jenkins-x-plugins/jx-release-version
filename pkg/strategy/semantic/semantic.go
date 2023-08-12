package semantic

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/zbindenren/cc"
)

var (
	ErrPreviousVersionTagNotFound = errors.New("the git repository has no tag for the previous version")
)

type Strategy struct {
	Dir                   string
	StripPrerelease       bool
	CommitHeadlinesString string
}

func (s Strategy) BumpVersion(previous semver.Version) (*semver.Version, error) {
	var (
		dir                   = s.Dir
		err                   error
		commitHeadlinesString = s.CommitHeadlinesString
		summary               *conventionalCommitsSummary
	)
	if commitHeadlinesString != "" {
		summary, err = s.parseCommitHeadlines(commitHeadlinesString)
		if err != nil {
			return nil, err
		}
	} else {
		if dir == "" {
			dir, err = os.Getwd()
			if err != nil {
				return nil, fmt.Errorf("failed to get current working directory: %w", err)
			}
		}

		repo, err := git.PlainOpen(dir)
		if err != nil {
			return nil, fmt.Errorf("failed to open git repository at %q: %w", dir, err)
		}

		tagCommit, err := s.extractTagCommit(repo, previous.String())
		if err != nil {
			return nil, err
		}

		summary, err = s.parseCommitsSince(repo, tagCommit)
		if err != nil {
			return nil, err
		}
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
		// let's try to prepend a `v` prefix...
		tagName = "v" + tagName
		previousTagRef, err = repo.Tag(tagName)
		if err == git.ErrTagNotFound {
			return nil, ErrPreviousVersionTagNotFound
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tag %q: %w", tagName, err)
	}

	previousTag, err := repo.TagObject(previousTagRef.Hash())
	if errors.Is(err, plumbing.ErrObjectNotFound) {
		// it's a lightweight tag, not an annotated tag
		tagCommit, err = repo.CommitObject(previousTagRef.Hash())
		if err != nil {
			return nil, fmt.Errorf("failed to get the commit with hash %q (from tag %q): %w", previousTagRef.Hash().String(), tagName, err)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get the annotated tag with hash %q (from tag name %q): %w", previousTagRef.Hash().String(), tagName, err)
	}

	if previousTag != nil {
		tagCommit, err = previousTag.Commit()
		if err != nil {
			return nil, fmt.Errorf("failed to get the commit from the annotated tag %q with hash %q: %w", previousTag.Name, previousTag.Hash.String(), err)
		}
	}

	if tagCommit == nil {
		return nil, fmt.Errorf("could not find a commit for tag %q", tagName)
	}

	log.Logger().Debugf("Previous version tag commit is %s", tagCommit.Hash)
	return tagCommit, nil
}

type conventionalCommitsSummary struct {
	conventionalCommitsCount int
	types                    map[string]bool
	breakingChanges          bool
}

func (s Strategy) parseCommitsSince(repo *git.Repository, firstCommit *object.Commit) (*conventionalCommitsSummary, error) {
	summary := conventionalCommitsSummary{
		types: map[string]bool{},
	}

	log.Logger().Debugf("Iterating over all commits since %s", firstCommit.Committer.When)
	commitIterator, err := repo.Log(&git.LogOptions{
		Since: &firstCommit.Committer.When,
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list commits since %s (commit %q): %w", firstCommit.Committer.When, firstCommit.Hash.String(), err)
	}
	defer commitIterator.Close()

	for {
		commit, err := commitIterator.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Logger().WithError(err).Debug("Skipping unretrievable commit")
			continue
		}

		log.Logger().Debugf("Checking commit %s with message %s", commit.Hash, commit.Message)
		if commit.Hash == firstCommit.Hash {
			log.Logger().Debugf("Found first commit %s and stopping iteration", commit.Hash)
			break
		}

		log.Logger().Debugf("Parsing commit %s", commit.Hash)
		c, err := cc.Parse(commit.Message)
		if err != nil {
			log.Logger().WithError(err).Debugf("Skipping non-conventional commit %s", commit.Hash)
			continue
		}

		summary.conventionalCommitsCount++
		summary.types[c.Header.Type] = true
		if len(c.BreakingMessage()) > 0 {
			summary.breakingChanges = true
		}
	}

	log.Logger().Debugf("Summary of conventional commits since %s: %#v", firstCommit.Committer.When, summary)
	return &summary, nil
}

func (s Strategy) parseCommitHeadlines(commitHeadlinesString string) (*conventionalCommitsSummary, error) {
	summary := conventionalCommitsSummary{
		types: map[string]bool{},
	}

	log.Logger().Debugf("Iterating over all commits headline passed as a string")

	commitHeadlines := regexp.MustCompile("\r?\n").Split(commitHeadlinesString, -1)

	for index, commitHeadline := range commitHeadlines {
		log.Logger().Debugf("Parsing commit headline number %d with message %s", index, commitHeadline)
		c, err := cc.Parse(commitHeadline)
		if err != nil {
			log.Logger().WithError(err).Debugf("Skipping non-conventional commit headline number %d", index)
			continue
		}

		summary.conventionalCommitsCount++
		summary.types[c.Header.Type] = true
		if len(c.BreakingMessage()) > 0 {
			summary.breakingChanges = true
		}
	}

	log.Logger().Debugf("Summary of conventional commits: %#v", summary)
	return &summary, nil
}
