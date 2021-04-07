### Linux

```shell
curl -L https://github.com/jenkins-x-plugins/jx-release-version/releases/download/v2.3.1/jx-release-version-linux-amd64.tar.gz | tar xzv 
sudo mv jx-release-version /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/jenkins-x-plugins/jx-release-version/releases/download/v2.3.1/jx-release-version-darwin-amd64.tar.gz | tar xzv
sudo mv jx-release-version /usr/local/bin
```

## Changes

### New Features

* lets optionally allow a flag to create and push a tag when generating the next release version (James Rawlings)

### Chores

* allow git token to be set when pushing tags (James Rawlings)
* move code into a pkg folder so we can structure code a bit better (James Rawlings)
