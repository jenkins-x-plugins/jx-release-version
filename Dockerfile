FROM scratch

ENTRYPOINT ["/semver-release-version"]

COPY ./bin/semver-release-version-linux /semver-release-version