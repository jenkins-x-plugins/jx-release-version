package fromfile

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPythonVersionReader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		filePath         string
		expected         string
		expectedErrorMsg string
	}{
		{
			name:     "standard file",
			filePath: "setup.py",
			expected: "1.2.11",
		},
		{
			name:     "nested file",
			filePath: "setup-nested.py",
			expected: "1.2.12",
		},
		{
			name:     "oneline file",
			filePath: "setup-oneline.py",
			expected: "1.2.13",
		},
		{
			name:     "double quotes",
			filePath: "setup-double-quotes.py",
			expected: "1.2.14",
		},
		{
			name:             "file does not exists",
			filePath:         "does-not-exists.yaml",
			expectedErrorMsg: "open testdata/does-not-exists.yaml: no such file or directory",
		},
		{
			name:             "invalid file",
			filePath:         "Chart.yaml",
			expectedErrorMsg: "setup call not found in file testdata/Chart.yaml",
		},
	}

	reader := PythonVersionReader{}
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
