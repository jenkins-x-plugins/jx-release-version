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

func TestGetGithubTag(t *testing.T) {

	c := config{
		owner: "rawlingsj",
		repo:  "test432317675",
	}
	v, err := getLatestGithubTag(c)

	assert.NoError(t, err)

	assert.Equal(t, "2.0.0", v, "error with getVersion for a Makefile")
}

func TestGetNewVersionFromTag(t *testing.T) {

	c := config{
		dryrun: false,
		dir:    "test-resources/make",
		owner:  "rawlingsj",
		repo:   "test432317675",
	}

	v, err := getNewVersionFromTag(c)

	assert.NoError(t, err)
	assert.Equal(t, "2.0.1", v, "error bumping a patch version")
}

func TestGetNewVersionFromTagCurrentRepo(t *testing.T) {

	c := config{
		dryrun: false,
		dir:    "test-resources/make",
	}

	v, err := getNewVersionFromTag(c)

	assert.NoError(t, err)
	assert.Equal(t, "1.0.0", v, "error bumping a patch version")
}
