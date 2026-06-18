package fromfile

import (
	"bufio"
	"os"
	"strings"
)

type MakefileVersionReader struct {
}

func (r MakefileVersionReader) String() string {
	return "makefile"
}

func (r MakefileVersionReader) SupportedFiles() []string {
	return []string{
		"Makefile",
	}
}

func (r MakefileVersionReader) ReadFileVersion(filePath string) (string, error) {
	f, err := os.Open(filePath) // #nosec G304 -- user-provided version file path
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "VERSION") || strings.HasPrefix(scanner.Text(), "VERSION ") || strings.HasPrefix(scanner.Text(), "VERSION:") || strings.HasPrefix(scanner.Text(), "VERSION=") {
			parts := strings.Split(scanner.Text(), "=")
			if len(parts) < 2 {
				continue
			}

			v := strings.TrimSpace(parts[1])
			if v != "" {
				return v, nil
			}
		}
	}

	return "", ErrFileHasNoVersion
}
