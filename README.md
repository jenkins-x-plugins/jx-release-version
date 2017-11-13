# semver release number

Returns a new release version based on previous git tags that can be used in a new release.

This is a simple binary that can be used in CD pipelines to read pom.xml or Makefile's and return an 'patch' incremented version.

If you need to bump the major or minor version simply increment the version in your Makefile / pom.xml


This helps in continuous delivery if you want an automatic release when a change is merged to master.  Traditional approaches mean the version is stored in a file that is checked and updated after each release.  If you want autotic releases this means you will get another release as a result of the version number update resulting in a cyclic release sitiation.  

Using a git tag to work out the next release version is better than traditional approaches of storing it in a a VERSION file or updating a pom.xml.  If a major or minor version increase is required then still update the file and `semver-release-number` will use you new version. 

### Examples

- If your project is new or has no existing git tags then running `semver-release-number` will return a default version of `0.0.1`

- If your latest git tag is `1.2.3` and you Makefile or pom.xml is `1.2.0-SNAPSHOT` then `semver-release-number` will return `1.2.4`

- If your latest git tag is `1.2.3` and your Makefile or pom.xml is `2.0.0` then `semver-release-number` will return `2.0.0`

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

### FAQ

__Why isn't a nodejs package.json supported?__

We use nodejs but make use of the semantic-release plugin which works out the next release versions instead

__Why only Makefiles and pom.xml supported?__

Right now we tend to only use golang, java and nodejs projects so if there's a file type missing please raise an issue or PR.