package main

import (
	"errors"
	"fmt"
	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
	version "github.com/hashicorp/go-version"
	"io/ioutil"
	"os"
	"strings"

	"bufio"
	"context"
	"encoding/xml"
	"flag"
	"golang.org/x/oauth2"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
)

type Project struct {
	Version string `xml:"version"`
}

type config struct {
	dryrun       bool
	debug        bool
	dir          string
	ghOwner      string
	ghRepository string
}

func main() {

	debug := flag.Bool("debug", false, "prints debug into to console")
	dir := flag.String("folder", ".", "the folder to look for files that contain a pom.xml or Makfile with the project version to bump")
	owner := flag.String("gh-owner", "", "a github repository owner if not running from within a git project  e.g. fabric8io")
	repo := flag.String("gh-repository", "", "a git repository if not running from within a git project  e.g. fabric8")

	flag.Parse()

	c := config{
		debug:        *debug,
		dir:          *dir,
		ghOwner:      *owner,
		ghRepository: *repo,
	}

	v, err := getNewVersionFromTag(c)
	if err != nil {
		fmt.Println("failed to get new version", err)
		os.Exit(-1)
	}
	fmt.Print(fmt.Sprintf("%s", v))
}

func getVersion(c config) (string, error) {
	m, err := ioutil.ReadFile(c.dir + string(filepath.Separator) + "Makefile")
	if err == nil {
		if c.debug {
			fmt.Println("Found Makefile")
		}
		scanner := bufio.NewScanner(strings.NewReader(string(m)))
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "VERSION") {
				parts := strings.Split(scanner.Text(), "=")

				v := strings.TrimSpace(parts[1])
				if v != "" {
					if c.debug {
						fmt.Println(fmt.Sprintf("existing Makefile version %v", v))
					}
					return v, nil
				}
			}
		}
	}

	p, err := ioutil.ReadFile(c.dir + string(filepath.Separator) + "pom.xml")
	if err == nil {
		if c.debug {
			fmt.Println("found pom.xml")
		}
		var project Project
		xml.Unmarshal(p, &project)
		if project.Version != "" {
			if c.debug {
				fmt.Println(fmt.Sprintf("existing version %v", project.Version))
			}
			return project.Version, nil
		}
	}
	return "", errors.New("no recognised file to obtain current version from")
}

func getLatestTag(c config) (string, error) {
	// if repo isn't provided by flags fall back to using current repo if run from a git project
	var versionsRaw []string
	if c.ghOwner != "" && c.ghRepository != "" {
		token := os.Getenv("GITHUB_AUTH_TOKEN")
		ctx := context.Background()
		var client *github.Client
		if token != "" {

			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token},
			)
			tc := oauth2.NewClient(ctx, ts)

			client = github.NewClient(tc)
		} else {
			if c.debug {
				fmt.Println("no GITHUB_AUTH_TOKEN env var found so using unauthenticated request")
			}
			client = github.NewClient(nil)
		}

		tags, _, err := client.Repositories.ListTags(ctx, c.ghOwner, c.ghRepository, nil)

		if err != nil {
			return "", err
		}
		if len(tags) == 0 {
			// if no current flags exist then lets start at 0.0.0
			return "0.0.0", errors.New("No existing tags found")
		}

		// build an array of all the tags
		versionsRaw = make([]string, len(tags))
		for i, tag := range tags {
			if c.debug {
				fmt.Println(fmt.Sprintf("found tag %s", tag.GetName()))
			}
			versionsRaw[i] = tag.GetName()
		}
	} else {
		_, err := exec.LookPath("git")
		if err != nil {
			return "", errors.New(fmt.Sprint("error running git: %v", err))
		}
		exec.Command("git", "fetch", "--tags")
		out, err := exec.Command("git", "tag").Output()
		if err != nil {
			return "", err
		}
		str := strings.TrimSuffix(string(out), "\n")
		tags := strings.Split(str, "\n")

		if len(tags) == 0 {
			// if no current flags exist then lets start at 0.0.0
			return "0.0.0", errors.New("No existing tags found")
		}

		// build an array of all the tags
		versionsRaw = make([]string, len(tags))
		for i, tag := range tags {
			if c.debug {
				fmt.Println(fmt.Sprintf("found tag %s", tag))
			}
			tag = strings.TrimPrefix(tag, "v")
			if tag != "" {
				versionsRaw[i] = tag
			}
		}

	}

	// turn the array into a new collection of versions that we can sort
	versions := make([]*version.Version, len(versionsRaw))
	for i, raw := range versionsRaw {
		v, _ := version.NewVersion(raw)
		versions[i] = v
	}

	if len(versions) == 0 {
		// if no current flags exist then lets start at 0.0.0
		return "0.0.0", errors.New("No existing tags found")
	}

	// return the latest tag
	sort.Sort(version.Collection(versions))
	latest := len(versions)
	if versions[latest-1] == nil {
		return "0.0.0", errors.New("No existing tags found")
	}
	return versions[latest-1].String(), nil
}

func getNewVersionFromTag(c config) (string, error) {

	// get the latest github tag
	tag, err := getLatestTag(c)
	if err != nil && tag == "" {
		return "", err
	}
	sv, err := semver.NewVersion(tag)
	if err != nil {
		return "", err
	}

	sv.BumpPatch()

	majorVersion := sv.Major
	minorVersion := sv.Minor
	patchVersion := sv.Patch

	// check if major or minor version has been changed
	baseVersion, err := getVersion(c)
	if err != nil {
		return fmt.Sprintf("%d.%d.%d", majorVersion, minorVersion, patchVersion), nil
	}

	// first use go-version to turn into a proper version, this handles 1.0-SNAPSHOT which semver doesn't
	tmpVersion, err := version.NewVersion(baseVersion)
	bsv, err := semver.NewVersion(tmpVersion.String())
	if err != nil {
		return "", err
	}
	baseMajorVersion := bsv.Major
	baseMinorVersion := bsv.Minor
	basePatchVersion := bsv.Patch

	if baseMajorVersion > majorVersion ||
		(baseMajorVersion == majorVersion &&
			(baseMinorVersion > minorVersion) || (baseMinorVersion == minorVersion && basePatchVersion > patchVersion)) {
		majorVersion = baseMajorVersion
		minorVersion = baseMinorVersion
		patchVersion = basePatchVersion
	}
	return fmt.Sprintf("%d.%d.%d", majorVersion, minorVersion, patchVersion), nil
}

// returns a string array containing the git owner and repo name for a given URL
func getCurrentGitOwnerRepo(url string) []string {
	var OwnerNameRegexp = regexp.MustCompile(`([^:]+)(/[^\/].+)?$`)

	matched2 := OwnerNameRegexp.FindStringSubmatch(url)
	s := strings.TrimSuffix(matched2[0], ".git")

	return strings.Split(s, "/")
}
