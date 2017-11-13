# semver release number

Returns a new release version based on previous git tags that can be used in a new release.

This is a simple binary that can be used in CD pipelines to read pom.xml or Makefile's and return an 'patch' incremented version.

If you need to bump the major or minor version simply increment the version in your Makefile / pom.xml


Using a git tag to work out the next release version is better than traditional approaches of storing it in a a VERSION file or updating a pom.xml.    

### Examples

If your project is new or has no existing git tags then running `semver-release-number` will return a default version of `0.0.1`

If your latest git tag is `1.2.3` then `semver-release-number` will return `1.2.4`

If your latest git tag is `1.2.3` and your Makefile or pom.xml is `2.0.0` then `semver-release-number` will return `2.0.0`

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