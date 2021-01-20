package main

import (
	"bufio"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/jenkins-x-plugins/jx-release-version/adapters"
	"github.com/jenkins-x-plugins/jx-release-version/domain"
	"github.com/jenkins-x-plugins/jx-release-version/semrel"

	"github.com/jenkins-x/jx-helpers/v3/pkg/cmdrunner"
	"github.com/jenkins-x/jx-helpers/v3/pkg/gitclient/cli"

	"github.com/Masterminds/semver"
	"github.com/hashicorp/go-version"
)

// Version is the build version
var Version string

// GitTag is the git tag of the build
var GitTag string

// BuildDate is the date when the build was created
var BuildDate string

type Project struct {
	Version string `xml:"version"`
}

type config struct {
	dryrun          bool
	debug           bool
	gitTag          bool
	dir             string
	ghOwner         string
	ghRepository    string
	samerelease     bool
	baseVersion     string
	minor           bool
	semanticRelease bool
}

func main() {

	debug := flag.Bool("debug", false, "prints debug into to console")
	dir := flag.String("folder", ".", "the folder to look for files that contain a pom.xml or Makefile with the project version to bump")
	owner := flag.String("gh-owner", "", "a github repository owner if not running from within a git project  e.g. fabric8io")
	repo := flag.String("gh-repository", "", "a git repository if not running from within a git project  e.g. fabric8")
	baseVersion := flag.String("base-version", "", "use this instead of Makefile, pom.xml, etc, e.g. -base-version=2.0.0-SNAPSHOT")
	samerelease := flag.Bool("same-release", false, "for support old releases: for example 7.0.x and tag for new release 7.1.x already exist, with `-same-release` argument next version from 7.0.x will be returned ")
	ver := flag.Bool("version", false, "prints the version")
	minor := flag.Bool("minor", false, "increase minor version instead of patch")
	gitTag := flag.Bool("use-git-tag", false, "use only git tag to derive next release version")
	semanticRelease := flag.Bool("semantic-release", false, "use conventional commits to derive next release version")
	flag.Parse()

	if *ver {
		printVersion()
		os.Exit(0)
	}

	c := config{
		debug:           *debug,
		gitTag:          *gitTag,
		dir:             *dir,
		ghOwner:         *owner,
		ghRepository:    *repo,
		samerelease:     *samerelease,
		baseVersion:     *baseVersion,
		minor:           *minor,
		semanticRelease: *semanticRelease,
	}

	if c.debug {
		if os.Getenv("JX_LOG_LEVEL") == "" {
			os.Setenv("JX_LOG_LEVEL", "debug")
		}
		fmt.Println("available environment:")
		for _, e := range os.Environ() {
			fmt.Println(e)
		}
	}

	gitHubClient := adapters.NewGitHubClient(c.debug)
	v, err := getNewVersionFromTag(c, gitHubClient)
	if err != nil {
		fmt.Println("failed to get new version", err)
		os.Exit(-1)
	}
	fmt.Printf("%s", v)
}

func printVersion() {
	fmt.Printf(`Version: %s
Git Tag: %s
Build Date: %s
`, Version, GitTag, BuildDate)
}

