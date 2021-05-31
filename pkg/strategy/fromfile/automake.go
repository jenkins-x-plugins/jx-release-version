package fromfile

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

var (
	configureRegexp = regexp.MustCompile(`AC_INIT\s*\(([^\s]+),\s*([.\d]+(-\w+)?).*\)`)
)

type AutomakeVersionReader struct {
}

func (r AutomakeVersionReader) String() string {
	return "automake"
}

func (r AutomakeVersionReader) SupportedFiles() []string {
	return []string{
		"configure.ac",
	}
}

func (r AutomakeVersionReader) ReadFileVersion(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "AC_INIT") {
			matched := configureRegexp.FindStringSubmatch(scanner.Text())
			if len(matched) < 3 {
				continue
			}

			v := strings.TrimSpace(matched[2])
			if v != "" {
				return v, nil
			}
		}
	}

	return "", ErrFileHasNoVersion
}
