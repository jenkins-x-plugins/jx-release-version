FROM alpine:3.14

COPY ./build/linux/jx-release-version /usr/bin/jx-release-version
COPY ./hack/github-actions-entrypoint.sh /usr/bin/github-actions-entrypoint.sh

ENTRYPOINT ["jx-release-version"]
