# semver release number

Returns a new release version based on previous git commits that can be used in a new release.

This is a simple binary that can be used in CD pipelines to read Java, Makefile files ans return an 'patch' incremented version.


If you need to bump the major or minor version simply increment the version in the Makefile / pom.xml in your master branch

e.g.

## Makefile

```$xslt
VERSION := 2.0.0-SNAPSHOT
```


## pom.xml

```xml
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/maven-v4_0_0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <groupId>io.example</groupId>
    <artifactId>example</artifactId>
    <version>1.0-0-SNAPSHOT</version>
    <packaging>pom</packaging>
</project>
```