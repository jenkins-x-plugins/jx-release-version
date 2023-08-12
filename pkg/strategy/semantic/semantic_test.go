package semantic

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/mholt/archiver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBumpVersion(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		strategy         Strategy
		previous         semver.Version
		expected         *semver.Version
		expectedErrorMsg string
	}{
		{
			name: "feat commit",
			strategy: Strategy{
				Dir: "testdata/git-repo",
			},
			previous: *semver.MustParse("2.0.0"),
			expected: semver.MustParse("2.1.0"),
		},
		{
			name: "breaking change",
			strategy: Strategy{
				Dir: "testdata/git-repo",
			},
			previous: *semver.MustParse("1.1.0"),
			expected: semver.MustParse("2.0.0"),
		},
		{
			name: "feat from commit headline",
			strategy: Strategy{
				CommitHeadlinesString: "feat: a feature",
			},
			previous: *semver.MustParse("2.0.0"),
			expected: semver.MustParse("2.1.0"),
		},
		{
			name: "feat from commit headlines",
			strategy: Strategy{
				CommitHeadlinesString: `chore: a chore
feat: a feature`,
			},
			previous: *semver.MustParse("2.0.0"),
			expected: semver.MustParse("2.1.0"),
		},
		{
			name: "breaking change from commit headline",
			strategy: Strategy{
				CommitHeadlinesString: "feat!: a breaking feature",
			},
			previous: *semver.MustParse("1.1.0"),
			expected: semver.MustParse("2.0.0"),
		},
		{
			name: "breaking change from commit headlines",
			strategy: Strategy{
				CommitHeadlinesString: `chore: a chore
feat!: a breaking feature`,
			},
			previous: *semver.MustParse("1.1.0"),
			expected: semver.MustParse("2.0.0"),
		},
		{
			name: "patch from unrecognized commit headline",
			strategy: Strategy{
				CommitHeadlinesString: "nothing",
			},
			previous: *semver.MustParse("1.1.0"),
			expected: semver.MustParse("1.1.1"),
		},
	}

	// the git repo is stored as a tar.gz archive to make it easy to commit
	gitRepoPath := filepath.Join("testdata", "git-repo")
	err := os.RemoveAll(gitRepoPath)
	require.NoErrorf(t, err, "failed to delete %s", gitRepoPath)
	err = archiver.Unarchive(filepath.Join("testdata", "git-repo.tar.gz"), gitRepoPath)
	require.NoErrorf(t, err, "failed to decompress git repository at %s", gitRepoPath)

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
