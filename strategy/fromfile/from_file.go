package fromfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
)

type Strategy struct {
	Dir      string
	FilePath string
}

func (s Strategy) ReadVersion() (*semver.Version, error) {
	var (
		dir = s.Dir
		err error
	)
	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	var (
		reader   FileVersionReader
		filePath string
	)
	if len(s.FilePath) > 0 {
		filePath = filepath.Join(dir, s.FilePath)
		reader, err = s.getReader()
	} else {
		reader, filePath, err = s.autoDetect(dir)
	}
	if err != nil {
		return nil, err
	}

	log.Logger().Debugf("Reading version from file %s using reader %s", filePath, reader)
	version, err := reader.ReadFileVersion(filePath)
	if err != nil {
		return nil, err
	}

	log.Logger().Debugf("Found version %s", version)
	return semver.NewVersion(version)
}

func (s Strategy) BumpVersion(_ semver.Version) (*semver.Version, error) {
	return s.ReadVersion()
}

func (s Strategy) autoDetect(dir string) (FileVersionReader, string, error) {
	for _, reader := range fileVersionReaders {
		for _, fileName := range reader.SupportedFiles() {
			filePath := filepath.Join(dir, fileName)
			if _, err := os.Stat(filePath); err == nil {
				log.Logger().Debugf("Using file %s to read version", filePath)
				return reader, filePath, nil
			}
		}
	}

	return nil, "", fmt.Errorf("could not find a file to read version from, in directory %s", dir)
}

func (s Strategy) getReader() (FileVersionReader, error) {
	for _, reader := range fileVersionReaders {
		for _, fileName := range reader.SupportedFiles() {
			if strings.HasSuffix(s.FilePath, fileName) {
				return reader, nil
			}
		}
	}

	return nil, fmt.Errorf("could not find a file version reader for %s", s.FilePath)
}

type FileVersionReader interface {
	ReadFileVersion(filePath string) (string, error)
	SupportedFiles() []string
	String() string
}

// fileVersionReaders is an ordered list of all readers to try
// when auto-detecting the file to use
var fileVersionReaders = []FileVersionReader{
	HelmChartVersionReader{},
	MakefileVersionReader{},
	AutomakeVersionReader{},
	CMakeVersionReader{},
	PythonVersionReader{},
	MavenPOMVersionReader{},
	JsPackageVersionReader{},
	GradleVersionReader{},
}
