### Linux

```shell
curl -L https://github.com/jenkins-x-plugins/jx-release-version/releases/download/v2.2.0/jx-release-version-linux-amd64.tar.gz | tar xzv 
sudo mv jx-release-version /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/jenkins-x-plugins/jx-release-version/releases/download/v2.2.0/jx-release-version-darwin-amd64.tar.gz | tar xzv
sudo mv jx-release-version /usr/local/bin
```

## Changes

### New Features

* add support for github actions (Vincent Behar)
* the output format is now configurable (Vincent Behar) [#52](https://github.com/jenkins-x-plugins/jx-release-version/issues/52) 

### Issues

* [#52](https://github.com/jenkins-x-plugins/jx-release-version/issues/52) Support 1 or 2 digit version syntax
