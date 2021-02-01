package auto

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/jenkins-x-plugins/jx-release-version/v2/strategy/fromtag"
	"github.com/jenkins-x-plugins/jx-release-version/v2/strategy/semantic"
	"github.com/mholt/archiver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadVersion(t *testing.T) {
	tests := []struct {
		name             string
		strategy         Strategy
		expected         *semver.Version
		expectedErrorMsg string
	}{
		{
			name: "empty git repo",
			strategy: Strategy{
				FromTagStrategy: fromtag.Strategy{
					Dir: "testdata/empty-git-repo",
				},
			},
			expected: semver.MustParse("0.0.0"),
		},
		{
			name: "non-empty git repo",
			strategy: Strategy{
				FromTagStrategy: fromtag.Strategy{
					Dir: "testdata/git-repo",
				},
			},
			expected: semver.MustParse("v2.0.0"),
		},
	}

	setupGitRepos(t)

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual, err := test.strategy.ReadVersion()
			if len(test.expectedErrorMsg) > 0 {
				require.EqualError(t, err, test.expectedErrorMsg)
				assert.Nil(t, actual)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expected, actual)
			}
		})
	}
}

func TestBumpVersion(t *testing.T) {
	tests := []struct {
		name             string
		strategy         Strategy
		previous         semver.Version
		expected         *semver.Version
		expectedErrorMsg string
	}{
		{
			name: "empty git repo",
			strategy: Strategy{
				SemanticStrategy: semantic.Strategy{
					Dir: "testdata/empty-git-repo",
				},
			},
			previous: *semver.MustParse("1.0.0"),
			expected: semver.MustParse("1.0.1"),
		},
		{
			name: "non-empty git repo",
			strategy: Strategy{
				SemanticStrategy: semantic.Strategy{
					Dir: "testdata/git-repo",
				},
			},
			previous: *semver.MustParse("1.0.0"),
			expected: semver.MustParse("2.0.0"),
		},
	}

	setupGitRepos(t)

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual, err := test.strategy.BumpVersion(test.previous)
			if len(test.expectedErrorMsg) > 0 {
				require.EqualError(t, err, test.expectedErrorMsg)
				assert.Nil(t, actual)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expected, actual)
			}
		})
	}
}

func setupGitRepos(t *testing.T) {
	// the git repos are stored as a tar.gz archive to make it easy to commit
	for _, repoName := range []string{"git-repo", "empty-git-repo"} {
		gitRepoPath := filepath.Join("testdata", repoName)
		err := os.RemoveAll(gitRepoPath)
		require.NoErrorf(t, err, "failed to delete %s", gitRepoPath)
		err = archiver.Unarchive(filepath.Join("testdata", fmt.Sprintf("%s.tar.gz", repoName)), gitRepoPath)
		require.NoErrorf(t, err, "failed to decompress git repository at %s", gitRepoPath)
	}
}