func getVersion(c config) (string, error) {
	if c.baseVersion != "" {
		return c.baseVersion, nil
	}
	chart, err := ioutil.ReadFile(filepath.Join(c.dir, "Chart.yaml"))
	if err == nil {
		if c.debug {
			fmt.Println("Found Chart.yaml")
		}
		scanner := bufio.NewScanner(strings.NewReader(string(chart)))
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "version") {
				parts := strings.Split(scanner.Text(), ":")
				v := strings.TrimSpace(parts[1])
				if v != "" {
					if c.debug {
						fmt.Println(fmt.Sprintf("existing Chart version %v", v))
					}
					return v, nil
				}
			}
		}
	}

	m, err := ioutil.ReadFile(filepath.Join(c.dir, "Makefile"))
	if err == nil {
		if c.debug {
			fmt.Println("Found Makefile")
		}
		scanner := bufio.NewScanner(strings.NewReader(string(m)))
		for scanner.Scan() {
			if strings.HasPrefix(scanner.Text(), "VERSION") || strings.HasPrefix(scanner.Text(), "VERSION ") || strings.HasPrefix(scanner.Text(), "VERSION:") || strings.HasPrefix(scanner.Text(), "VERSION=") {
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

	am, err := ioutil.ReadFile(filepath.Join(c.dir, "configure.ac"))
	if err == nil {
		if c.debug {
			fmt.Println("configure.ac")
		}

		scanner := bufio.NewScanner(strings.NewReader(string(am)))
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "AC_INIT") {
				re := regexp.MustCompile("AC_INIT\\s*\\(([^\\s]+),\\s*([.\\d]+(-\\w+)?).*\\)")
				matched := re.FindStringSubmatch(scanner.Text())
				v := strings.TrimSpace(matched[2])
				if v != "" {
					if c.debug {
						fmt.Println(fmt.Sprintf("existing configure.ac version %v", v))
					}
					return v, nil
				}
			}
		}
	}

	cm, err := ioutil.ReadFile(filepath.Join(c.dir, "CMakeLists.txt"))
	if err == nil {
		if c.debug {
			fmt.Println("CMakeLists.txt")
		}

		scanner := bufio.NewScanner(strings.NewReader(string(cm)))
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), " VERSION ") {
				re := regexp.MustCompile("project\\s*(([^\\s]+)\\s+VERSION\\s+([.\\d]+(-\\w+)?).*)")
				matched := re.FindStringSubmatch(scanner.Text())
				v := strings.TrimSpace(matched[3])
				if v != "" {
					if c.debug {
						fmt.Println(fmt.Sprintf("existing CMakeLists.txt version %v", v))
					}
					return v, nil
				}
			}
		}
	}

	s, err := ioutil.ReadFile(c.dir + string(filepath.Separator) + "setup.py")
	if err == nil {
		if c.debug {
			fmt.Println("Found setup.py")
		}

		// Regex to find the call to `setup(..., version='1.2.3', ...)`
		re := regexp.MustCompile("setup\\((.|\\n)*version\\s*=\\s*'(\\d|\\.)*'([^\\)]|\\n)*\\)")
		setup_call_bytes := re.Find(s)

		if len(setup_call_bytes) > 0 {

			// Regex to find the argument `version='1.2.3'`
			version_re := regexp.MustCompile("version\\s*=\\s*'(\\d*|\\.)*'")

			version := string(version_re.Find(setup_call_bytes))

			parts := strings.Split(strings.Replace(version, " ", "", -1), "=")
			v := strings.TrimPrefix(strings.TrimSuffix(parts[1], "'"), "'")

			if v != "" {
				if c.debug {
					fmt.Println(fmt.Sprintf("existing Makefile version %v", v))
				}
				return v, nil
			}
		}
	}

	p, err := ioutil.ReadFile(filepath.Join(c.dir, "pom.xml"))
	if err == nil {
		if c.debug {
			fmt.Println("found pom.xml")
		}
		var project Project
		_ = xml.Unmarshal(p, &project)
		if project.Version != "" {
			if c.debug {
				fmt.Println(fmt.Sprintf("existing version %v", project.Version))
			}
			return project.Version, nil
		}
	}

	pkg, err := ioutil.ReadFile(filepath.Join(c.dir, "package.json"))
	if err == nil {
		if c.debug {
			fmt.Println("found package.json")
		}
		var project Project
		_ = json.Unmarshal(pkg, &project)
		if project.Version != "" {
			if c.debug {
				fmt.Println(fmt.Sprintf("existing version %v", project.Version))
			}
			return project.Version, nil
		}
	}

	// Groovy syntax
	gradle, err := ioutil.ReadFile(filepath.Join(c.dir, "build.gradle"))
	if err == nil {
		if c.debug {
			fmt.Println("found build.gradle")
		}
		v := getGradleVersion(c, gradle)
		if v != "" {
			if c.debug {
				fmt.Println(fmt.Sprintf("existing version %v", v))
			}
			return v, nil
		}
	}
	// Kotlin syntax
	gradle_kts, err := ioutil.ReadFile(filepath.Join(c.dir, "build.gradle.kts"))
	if err == nil {
		if c.debug {
			fmt.Println("found build.gradle.kts")
		}
		v := getGradleVersion(c, gradle_kts)
		if v != "" {
			if c.debug {
				fmt.Println(fmt.Sprintf("existing version %v", v))
			}
			return v, nil
		}
	}

	return "0.0.0", errors.New("no recognised file to obtain current version from")
}

func getGradleVersion(c config, gradle []byte) string {
	scanner := bufio.NewScanner(strings.NewReader(string(gradle)))
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "version") {
			re := regexp.MustCompile("^version\\s*=\\s*['\"]([.\\d]+(-\\w+)?)['\"]")
			matched := re.FindStringSubmatch(scanner.Text())
			if len(matched) > 0 {
				return strings.TrimSpace(matched[1])
			}
		}
	}
	return ""
}

