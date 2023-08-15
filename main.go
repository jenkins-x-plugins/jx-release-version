package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/semver/v3"
	"github.com/Masterminds/sprig/v3"
	"github.com/jenkins-x-plugins/jx-release-version/v2/pkg/strategy"
	"github.com/jenkins-x-plugins/jx-release-version/v2/pkg/strategy/auto"
	"github.com/jenkins-x-plugins/jx-release-version/v2/pkg/strategy/fromfile"
	"github.com/jenkins-x-plugins/jx-release-version/v2/pkg/strategy/fromtag"
	"github.com/jenkins-x-plugins/jx-release-version/v2/pkg/strategy/increment"
	"github.com/jenkins-x-plugins/jx-release-version/v2/pkg/strategy/manual"
	"github.com/jenkins-x-plugins/jx-release-version/v2/pkg/strategy/semantic"
	"github.com/jenkins-x-plugins/jx-release-version/v2/pkg/tag"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
)

var (
	// these are set at compile time through LD Flags
	Version  = "dev"
	Revision = "unknown"
	Date     = "now"

	// CLI flag options
	options struct {
		printVersion    bool
		debug           bool
		dir             string
		previousVersion string
		commitHeadlines string
		nextVersion     string
		outputFormat    string
		tag             bool
		tagPrefix       string
		pushTag         bool
		fetchTags       bool
		gitName         string
		gitEmail        string
	}
)

func init() {
	wd, _ := os.Getwd()
	flag.StringVar(&options.dir, "dir", wd, "The directory that contains the git repository. Default to the current working directory.")
	flag.StringVar(&options.previousVersion, "previous-version", getEnvWithDefault("PREVIOUS_VERSION", "auto"), "The strategy to detect the previous version: auto, from-tag, from-file or manual. Default to the PREVIOUS_VERSION env var.")
	flag.StringVar(&options.commitHeadlines, "commit-headlines", getEnvWithDefault("COMMIT_HEADLINES", ""), "The commit headline(s) to use for semantic next version instead of the commit()s of a repository. Default to empty.")
	flag.StringVar(&options.nextVersion, "next-version", getEnvWithDefault("NEXT_VERSION", "auto"), "The strategy to calculate the next version: auto, semantic, from-file, increment or manual. Default to the NEXT_VERSION env var.")
	flag.StringVar(&options.outputFormat, "output-format", getEnvWithDefault("OUTPUT_FORMAT", "{{.Major}}.{{.Minor}}.{{.Patch}}"), "The output format of the next version. Default to the OUTPUT_FORMAT env var.")
	flag.BoolVar(&options.debug, "debug", os.Getenv("JX_LOG_LEVEL") == "debug", "Print debug logs. Enabled by default if the JX_LOG_LEVEL env var is set to 'debug'.")
	flag.BoolVar(&options.printVersion, "version", false, "Just print the version and do nothing.")
	flag.BoolVar(&options.tag, "tag", os.Getenv("TAG") == "true", "Perform a git tag")
	flag.StringVar(&options.tagPrefix, "tag-prefix", getEnvWithDefault("TAG_PREFIX", "v"), "Prefix to use for the git tag")
	flag.BoolVar(&options.pushTag, "push-tag", getEnvWithDefault("PUSH_TAG", "true") == "true", "Use with tag flag, pushes a git tag to the remote branch")
	flag.BoolVar(&options.fetchTags, "fetch-tags", getEnvWithDefault("FETCH_TAGS", "") == "true", "Fetch tags from the remote origin before detecting the previous version")
	flag.StringVar(&options.gitName, "git-user", getEnvWithDefault("GIT_NAME", ""), "Name is the personal name of the author and the committer of a commit, use to override Git config")
	flag.StringVar(&options.gitEmail, "git-email", getEnvWithDefault("GIT_EMAIL", ""), "Email is the email of the author and the committer of a commit, use to override Git config")
}

