### Linux

```shell
curl -L https://github.com/jenkins-x-plugins/jx-release-version/releases/download/v2.2.6/jx-release-version-linux-amd64.tar.gz | tar xzv 
sudo mv jx-release-version /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/jenkins-x-plugins/jx-release-version/releases/download/v2.2.6/jx-release-version-darwin-amd64.tar.gz | tar xzv
sudo mv jx-release-version /usr/local/bin
```

## Changes

### Bug Fixes

* commits iteration (Vincent Behar)

### Documentation

* added some documentation for the -strip-prerelease flag (Gareth Evans)

### Chores

* move strip-release from global flag to next-version (Gareth Evans)
* add ability to strip prerelease information from the version (Gareth Evans)
