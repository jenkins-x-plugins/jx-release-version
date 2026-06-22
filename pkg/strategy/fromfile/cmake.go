package fromfile

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

var (
	cmakeRegexp = regexp.MustCompile(`project\s*((\S+)\s+VERSION\s+([.\d]+(-\w+)?).*)`)
)

type CMakeVersionReader struct {
}

func (r CMakeVersionReader) String() string {
	return "cmake"
}

func (r CMakeVersionReader) SupportedFiles() []string {
	return []string{
		"CMakeLists.txt",
	}
}

func (r CMakeVersionReader) ReadFileVersion(filePath string) (string, error) {
	f, err := os.Open(filePath) // #nosec G304 -- user-provided version file path
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), " VERSION ") {
			matched := cmakeRegexp.FindStringSubmatch(scanner.Text())
			if len(matched) < 4 {
				continue
			}
			v := strings.TrimSpace(matched[3])
			if v != "" {
				return v, nil
			}
		}
	}

	return "", ErrFileHasNoVersion
}