func main() {
	flag.Parse()

	if options.printVersion {
		fmt.Printf("Version %s - Revision %s - Date %s", Version, Revision, Date)
		return
	}

	if options.debug {
		os.Setenv("JX_LOG_LEVEL", "debug")
		log.Logger().Debugf("jx-release-version %s running in debug mode in %s", Version, options.dir)
	}

	previousVersion, err := versionReader().ReadVersion()
	if err != nil {
		log.Logger().Fatalf("Failed to read previous version using %q: %v", options.previousVersion, err)
	}
	log.Logger().Debugf("Previous version: %s", previousVersion.String())

	nextVersion, err := versionBumper().BumpVersion(*previousVersion)
	if err != nil {
		log.Logger().Fatalf("Failed to bump version using %q: %v", options.nextVersion, err)
	}
	log.Logger().Debugf("Next version: %s", nextVersion.String())

	output, err := formatVersion(*nextVersion)
	if err != nil {
		log.Logger().Fatalf("Failed to format version %q with %q: %v", *nextVersion, options.outputFormat, err)
	}

	fmt.Print(output)

	if options.tag {
		tagOptions := tag.Tag{
			FormattedVersion: options.tagPrefix + output,
			Dir:              options.dir,
			PushTag:          options.pushTag,
			GitName:          options.gitName,
			GitEmail:         options.gitEmail,
		}
		err = tagOptions.TagRemote()
		if err != nil {
			log.Logger().Fatalf("Failed to tag using version %s: %v", output, err)
		}
	}
}

func versionReader() strategy.VersionReader {
	var (
		versionReader             strategy.VersionReader
		strategyName, strategyArg string
	)

	parts := strings.SplitN(options.previousVersion, ":", 2)
	strategyName = parts[0]
	if len(parts) > 1 {
		strategyArg = parts[1]
	}

	switch strategyName {
	case "auto", "":
		versionReader = auto.Strategy{
			FromTagStrategy: fromtag.Strategy{
				Dir:       options.dir,
				FetchTags: options.fetchTags,
			},
		}
	case "from-tag":
		versionReader = fromtag.Strategy{
			Dir:        options.dir,
			TagPattern: strategyArg,
			FetchTags:  options.fetchTags,
		}
	case "from-file":
		versionReader = fromfile.Strategy{
			Dir:      options.dir,
			FilePath: strategyArg,
		}
	case "manual":
		versionReader = manual.Strategy{
			Version: strategyArg,
		}
	default:
		versionReader = manual.Strategy{
			Version: options.previousVersion,
		}
	}

	log.Logger().Debugf("Using %q version reader (with %q)", strategyName, strategyArg)
	return versionReader
}

func versionBumper() strategy.VersionBumper {
	var (
		versionBumper             strategy.VersionBumper
		strategyName, strategyArg string
	)

	parts := strings.SplitN(options.nextVersion, ":", 2)
	strategyName = parts[0]
	if len(parts) > 1 {
		strategyArg = parts[1]
	}

	switch strategyName {
	case "auto", "":
		versionBumper = auto.Strategy{
			SemanticStrategy: semantic.Strategy{
				Dir:                   options.dir,
				StripPrerelease:       strings.Contains(strategyArg, "strip-prerelease"),
				CommitHeadlinesString: options.commitHeadlines,
			},
		}
	case "semantic":
		versionBumper = semantic.Strategy{
			Dir:                   options.dir,
			StripPrerelease:       strings.Contains(strategyArg, "strip-prerelease"),
			CommitHeadlinesString: options.commitHeadlines,
		}
	case "from-file":
		versionBumper = fromfile.Strategy{
			Dir:      options.dir,
			FilePath: strategyArg,
		}
	case "increment":
		versionBumper = increment.Strategy{
			ComponentToIncrement: strategyArg,
		}
	case "manual":
		versionBumper = manual.Strategy{
			Version: strategyArg,
		}
	default:
		versionBumper = manual.Strategy{
			Version: options.previousVersion,
		}
	}

	log.Logger().Debugf("Using %q version bumper (with %q)", strategyName, strategyArg)
	return versionBumper
}

func formatVersion(version semver.Version) (string, error) {
	outputTemplate, err := template.New("output").Funcs(sprig.TxtFuncMap()).Parse(options.outputFormat)
	if err != nil {
		return "", err
	}

	output := new(strings.Builder)
	err = outputTemplate.Execute(output, version)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

func getEnvWithDefault(key string, defaultVal string) string {
	if val, found := os.LookupEnv(key); found {
		return val
	}
	return defaultVal
}
