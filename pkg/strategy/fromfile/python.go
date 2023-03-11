package fromfile

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	// pythonRegexp is used to find the call to `setup(..., version='1.2.3', ...)`
	pythonSetupRegexp = regexp.MustCompile(`setup\((.|\n)*version\s*=\s*['"]{1}(\d|\.)*['"]{1}([^\)]|\n)*\)`)
	// pythonVersionRegexp is used to find the argument `version='1.2.3'` or `version="1.2.3"`
	pythonVersionRegexp = regexp.MustCompile(`version\s*=\s*['"]{1}(\d*|\.)*['"]{1}`)
)

type PythonVersionReader struct {
}

func (r PythonVersionReader) String() string {
	return "python"
}

func (r PythonVersionReader) SupportedFiles() []string {
	return []string{
		"setup.py",
	}
}

func (r PythonVersionReader) ReadFileVersion(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	setupCall := pythonSetupRegexp.Find(content)
	if len(setupCall) == 0 {
		return "", fmt.Errorf("setup call not found in file %s", filePath)
	}

	version := string(pythonVersionRegexp.Find(setupCall))

	parts := strings.Split(strings.Replace(version, " ", "", -1), "=")
	if len(parts) < 2 {
		return "", fmt.Errorf("version value not found in file %s", filePath)
	}

	v := strings.TrimPrefix(strings.TrimSuffix(parts[1], "'"), "'")
	v = strings.TrimPrefix(strings.TrimSuffix(v, "\""), "\"")
	if v == "" {
		return "", fmt.Errorf("empty version found in file %s", filePath)
	}

	return v, nil
}
