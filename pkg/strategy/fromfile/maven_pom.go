package fromfile

import (
	"encoding/xml"
	"fmt"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"os"
	"os/exec"
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
	log.Logger().Debugf("using path " + os.Getenv("PATH"))
	path, err := exec.LookPath("mvn")
	if err != nil {
		log.Logger().Debugf("Maven does not appear to be installed, reading directly from %s", filePath)
		return r.readDirectlyFromPom(filePath)
	}

	log.Logger().Debugf("Maven is installed into path %s", path)

	cmd := exec.Command("mvn",
		"-f",
		filePath,
		"-B",            // batch mode
		"-ntp",          // do not display transfer output
		"-q",            // quiet
		"-DforceStdout", // force version to displayed
		"-Dexpression=project.version",
		"help:evaluate",
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("unable to evaluate project.version %s", err)
	}

	return string(out), nil
}

func (r MavenPOMVersionReader) readDirectlyFromPom(filePath string) (string, error) {
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
