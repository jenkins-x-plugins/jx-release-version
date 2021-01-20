package semrel

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/jenkins-x/jx-helpers/v3/pkg/gitclient"
	"github.com/jenkins-x/jx-helpers/v3/pkg/gitclient/gitlog"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
)

var commitPattern = regexp.MustCompile("^(\\w*)(?:\\((.*)\\))?\\: (.*)$")
var breakingPattern = regexp.MustCompile("BREAKING CHANGES?")

type change struct {
	Major, Minor, Patch bool
}

type conventionalCommit struct {
	*gitlog.Commit
	MessageLines []string
	Type         string
	Scope        string
	MessageBody  string
	Change       change
}

type release struct {
	SHA     string
	Version *semver.Version
}

// GetNewVersion uses the conventional commits in the range of latestTagRev..endSha to increment the version from latestTag
func GetNewVersion(dir string, git gitclient.Interface) (*semver.Version, error) {
	latestCommitSHA, err := gitclient.GetLatestCommitSha(git, dir)
	if err != nil {
		return nil, fmt.Errorf("failed getting latest commit SHA in %s: %w", dir, err)
	}
	latestTagRev, latestTag, err := gitclient.GetCommitPointedToByLatestTag(git, dir)
	if err != nil {
		return nil, fmt.Errorf("failed getting latest tag in %s: %w", dir, err)
	}
	version, err := semver.NewVersion(strings.TrimPrefix(latestTag, "v"))
	if err != nil {
		return nil, fmt.Errorf("failed parsing %s as semantic version: %w", latestTag, err)
	}
	release := release{
		SHA:     latestTagRev,
		Version: version,
	}

	out, err := git.Command(dir, "--no-pager", "log", fmt.Sprintf("%s..%s", latestTagRev, latestCommitSHA), "--reverse", "--decorate=no", "--no-color")
	if err != nil {
		return nil, fmt.Errorf("failed getting commits in range %s..%s: %w", latestTagRev, latestCommitSHA, err)
	}

	commits := make([]*conventionalCommit, 0)
	gitCommits := gitlog.ParseGitLog(out)
	log.Logger().Debugf("got %d commits", len(gitCommits))
	for _, c := range gitCommits {
		commit := c
		commits = append(commits, parseCommit(commit))
	}

	return applyChange(release.Version, calculateChange(commits, &release)), nil
}

func calculateChange(commits []*conventionalCommit, latestRelease *release) change {
	var change change
	for _, commit := range commits {
		if latestRelease.SHA == commit.SHA {
			break
		}
		change.Major = change.Major || commit.Change.Major
		change.Minor = change.Minor || commit.Change.Minor
		change.Patch = change.Patch || commit.Change.Patch
	}
	return change
}

func applyChange(version *semver.Version, change change) *semver.Version {
	log.Logger().Debugf("applying change %+v", change)
	if version.Major() == 0 {
		change.Major = true
	}
	if !change.Major && !change.Minor && !change.Patch {
		log.Logger().Debugf("unable to determine change so defaulting to patch")
		change.Patch = true
	}
	var newVersion semver.Version
	preRel := version.Prerelease()
	if preRel == "" {
		switch {
		case change.Major:
			newVersion = version.IncMajor()
			break
		case change.Minor:
			newVersion = version.IncMinor()
			break
		case change.Patch:
			newVersion = version.IncPatch()
			break
		}
		return &newVersion
	}
	preRelVer := strings.Split(preRel, ".")
	if len(preRelVer) > 1 {
		idx, err := strconv.ParseInt(preRelVer[1], 10, 32)
		if err != nil {
			idx = 0
		}
		preRel = fmt.Sprintf("%s.%d", preRelVer[0], idx+1)
	} else {
		preRel += ".1"
	}
	newVersion, _ = version.SetPrerelease(preRel)
	return &newVersion
}

func parseCommit(commit *gitlog.Commit) *conventionalCommit {
	c := &conventionalCommit{
		Commit: commit,
	}
	log.Logger().Debugf("parsing message '%s'", commit.Comment)
	c.MessageLines = strings.Split(commit.Comment, "\n")
	found := commitPattern.FindAllStringSubmatch(c.MessageLines[0], -1)
	if len(found) < 1 {
		return c
	}
	c.Type = strings.ToLower(found[0][1])
	c.Scope = found[0][2]
	c.MessageBody = found[0][3]
	c.Change = change{
		Major: breakingPattern.MatchString(commit.Comment),
		Minor: c.Type == "feat",
		Patch: c.Type == "fix",
	}
	return c
}
