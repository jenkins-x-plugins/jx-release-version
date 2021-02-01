package increment

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBumpVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		componentToIncrement string
		previous             semver.Version
		expected             *semver.Version
		expectedErrorMsg     string
	}{
		{
			name:                 "increment major",
			componentToIncrement: "major",
			previous:             *semver.MustParse("1.2.3"),
			expected:             semver.MustParse("2.0.0"),
		},
		{
			name:                 "increment minor",
			componentToIncrement: "minor",
			previous:             *semver.MustParse("1.2.3"),
			expected:             semver.MustParse("1.3.0"),
		},
		{
			name:                 "increment patch",
			componentToIncrement: "patch",
			previous:             *semver.MustParse("1.2.3"),
			expected:             semver.MustParse("1.2.4"),
		},
		{
			name:                 "increment patch by default",
			componentToIncrement: "",
			previous:             *semver.MustParse("1.2.3"),
			expected:             semver.MustParse("1.2.4"),
		},
		{
			name:                 "case insensitive",
			componentToIncrement: "MiNoR",
			previous:             *semver.MustParse("1.2.3"),
			expected:             semver.MustParse("1.3.0"),
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			s := Strategy{
				ComponentToIncrement: test.componentToIncrement,
			}
			actual, err := s.BumpVersion(test.previous)
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
