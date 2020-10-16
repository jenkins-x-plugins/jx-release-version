package main

import (
	"context"
	"testing"

	"github.com/jenkins-x-plugins/jx-release-version/adapters"
	"github.com/jenkins-x-plugins/jx-release-version/domain"
	"github.com/jenkins-x-plugins/jx-release-version/mocks"

	"github.com/stretchr/testify/assert"
)

func TestMakefile(t *testing.T) {
	c := config{
		dir: "test-resources/make",
	}

	v, err := getVersion(c)

	assert.NoError(t, err)

	assert.Equal(t, "1.2.0-SNAPSHOT", v, "error with getVersion for a Makefile")
}

func TestAutomakefile(t *testing.T) {
	c := config{
		dir: "test-resources/automake",
	}

	v, err := getVersion(c)

	assert.NoError(t, err)

	assert.Equal(t, "1.2.0-SNAPSHOT", v, "error with getVersion for a configure.ac")
}

func TestCMakefile(t *testing.T) {

	c := config{
		dir: "test-resources/cmake",
	}

	v, err := getVersion(c)

	assert.NoError(t, err)

	assert.Equal(t, "1.2.0-SNAPSHOT", v, "error with getVersion for a CMakeLists.txt")
}

func TestPomXML(t *testing.T) {
	c := config{
		dir: "test-resources/java",
	}
	v, err := getVersion(c)

	assert.NoError(t, err)

	assert.Equal(t, "1.0-SNAPSHOT", v, "error with getVersion for a pom.xml")
}

func TestPackageJSON(t *testing.T) {
	c := config{
		dir: "test-resources/package",
	}
	v, err := getVersion(c)

	assert.NoError(t, err)

	assert.Equal(t, "1.2.3", v, "error with getVersion for a package.json")
}

// TODO enable this. It seems that meta-pipeline is bumping the version of the Chart.yaml
// when the release pipeline is running, this is causing this test to fail.
/*
func TestChart(t *testing.T) {
	c := config{
		dir: "test-resources/helm",
	}
	v, err := getVersion(c)

	assert.NoError(t, err)

	assert.Equal(t, "0.0.1-SNAPSHOT", v, "error with getVersion for a Chart.yaml")
}
*/

func TestSetupPyStandard(t *testing.T) {

	c := config{
		dir: "test-resources/python/standard",
	}
	v, err := getVersion(c)

	assert.NoError(t, err)

	assert.Equal(t, "4.5.6", v, "error with getVersion for a setup.py")
}

func TestSetupPyNested(t *testing.T) {

	c := config{
		dir: "test-resources/python/nested",
	}
	v, err := getVersion(c)

	assert.NoError(t, err)

	assert.Equal(t, "4.5.6", v, "error with getVersion for a setup.py")
}

func TestSetupPyOneLine(t *testing.T) {

	c := config{
		dir: "test-resources/python/one_line",
	}
	v, err := getVersion(c)

	assert.NoError(t, err)

	assert.Equal(t, "4.5.6", v, "error with getVersion for a setup.py")
}

func TestGetGitTag(t *testing.T) {
	c := config{
		ghOwner:      "jenkins-x",
		ghRepository: "jx-release-version",
	}

	gitHubClient := adapters.NewGitHubClient(c.debug)

	expectedVersion, err := getLatestTag(c, gitHubClient)
	assert.NoError(t, err)

	c = config{}

	v, err := getLatestTag(c, gitHubClient)

	assert.NoError(t, err)

	assert.Equal(t, expectedVersion, v, "error with getLatestTag for a Makefile")
}

func TestGetNewVersionFromTagCurrentRepo(t *testing.T) {
	c := config{
		dryrun: false,
		dir:    "test-resources/make",
	}

	tags := createTags()

	mockClient := &mocks.GitClient{}
	mockClient.On("ListTags", context.Background(), c.ghOwner, c.ghRepository).Return(tags, nil)
	v, err := getNewVersionFromTag(c, mockClient)

	assert.NoError(t, err)
	assert.Equal(t, "1.2.0", v, "error bumping a patch version")
}

func TestGetNewMinorVersionFromGitHubTag(t *testing.T) {

	c := config{
		ghOwner:      "rawlingsj",
		ghRepository: "semver-release-version",
		minor:        true,
	}

	tags := createTags()

	mockClient := &mocks.GitClient{}
	mockClient.On("ListTags", context.Background(), c.ghOwner, c.ghRepository).Return(tags, nil)

	v, err := getNewVersionFromTag(c, mockClient)

	assert.NoError(t, err)
	assert.Equal(t, "1.1.0", v, "error bumping a minor version")
}

func TestGetNewPatchVersionFromGitHubTag(t *testing.T) {

	c := config{
		ghOwner:      "rawlingsj",
		ghRepository: "semver-release-version",
	}

	tags := createTags()

	mockClient := &mocks.GitClient{}
	mockClient.On("ListTags", context.Background(), c.ghOwner, c.ghRepository).Return(tags, nil)

	v, err := getNewVersionFromTag(c, mockClient)

	assert.NoError(t, err)
	assert.Equal(t, "1.0.18", v, "error bumping a patch version")
}

func createTags() []domain.Tag {
	var tags []domain.Tag
	tags = append(tags, domain.Tag{Name: "v1.0.0"})
	tags = append(tags, domain.Tag{Name: "v1.0.1"})
	tags = append(tags, domain.Tag{Name: "v1.0.10"})
	tags = append(tags, domain.Tag{Name: "v1.0.11"})
	tags = append(tags, domain.Tag{Name: "v1.0.12"})
	tags = append(tags, domain.Tag{Name: "v1.0.13"})
	tags = append(tags, domain.Tag{Name: "v1.0.14"})
	tags = append(tags, domain.Tag{Name: "v1.0.15"})
	tags = append(tags, domain.Tag{Name: "v1.0.16"})
	tags = append(tags, domain.Tag{Name: "v1.0.17"})
	tags = append(tags, domain.Tag{Name: "v1.0.2"})
	tags = append(tags, domain.Tag{Name: "v1.0.3"})
	tags = append(tags, domain.Tag{Name: "v1.0.4"})
	tags = append(tags, domain.Tag{Name: "v1.0.5"})
	tags = append(tags, domain.Tag{Name: "v1.0.6"})
	tags = append(tags, domain.Tag{Name: "v1.0.7"})
	tags = append(tags, domain.Tag{Name: "v1.0.8"})
	tags = append(tags, domain.Tag{Name: "v1.0.9"})

	return tags
}
