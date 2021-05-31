package fromfile

import (
	"encoding/json"
	"os"
)

type JsPackageVersionReader struct {
}

func (r JsPackageVersionReader) String() string {
	return "javascript-package.json"
}

func (r JsPackageVersionReader) SupportedFiles() []string {
	return []string{
		"package.json",
	}
}

func (r JsPackageVersionReader) ReadFileVersion(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var pkg JsPackage
	err = json.NewDecoder(f).Decode(&pkg)
	if err != nil {
		return "", err
	}

	if pkg.Version == "" {
		return "", ErrFileHasNoVersion
	}

	return pkg.Version, nil
}

type JsPackage struct {
	Version string `json:"version"`
}
