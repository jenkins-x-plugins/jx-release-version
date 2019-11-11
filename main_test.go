package main

import (
	"testing"

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

func TestChart(t *testing.T) {
	c := config{
		dir: "test-resources/helm",
	}
	v, err := getVersion(c)

	assert.NoError(t, err)

	assert.Equal(t, "0.0.1-SNAPSHOT", v, "error with getVersion for a Chart.yaml")
}

func TestGetGitTag(t *testing.T) {
	c := config{
		ghOwner:      "jenkins-x",
		ghRepository: "jx-release-version",
	}
	expectedVersion, err := getLatestTag(c)
	assert.NoError(t, err)

	c = config{}
	v, err := getLatestTag(c)

	assert.NoError(t, err)

	assert.Equal(t, expectedVersion, v, "error with getLatestTag for a Makefile")
}

func TestGetNewVersionFromTagCurrentRepo(t *testing.T) {
	c := config{
		dryrun: false,
		dir:    "test-resources/make",
	}

	v, err := getNewVersionFromTag(c)

	assert.NoError(t, err)
	assert.Equal(t, "1.2.0", v, "error bumping a patch version")
}
