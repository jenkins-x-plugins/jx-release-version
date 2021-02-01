package fromfile

import (
	"path/filepath"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		strategy         Strategy
		expected         *semver.Version
		expectedErrorMsg string
	}{
		{
			name: "auto detect",
			strategy: Strategy{
				Dir: "testdata",
			},
			expected: semver.MustParse("1.2.3"),
		},
		{
			name: "Helm Chart",
			strategy: Strategy{
				Dir:      "testdata",
				FilePath: "Chart.yaml",
			},
			expected: semver.MustParse("1.2.3"),
		},
		{
			name: "Helm Chart with dir in filePath",
			strategy: Strategy{
				Dir:      "",
				FilePath: filepath.Join("testdata", "Chart.yaml"),
			},
			expected: semver.MustParse("1.2.3"),
		},
		{
			name: "Makefile",
			strategy: Strategy{
				Dir:      "testdata",
				FilePath: "Makefile",
			},
			expected: semver.MustParse("1.2.4"),
		},
		{
			name: "Automake",
			strategy: Strategy{
				Dir:      "testdata",
				FilePath: "configure.ac",
			},
			expected: semver.MustParse("1.2.5"),
		},
		{
			name: "CMake",
			strategy: Strategy{
				Dir:      "testdata",
				FilePath: "CMakeLists.txt",
			},
			expected: semver.MustParse("1.2.6"),
		},
		{
			name: "Python",
			strategy: Strategy{
				Dir:      "testdata",
				FilePath: "setup.py",
			},
			expected: semver.MustParse("1.2.11"),
		},
		{
			name: "Maven POM",
			strategy: Strategy{
				Dir:      "testdata",
				FilePath: "pom.xml",
			},
			expected: semver.MustParse("1.2.9"),
		},
		{
			name: "Javascript package.json",
			strategy: Strategy{
				Dir:      "testdata",
				FilePath: "package.json",
			},
			expected: semver.MustParse("1.2.10"),
		},
		{
			name: "Gradle (groovy)",
			strategy: Strategy{
				Dir:      "testdata",
				FilePath: "build.gradle",
			},
			expected: semver.MustParse("1.2.7"),
		},
		{
			name: "Gradle (kotlin)",
			strategy: Strategy{
				Dir:      "testdata",
				FilePath: "build.gradle.kts",
			},
			expected: semver.MustParse("1.2.8"),
		},
		{
			name: "unknown file",
			strategy: Strategy{
				Dir:      "testdata",
				FilePath: "something.else",
			},
			expectedErrorMsg: "could not find a file version reader for something.else",
		},
	}

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

func TestAutoDetect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		dir              string
		expectedReader   FileVersionReader
		expectedFilePath string
		expectedErrorMsg string
	}{
		{
			name:             "testdata",
			dir:              "testdata",
			expectedReader:   HelmChartVersionReader{},
			expectedFilePath: "testdata/Chart.yaml",
		},
		{
			name:             "no match",
			dir:              ".",
			expectedReader:   nil,
			expectedErrorMsg: "could not find a file to read version from, in directory .",
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			s := Strategy{
				Dir: test.dir,
			}

			actualReader, actualFilePath, err := s.autoDetect(test.dir)
			if len(test.expectedErrorMsg) > 0 {
				require.EqualError(t, err, test.expectedErrorMsg)
				assert.Nil(t, actualReader)
				assert.Empty(t, actualFilePath)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectedReader, actualReader)
				assert.Equal(t, test.expectedFilePath, actualFilePath)
			}
		})
	}

}
