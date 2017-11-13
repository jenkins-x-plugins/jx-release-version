FROM scratch

ENTRYPOINT ["/semver-release-number"]

COPY ./bin/semver-release-number-linux-amd64 /semver-release-number