# Jenkins X Release Version

Returns a new release version based on previous git tags that can be used in a new release.

This is a simple binary that can be used in CD pipelines to read pom.xml or Makefile's and return a 'patch' incremented version.

If you need to bump the major or minor version simply increment the version in your Makefile / pom.xml


This helps in continuous delivery if you want an automatic release when a change is merged to master.  Traditional approaches mean the version is stored in a file that is checked and updated after each release.  If you want automatic releases this means you will get another release triggered from the version update resulting in a cyclic release sitiation.  

Using a git tag to work out the next release version is better than traditional approaches of storing it in a VERSION file or updating a pom.xml.  If a major or minor version increase is required then still update the file and `jx-release-version` will use you new version.

Please note that `jx-release-version` is not called from the Tekton-style build pipelines, these use `jx step` instead.

## Prerequisits

- `git` to be available on your `$PATH`

## Examples

- If your project is new or has no existing git tags then running `jx-release-version` will return a default version of `0.0.1`

- If your latest git tag is `1.2.3` and you Makefile or pom.xml is `1.2.0-SNAPSHOT` then `jx-release-version` will return `1.2.4`

- If your latest git tag is `1.2.3` and your Makefile or pom.xml is `2.0.0` then `jx-release-version` will return `2.0.0`

- If you need to support an old release for example 7.0.x and tags for new realese 7.1.x already exist, the `-same-release` flag  will help to obtain version from 7.0.x release. If the pom file version is 7.0.0-SNAPSHOT and both the 7.1.0 and 7.0.2 tags exist the command `jx-release-version` will return 7.1.1 but if we run `jx-release-version -same-release` it will return 7.0.3

- If you need to get a release version `1.1.0` for older release and your last tag is `1.2.3` please change your Makefile or pom.xml to `1.1.0-SNAPSHOT` and run `jx-release-version -same-release`

## Example Makefile

```$xslt
VERSION := 2.0.0-SNAPSHOT
```

## Example pom.xml

```xml
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/maven-v4_0_0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <groupId>io.example</groupId>
    <artifactId>example</artifactId>
    <version>1.0-0-SNAPSHOT</version>
    <packaging>pom</packaging>
</project>
```

Then in your release pipeline do something like this:

```sh
    ➜ RELEASE_VERSION=$(jx-release-version)
    ➜ echo "New release version ${RELEASE_VERSION}
    ➜ mvn versions:set -DnewVersion=${RELEASE_VERSION}
    ➜ git commit -a -m 'release ${RELEASE_VERSION}'
    ➜ git tag -fa v${RELEASE_VERSION} -m 'Release version ${RELEASE_VERSION}'
    ➜ git push origin v${RELEASE_VERSION}
```

### CLI arguments

```sh
  -base-version string
    	use this instead of Makefile, pom.xml, etc, e.g. -base-version=2.0.0-SNAPSHOT
  -debug
    prints debug into to console
  -folder string
    the folder to look for files that contain a pom.xml or Makefile with the project version to bump (default ".")
  -gh-owner string
    a github repository owner if not running from within a git project  e.g. fabric8io
  -gh-repository string
    a git repository if not running from within a git project  e.g. fabric8
  -same-release -same-release
    for support old releases: for example 7.0.x and tag for new release 7.1.x already exist, with -same-release argument next version from 7.0.x will be returned
  -version
    prints the version
```

### FAQ

__Why isn't a nodejs package.json supported?__

We use nodejs but make use of the semantic-release plugin which works out the next release versions instead

__Why only Makefiles and pom.xml supported?__

Right now we tend to only use golang, java and nodejs projects so if there's a file type missing please raise an issue or PR.
