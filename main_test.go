package main

import (
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

func TestGetGithubTag(t *testing.T) {

	c := config{
		ghOwner:      "rawlingsj",
		ghRepository: "test432317675",
	}
	v, err := getLatestTag(c)

	assert.NoError(t, err)

	assert.Equal(t, "2.0.0", v, "error with getLatestGithubTag for a Makefile")
}

func TestGetGitTag(t *testing.T) {

	// first get the expeted version from github as test above passed
	c := config{
		ghOwner:      "rawlingsj",
		ghRepository: "semver-release-version",
	}
	expectedVersion, err := getLatestTag(c)
	assert.NoError(t, err)

	c = config{
		debug: true,
	}
	v, err := getLatestTag(c)

	assert.NoError(t, err)

	assert.Equal(t, expectedVersion, v, "error with getLatestGithubTag for a Makefile")
}

func TestGetNewVersionFromTag(t *testing.T) {

	c := config{
		dryrun:       false,
		debug:        true,
		dir:          "test-resources/make",
		ghOwner:      "rawlingsj",
		ghRepository: "test432317675",
	}

	v, err := getNewVersionFromTag(c)

	assert.NoError(t, err)
	assert.Equal(t, "2.0.1", v, "error bumping a patch version")
}

func TestGetNewVersionFromTagCurrentRepo(t *testing.T) {

	c := config{
		dryrun: false,
		debug:  true,
		dir:    "test-resources/make",
	}

	v, err := getNewVersionFromTag(c)

	assert.NoError(t, err)
	assert.Equal(t, "1.2.0", v, "error bumping a patch version")
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
