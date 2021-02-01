package fromfile

import (
	"encoding/xml"
	"fmt"
	"os"
)

type MavenPOMVersionReader struct {
}

func (r MavenPOMVersionReader) String() string {
	return "maven POM"
}

func (r MavenPOMVersionReader) SupportedFiles() []string {
	return []string{
		"pom.xml",
	}
}

func (r MavenPOMVersionReader) ReadFileVersion(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var pom MavenPOM
	err = xml.NewDecoder(f).Decode(&pom)
	if err != nil {
		return "", err
	}

	if pom.Version == "" {
		return "", fmt.Errorf("version not found in file %s", filePath)
	}

	return pom.Version, nil
}

type MavenPOM struct {
	Version string `xml:"version"`
}
