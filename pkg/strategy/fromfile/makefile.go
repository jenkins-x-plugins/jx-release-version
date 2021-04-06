package fromfile

import (
	"bufio"
	"fmt"
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
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

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

	return "", fmt.Errorf("version not found in file %s", filePath)
}
