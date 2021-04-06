package fromfile

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

var (
	// pythonRegexp is used to find the call to `setup(..., version='1.2.3', ...)`
	pythonSetupRegexp = regexp.MustCompile(`setup\((.|\n)*version\s*=\s*'(\d|\.)*'([^\)]|\n)*\)`)
	// pythonVersionRegexp is used to find the argument `version='1.2.3'`
	pythonVersionRegexp = regexp.MustCompile(`version\s*=\s*'(\d*|\.)*'`)
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
	content, err := ioutil.ReadFile(filePath)
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
	if v == "" {
		return "", fmt.Errorf("empty version found in file %s", filePath)
	}

	return v, nil
}