func getLatestTag(c config, gitClient domain.GitClient) (string, error) {
	// Get base version from file, will fallback to 0.0.0 if not found.
	baseVersion := "0.0.0"
	if !c.gitTag {
		baseVersion, _ = getVersion(c)
	}

	// if repo isn't provided by flags fall back to using current repo if run from a git project
	var versionsRaw []string
	if c.ghOwner != "" && c.ghRepository != "" {
		ctx := context.Background()

		tags, err := gitClient.ListTags(ctx, c.ghOwner, c.ghRepository)

		if err != nil {
			return "", err
		}
		if len(tags) == 0 {
			// if no current flags exist then lets start at base version
			return baseVersion, errors.New("No existing tags found")
		}

		// build an array of all the tags
		versionsRaw = make([]string, len(tags))
		for i, tag := range tags {
			if c.debug {
				fmt.Println(fmt.Sprintf("found remote tag %s", tag.Name))
			}
			versionsRaw[i] = tag.Name
		}
	} else {
		_, err := exec.LookPath("git")
		if err != nil {
			return "", fmt.Errorf("error running git: %v", err)
		}
		cmd := exec.Command("git", "fetch", "--tags", "-v")
		cmd.Env = append(cmd.Env, os.Environ()...)
		cmd.Dir = c.dir
		err = cmd.Run()
		if err != nil {
			return baseVersion, fmt.Errorf("error fetching tags: %v", err)
		}

		cmd = exec.Command("git", "tag")
		cmd.Dir = c.dir
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}
		str := strings.TrimSuffix(string(out), "\n")
		tags := strings.Split(str, "\n")

		if len(tags) == 0 {
			// if no current flags exist then lets start at base version
			return baseVersion, errors.New("No existing tags found")
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
	var versions []*version.Version
	for _, raw := range versionsRaw {
		// if same-release argument is set work only with versions which Major and Minor versions are the same
		if c.samerelease {
			same, _ := isMajorMinorTheSame(baseVersion, raw)
			if same {
				v, _ := version.NewVersion(raw)
				if v != nil {
					versions = append(versions, v)
				}
			}
		} else {
			v, _ := version.NewVersion(raw)
			if v != nil {
				versions = append(versions, v)
			}
		}
	}

	if len(versions) == 0 {
		// if no current flags exist then lets start at base version
		return baseVersion, errors.New("No existing tags found")
	}

	// return the latest tag
	col := version.Collection(versions)
	if c.debug {
		fmt.Printf("version collection %v \n", col)
	}

	sort.Sort(col)
	latest := len(versions)
	if versions[latest-1] == nil {
		return baseVersion, errors.New("No existing tags found")
	}
	return versions[latest-1].String(), nil
}

func getNewVersionFromTag(c config, gitClient domain.GitClient) (string, error) {
	var sv *semver.Version

	if c.semanticRelease {
		var err error
		sv, err = getSemanticReleaseVersion(c)
		if err != nil {
			return "", err
		}
	}

	if sv == nil {
		// get the latest github tag
		tag, err := getLatestTag(c, gitClient)
		if err != nil && tag == "" {
			return "", err
		}
		sv, err = semver.NewVersion(tag)
		if err != nil {
			return "", err
		}

		var newVersion semver.Version
		if c.minor {
			newVersion = sv.IncMinor()
		} else {
			newVersion = sv.IncPatch()
		}
		sv = &newVersion
	}

	majorVersion := sv.Major()
	minorVersion := sv.Minor()
	patchVersion := sv.Patch()

	// check if major or minor version has been changed
	baseVersion, err := getVersion(c)
	if err != nil {
		return fmt.Sprintf("%d.%d.%d", majorVersion, minorVersion, patchVersion), nil
	}

	// first use go-version to turn into a proper version, this handles 1.0-SNAPSHOT which semver doesn't
	tmpVersion, err := version.NewVersion(baseVersion)
	if err != nil {
		return fmt.Sprintf("%d.%d.%d", majorVersion, minorVersion, patchVersion), nil
	}
	bsv, err := semver.NewVersion(tmpVersion.String())
	if err != nil {
		return "", err
	}
	baseMajorVersion := bsv.Major()
	baseMinorVersion := bsv.Minor()
	basePatchVersion := bsv.Patch()

	if baseMajorVersion > majorVersion ||
		(baseMajorVersion == majorVersion &&
			(baseMinorVersion > minorVersion) || (baseMinorVersion == minorVersion && basePatchVersion > patchVersion)) {
		majorVersion = baseMajorVersion
		minorVersion = baseMinorVersion
		patchVersion = basePatchVersion
	}
	return fmt.Sprintf("%d.%d.%d", majorVersion, minorVersion, patchVersion), nil
}

func getSemanticReleaseVersion(c config) (*semver.Version, error) {
	cmdRunner := cmdrunner.QuietCommandRunner
	if c.debug {
		cmdRunner = cmdrunner.DefaultCommandRunner
	}
	gitClient := cli.NewCLIClient("git", cmdRunner)
	return semrel.GetNewVersion(c.dir, gitClient)
}

func isMajorMinorTheSame(v1 string, v2 string) (bool, error) {
	sv1, err1 := semver.NewVersion(v1)
	if err1 != nil {
		return false, err1
	}
	sv2, err2 := semver.NewVersion(v2)
	if err2 != nil {
		return false, err2
	}
	if sv1.Major() != sv2.Major() {
		return false, nil
	}
	if sv1.Minor() != sv2.Minor() {
		return false, nil
	}
	return true, nil
}
