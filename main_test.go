package main

import (
	"context"
	"github.com/jenkins-x/jx-release-version/domain"
	"github.com/jenkins-x/jx-release-version/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakefile(t *testing.T) {

	c := config{
		dir: "test-resources/make",
	}

	v, err := getVersion(c)

	assert.NoError(t, err)

	assert.Equal(t, "1.2.0-SNAPSHOT", v, "error with getVersion for a Makefile")
}

func TestPomXML(t *testing.T) {

	c := config{
		dir: "test-resources/java",
	}
	v, err := getVersion(c)

	assert.NoError(t, err)

	assert.Equal(t, "1.0-SNAPSHOT", v, "error with getVersion for a pom.xml")
}

func TestChart(t *testing.T) {

	c := config{
		dir: "test-resources/helm",
	}
	v, err := getVersion(c)

	assert.NoError(t, err)

	assert.Equal(t, "0.0.1-SNAPSHOT", v, "error with getVersion for a pom.xml")
}

//func TestGetGithubTag(t *testing.T) {
//
//	c := config{
//		ghOwner:      "rawlingsj",
//		ghRepository: "test432317675",
//	}
//	v, err := getLatestTag(c)
//
//	assert.NoError(t, err)
//
//	assert.Equal(t, "2.0.0", v, "error with getLatestGithubTag for a Makefile")
//}

func TestGetGitTag(t *testing.T) {

	// first get the expeted version from github as test above passed
	c := config{
		ghOwner:      "rawlingsj",
		ghRepository: "semver-release-version",
	}

	tags := createTags()

	mockClient := &mocks.GitClient{}
	mockClient.On("ListTags", context.Background(), c.ghOwner, c.ghRepository).Return(tags, nil)
	//gitHubClient := NewGitHubClient(c)
	expectedVersion, err := getLatestTag(c, mockClient)
	assert.NoError(t, err)

	c = config{
		debug: true,
	}

	v, err := getLatestTag(c, nil)

	assert.NoError(t, err)

	assert.Equal(t, expectedVersion, v, "error with getLatestGithubTag for a Makefile")
}

//func TestGetNewVersionFromTag(t *testing.T) {
//
//	c := config{
//		dryrun:       false,
//		debug:        true,
//		dir:          "test-resources/make",
//		ghOwner:      "rawlingsj",
//		ghRepository: "test432317675",
//	}
//
//	v, err := getNewVersionFromTag(c)
//
//	assert.NoError(t, err)
//	assert.Equal(t, "2.0.1", v, "error bumping a patch version")
//}

func TestGetNewVersionFromTagCurrentRepo(t *testing.T) {

	c := config{
		dryrun: false,
		debug:  true,
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
		debug:        true,
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
		debug:        true,
	}

	tags := createTags()

	mockClient := &mocks.GitClient{}
	mockClient.On("ListTags", context.Background(), c.ghOwner, c.ghRepository).Return(tags, nil)

	v, err := getNewVersionFromTag(c, mockClient)

	assert.NoError(t, err)
	assert.Equal(t, "1.0.18", v, "error bumping a patch version")
}

func TestGetGitOwner(t *testing.T) {

	rs := getCurrentGitOwnerRepo("git@github.com:rawlingsj/semver-release-version.git")

	assert.Equal(t, "rawlingsj", rs[0])
	assert.Equal(t, "semver-release-version", rs[1])

	//rs = getCurrentGitOwnerRepo("https://github.com/rawlingsj/semver-release-number.git")

	//assert.Equal(t, "rawlingsj", rs[0])
	//assert.Equal(t, "semver-release-number", rs[1])

	//assertParseGitRepositoryInfo("git://host.xz/org/repo", "host.xz", "org", "repo");
	//assertParseGitRepositoryInfo("git://host.xz/org/repo.git", "host.xz", "org", "repo");
	//assertParseGitRepositoryInfo("git://host.xz/org/repo.git/", "host.xz", "org", "repo");
	//assertParseGitRepositoryInfo("git://github.com/jstrachan/npm-pipeline-test-project.git", "github.com", "jstrachan", "npm-pipeline-test-project");
	//assertParseGitRepositoryInfo("https://github.com/fabric8io/foo.git", "github.com", "fabric8io", "foo");
	//assertParseGitRepositoryInfo("https://github.com/fabric8io/foo", "github.com", "fabric8io", "foo");
	//assertParseGitRepositoryInfo("git@github.com:jstrachan/npm-pipeline-test-project.git", "github.com", "jstrachan", "npm-pipeline-test-project");
	//assertParseGitRepositoryInfo("git@github.com:bar/foo.git", "github.com", "bar", "foo");
	//assertParseGitRepositoryInfo("git@github.com:bar/foo", "github.com", "bar", "foo");
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
