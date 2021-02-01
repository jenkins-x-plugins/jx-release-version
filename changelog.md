### Linux

```shell
curl -L https://github.com/jenkins-x-plugins/jx-release-version/releases/download/v2.1.0/jx-release-version-linux-amd64.tar.gz | tar xzv 
sudo mv jx-release-version /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/jenkins-x-plugins/jx-release-version/releases/download/v2.1.0/jx-release-version-darwin-amd64.tar.gz | tar xzv
sudo mv jx-release-version /usr/local/bin
```

## Changes

### New Features

* the output format is now configurable (Vincent Behar) [#52](https://github.com/jenkins-x-plugins/jx-release-version/issues/52) 

### Chores

* restore release pipeline (Vincent Behar)

### Issues

* [#52](https://github.com/jenkins-x-plugins/jx-release-version/issues/52) Support 1 or 2 digit version syntax
