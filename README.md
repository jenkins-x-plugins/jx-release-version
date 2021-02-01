# Jenkins X Release Version

**This is a simple and standalone binary that can be used in CD pipelines to calculate the next release version**.

By default it:
- retrieves the previous version - which is the latest git tag matching the [semver](https://semver.org/) spec
- and then use [conventional commits](https://www.conventionalcommits.org/) to calculate the next release version

But it also supports other strategies to read the previous version and calculate the next version.

## Usage

Just run `jx-release-version` in your project's top directory, and it should just print the next release version. It won't write anything to disk.

It accepts the following CLI flags:
- `-dir`: the location on the filesystem of your project's top directory - default to the current working directory.
- `-previous-version`: the [strategy to use to read the previous version](#reading-the-previous-version). Can also be set using the `PREVIOUS_VERSION` environment variable.
- `-next-version`: the [strategy to use to calculate the next version](#calculating—the-next-version). Can also be set using the `NEXT_VERSION` environment variable.
- `-debug`: if enabled, will print debug logs to stdout in addition to the next version. It can also be enabled by setting the `JX_LOG_LEVEL` environment variable to `debug`.

### Features

- standalone - no dependencies required. It uses an embedded [git implementation](https://github.com/go-git/go-git) to read the [Git](https://git-scm.com/) repository's information.
- simple configuration through CLI flags or environment variables.
- by default works even on an empty git repository
- multiple strategies to read the previous version and/or calculate the next version.

## Reading the previous version

There are different ways to read the previous version:

### Auto

The `auto` strategy is the default one. It tries to find the latest git tag, or if there there are no git tags, it just use `0.0.0` as the previous version.

**Usage**:
- `jx-release-version -previous-version=auto`
- `jx-release-version` - the `auto` strategy is already the default one

### From tag

The `from-tag` strategy uses the latest git tag as the previous version. Note that it only uses tags which matches the [semver](https://semver.org/) spec - other tags are just ignored.

Optionnaly, it can filter tags based on a given pattern: if you use `from-tag:v1` it will use the latest tag matching the `v1` pattern. Note that it uses [Go's stdlib regexp](https://golang.org/pkg/regexp/) - you can see the [syntax](https://golang.org/pkg/regexp/syntax/).
This feature can be used to maintain 2 major versions in parallel: for each, you just configure the right pattern, so that `jx-release-version` retrieves the right previous version, and bump it as it should.

Note that if it can't find a tag, it will fail.

**Usage**:
- `jx-release-version -previous-version=from-tag`
- `jx-release-version -previous-version=from-tag:v1`

### From file

The `from-file` strategy will read the previous version from a file. Supported formats are:
- **Helm Charts**, using the `Chart.yaml` file
- **Makefile**, using the `Makefile` file
- **Automake**, using the `configure.ac` file
- **CMake**, using the `CMakeLists.txt` file
- **Python**, using the `setup.py` file
- **Maven**, using the `pom.xml` file
- **Javascript**, using the `package.json` file
- **Gradle**, using the `build.gradle` or `build.gradle.kts` file

**Usage**:
- if you use `jx-release-version -previous-version=from-file` it will auto detect which file to use, trying the supported formats in the order in which they are listed.
- if you specify a file, it will use it to find the previous version. For example:
  - `jx-release-version -previous-version=from-file:pom.xml`
  - `jx-release-version -previous-version=from-file:charts/my-chart/Chart.yaml`
  - `jx-release-version -previous-version=from-file:Chart.yaml -dir=charts/my-chart`

### Manual

The `manual` strategy can be used if you already know the previous version, and just want `jx-release-version` to use it.

**Usage**:
- `jx-release-version -previous-version=manual:1.2.3`
- `jx-release-version -previous-version=1.2.3` - the `manual` prefix is optional

## Calculating the next version

There are different ways to calculate the next version:

### Auto

The `auto` strategy is the default one. It tries to use the `semantic` strategy, but if it can't find a tag for the previous version, it will fallback to incrementing the patch component.

**Usage**:
- `jx-release-version -next-version=auto`
- `jx-release-version` - the `auto` strategy is already the default one

### Semantic release

The `semantic` strategy finds all commits between the previous version's git tag and the current HEAD, and then uses [conventional commits](https://www.conventionalcommits.org/) to parse them. If it finds:
- at least 1 commit with a `BREAKING CHANGE: ` footer, then it will bump the major component of the version
- at least 1 commit with a `feat:` prefix, then it will bump the minor component of the version
- otherwise it will bump the patch component of the version

Note that if it can't find a tag for the previous version, it will fail.

**Usage**:
- `jx-release-version -next-version=semantic`

### From file

The `from-file` strategy will read the next version from a file. Supported formats are:
- **Helm Charts**, using the `Chart.yaml` file
- **Makefile**, using the `Makefile` file
- **Automake**, using the `configure.ac` file
- **CMake**, using the `CMakeLists.txt` file
- **Python**, using the `setup.py` file
- **Maven**, using the `pom.xml` file
- **Javascript**, using the `package.json` file
- **Gradle**, using the `build.gradle` or `build.gradle.kts` file

**Usage**:
- if you use `jx-release-version -next-version=from-file` it will auto detect which file to use, trying the supported formats in the order in which they are listed.
- if you specify a file, it will use it to find the next version. For example:
  - `jx-release-version -next-version=from-file:pom.xml`
  - `jx-release-version -next-version=from-file:charts/my-chart/Chart.yaml`
  - `jx-release-version -next-version=from-file:Chart.yaml -dir=charts/my-chart`

### Increment

The `increment` strategy can be used if you want to increment a specific component of the version.

**Usage**:
- `jx-release-version -next-version=increment:major`
- `jx-release-version -next-version=increment:minor`
- `jx-release-version -next-version=increment:patch`
- `jx-release-version -next-version=increment` - by default it will increment the patch component

### Manual

The `manual` strategy can be used if you already know the next version, and just want `jx-release-version` to use it.

**Usage**:
- `jx-release-version -next-version=manual:1.2.3`
- `jx-release-version -next-version=1.2.3` - the `manual` prefix is optional
