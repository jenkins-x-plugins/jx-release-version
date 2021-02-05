### Linux

```shell
curl -L https://github.com/jenkins-x-plugins/jx-release-version/releases/download/v2.2.2/jx-release-version-linux-amd64.tar.gz | tar xzv 
sudo mv jx-release-version /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/jenkins-x-plugins/jx-release-version/releases/download/v2.2.2/jx-release-version-darwin-amd64.tar.gz | tar xzv
sudo mv jx-release-version /usr/local/bin
```

## Changes

### Documentation

* document tekton and github integrations (Vincent Behar)

### Chores

* upgrade go dependencies (jenkins-x-bot)
