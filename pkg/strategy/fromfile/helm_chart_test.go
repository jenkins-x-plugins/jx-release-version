package fromfile

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelmChartVersionReader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		filePath         string
		expected         string
		expectedErrorMsg string
	}{
		{
			name:     "valid file",
			filePath: "Chart.yaml",
			expected: "1.2.3",
		},
		{
			name:             "file does not exists",
			filePath:         "does-not-exists.yaml",
			expectedErrorMsg: "open testdata/does-not-exists.yaml: no such file or directory",
		},
		{
			name:             "invalid file",
			filePath:         "setup.py",
			expectedErrorMsg: "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `from se...` into fromfile.HelmChart",
		},
	}

	reader := HelmChartVersionReader{}
	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual, err := reader.ReadFileVersion(filepath.Join("testdata", test.filePath))
			if len(test.expectedErrorMsg) > 0 {
				require.EqualError(t, err, test.expectedErrorMsg)
				assert.Empty(t, actual)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expected, actual)
			}
		})
	}
}
