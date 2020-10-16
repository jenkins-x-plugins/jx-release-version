FROM gcr.io/jenkinsxio/jx-cli-base:0.0.21

ENTRYPOINT ["jx-release-version"]

COPY ./build/linux/jx-release-version /usr/bin/jx-release-version