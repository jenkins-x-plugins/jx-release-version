package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/jenkins-x-plugins/jx-release-version/v2/strategy"
	"github.com/jenkins-x-plugins/jx-release-version/v2/strategy/auto"
	"github.com/jenkins-x-plugins/jx-release-version/v2/strategy/fromfile"
	"github.com/jenkins-x-plugins/jx-release-version/v2/strategy/fromtag"
	"github.com/jenkins-x-plugins/jx-release-version/v2/strategy/increment"
	"github.com/jenkins-x-plugins/jx-release-version/v2/strategy/manual"
	"github.com/jenkins-x-plugins/jx-release-version/v2/strategy/semantic"
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
		nextVersion     string
	}
)

func init() {
	wd, _ := os.Getwd()
	flag.StringVar(&options.dir, "dir", wd, "The directory that contains the git repository. Default to the current working directory.")
	flag.StringVar(&options.previousVersion, "previous-version", getEnvWithDefault("PREVIOUS_VERSION", "auto"), "The strategy to detect the previous version: auto, from-tag, from-file or manual. Default to the PREVIOUS_VERSION env var.")
	flag.StringVar(&options.nextVersion, "next-version", getEnvWithDefault("NEXT_VERSION", "auto"), "The strategy to calculate the next version: auto, semantic, from-file, increment or manual. Default to the NEXT_VERSION env var.")
	flag.BoolVar(&options.debug, "debug", os.Getenv("JX_LOG_LEVEL") == "debug", "Print debug logs. Enabled by default if the JX_LOG_LEVEL env var is set to 'debug'.")
	flag.BoolVar(&options.printVersion, "version", false, "Just print the version and do nothing.")
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

	// ensure we don't keep pre-release or metadata information
	fmt.Printf("%d.%d.%d", nextVersion.Major(), nextVersion.Minor(), nextVersion.Patch())
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
				Dir: options.dir,
			},
		}
	case "from-tag":
		versionReader = fromtag.Strategy{
			Dir:        options.dir,
			TagPattern: strategyArg,
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
				Dir: options.dir,
			},
		}
	case "semantic":
		versionBumper = semantic.Strategy{
			Dir: options.dir,
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

func getEnvWithDefault(key string, defaultVal string) string {
	if val, found := os.LookupEnv(key); found {
		return val
	}
	return defaultVal
}
