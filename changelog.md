### Linux

```shell
curl -L https://github.com/jenkins-x-plugins/jx-release-version/releases/download/v2.2.4/jx-release-version-linux-amd64.tar.gz | tar xzv 
sudo mv jx-release-version /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/jenkins-x-plugins/jx-release-version/releases/download/v2.2.4/jx-release-version-darwin-amd64.tar.gz | tar xzv
sudo mv jx-release-version /usr/local/bin
```

## Changes

### Chores

* upgrade go dependencies (jenkins-x-bot)

### Other Changes

These commits did not use [Conventional Commits](https://conventionalcommits.org/) formatted messages:

* add gradle-properties support (pg2000)
