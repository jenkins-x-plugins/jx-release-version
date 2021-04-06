package fromfile

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type HelmChartVersionReader struct {
}

func (r HelmChartVersionReader) String() string {
	return "helm-chart"
}

func (r HelmChartVersionReader) SupportedFiles() []string {
	return []string{
		"Chart.yaml",
	}
}

func (r HelmChartVersionReader) ReadFileVersion(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var chart HelmChart
	err = yaml.NewDecoder(f).Decode(&chart)
	if err != nil {
		return "", err
	}

	if chart.Version == "" {
		return "", fmt.Errorf("version not found in file %s", filePath)
	}

	return chart.Version, nil
}

type HelmChart struct {
	Version string `yaml:"version"`
}
