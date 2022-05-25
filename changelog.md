### Linux

```shell
curl -L https://github.com/jenkins-x-plugins/jx-release-version/releases/download/v2.5.2/jx-release-version-linux-amd64.tar.gz | tar xzv 
sudo mv jx-release-version /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/jenkins-x-plugins/jx-release-version/releases/download/v2.5.2/jx-release-version-darwin-amd64.tar.gz | tar xzv
sudo mv jx-release-version /usr/local/bin
```

## Changes

### Bug Fixes

* support for tekton 0.28+ (Vincent Behar)
* don't fail when fetching tags if already up-to-date (Vincent Behar)

### Chores

* upgrade image to 2.5.1 (jenkins-x-bot)
