package fromfile

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

var (
	gradleRegexp = regexp.MustCompile(`^version\s*=\s*['"]([.\d]+(-\w+)?)['"]`)
)

type GradleVersionReader struct {
}

func (r GradleVersionReader) String() string {
	return "gradle"
}

func (r GradleVersionReader) SupportedFiles() []string {
	return []string{
		"build.gradle",      // groovy syntax
		"build.gradle.kts",  // kotlin syntax
		"gradle.properties", // gradle properties syntax
	}
}

func (r GradleVersionReader) ReadFileVersion(filePath string) (string, error) {
	f, err := os.Open(filePath) // #nosec G304 -- user-provided version file path
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "version") {
			matched := gradleRegexp.FindStringSubmatch(scanner.Text())
			if len(matched) < 2 {
				continue
			}

			v := strings.TrimSpace(matched[1])
			if v != "" {
				return v, nil
			}
		}
	}

	return "", ErrFileHasNoVersion
}
