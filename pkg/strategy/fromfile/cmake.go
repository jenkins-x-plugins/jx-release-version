package fromfile

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	cmakeRegexp = regexp.MustCompile(`project\s*(([^\s]+)\s+VERSION\s+([.\d]+(-\w+)?).*)`)
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
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

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

	return "", fmt.Errorf("version not found in file %s", filePath)
}
